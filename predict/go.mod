module predict

go 1.15

replace (
	github.com/c3sr/tracer => ../../tracer
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)

require (
	github.com/c3sr/config v1.0.1
	github.com/c3sr/dlframework v1.3.2
	github.com/go-openapi/loads v0.20.2
	github.com/jaegertracing/jaeger-client-go v2.29.1+incompatible // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
)
