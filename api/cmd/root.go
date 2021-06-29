package cmd

import (
	"context"
	"github.com/apex/log"
	"github.com/c3sr/config"
	dlframework "github.com/c3sr/dlframework/http"
	"github.com/c3sr/dlframework/httpapi/restapi"
	"github.com/c3sr/dlframework/httpapi/restapi/operations"
	"github.com/c3sr/tracer"
	tracermiddleware "github.com/c3sr/tracer/middleware"
	"github.com/c3sr/uuid"
	"github.com/go-openapi/loads"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os/signal"
	"time"

	"fmt"
	"net/http"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "Run the prediction API endpoints",
	Long:  "Run the prediction API endpoints",
	Run:   serve,
}

const DefaultPort = "80"

func serve(cmd *cobra.Command, args []string) {
	port := viper.Get("Port")
	e := echo.New()
	handler, err := getDlframeworkHandler()

	if err != nil {
		panic(fmt.Errorf("unable to create handler: %s", err))
	}

	configureEcho(e)
	addRoutes(e, handler)

	go func() {
		err = e.Start(fmt.Sprintf(":%s", port))
		if err != nil {
			panic(fmt.Errorf("unable to listen: %s", err))
			return
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Failed to gracefully shutdown server")
	}
}

func configureEcho(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: middleware.DefaultLoggerConfig.Format,
	}))
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.NewV4()
		},
	}))
	e.Use(AllowTracedHeaders())
}

func addRoutes(e *echo.Echo, handler http.Handler) {
	api := e.Group("/api")
	api.Any("/predict*",
		echo.WrapHandler(handler),
		tracermiddleware.FromHTTPRequest(tracer.Std(), "api_request"),
		tracermiddleware.ToHTTPResponse(tracer.Std()),
	)
	api.Any("/*", echo.WrapHandler(handler))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		opts := []config.Option{
			config.AppName("carml"),
		}
		config.Init(opts...)
	})

	viper.SetEnvPrefix("c3sr")
	viper.BindEnv("Port", "port")
	viper.SetDefault("Port", DefaultPort)
	viper.BindEnv("app.name", "app_name")
	viper.SetDefault("app.name", "carml")
	viper.BindEnv("registry.provider", "registry_provider")
	viper.BindEnv("registry.endpoints", "registry_endpoints")
	viper.BindEnv("tracer.enabled", "tracer_enabled")
	viper.BindEnv("tracer.provider", "tracer_provider")
	viper.BindEnv("tracer.endpoints", "tracer_endpoints")
	viper.BindEnv("tracer.level", "tracer_level")
	viper.AutomaticEnv()
}

func getDlframeworkHandler() (http.Handler, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	api := operations.NewDlframeworkAPI(swaggerSpec)
	handler := dlframework.ConfigureAPI(api)

	return handler, nil
}

func AllowTracedHeaders() echo.MiddlewareFunc {
	headers := strings.Join(
		[]string{
			"Origin",
			"Accept",
			"X-Requested-With",
			"X-B3-TraceId",
			"X-B3-ParentSpanId",
			"X-B3-SpanId",
			"X-B3-Sampled",
			"trace.traceid",
			"trace.spanid",
		},
		", ")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderAccessControlAllowHeaders)
			rid += headers
			res.Header().Set(echo.HeaderAccessControlAllowHeaders, rid)

			return next(c)
		}
	}
}
