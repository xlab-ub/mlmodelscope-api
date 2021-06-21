package cmd

import (
	"fmt"
	dlframework "github.com/c3sr/dlframework/http"
	"github.com/c3sr/dlframework/httpapi/restapi"
	"github.com/c3sr/dlframework/httpapi/restapi/operations"
	"github.com/go-openapi/loads"
	"github.com/spf13/viper"
	"net/http"
	"os"

	"github.com/c3sr/config"
	"github.com/spf13/cobra"
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

	handler, err := getDlframeworkHandler()
	if err != nil {
		panic(fmt.Errorf("unable to create handler: %s", err))
	}

	http.Handle("/", handler)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(fmt.Errorf("unable to listen: %s", err))
	}
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
	viper.BindEnv("registry.provider", "registry_provider")
	viper.BindEnv("registry.endpoints", "registry_endpoints")
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
