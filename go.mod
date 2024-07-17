module github.com/huangchunlong818/go-rabbitmq

go 1.20

replace go.opentelemetry.io/otel v1.28.0 => go.opentelemetry.io/otel v1.24.0

replace go.opentelemetry.io/otel/trace v1.28.0 => go.opentelemetry.io/otel/trace v1.24.0

replace go.opentelemetry.io/otel/sdk v1.28.0 => go.opentelemetry.io/otel/sdk v1.24.0

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.28.0 => go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0 => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.28.0 => go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.21.0

replace go.opentelemetry.io/otel/sdk/metric v1.28.0 => go.opentelemetry.io/otel/sdk/metric v1.24.0

replace go.opentelemetry.io/otel/metric v1.28.0 => go.opentelemetry.io/otel/metric v1.24.0

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.28.0 => go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.44.0

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.28.0 => go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.44.0

replace go.opentelemetry.io/proto/otlp v1.3.1 => go.opentelemetry.io/proto/otlp v1.0.0

replace go.opentelemetry.io/contrib/instrumentation/host v0.53.0 => go.opentelemetry.io/contrib/instrumentation/host v0.46.1

replace go.opentelemetry.io/contrib/instrumentation/runtime v0.53.0 => go.opentelemetry.io/contrib/instrumentation/runtime v0.46.1

replace google.golang.org/grpc v1.65.0 => google.golang.org/grpc v1.63.2

replace github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0 => github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1

require (
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/huangchunlong818/go-trace v0.1.1
	github.com/rs/xid v1.5.0
	github.com/streadway/amqp v1.1.0
	go.opentelemetry.io/otel/trace v1.24.0
)

require (
	github.com/aliyun-sls/opentelemetry-go-provider-sls v0.12.0 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20231016141302-07b5767bb0ed // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/sethvargo/go-envconfig v0.9.0 // indirect
	github.com/shirou/gopsutil/v3 v3.23.11 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/tklauser/go-sysconf v0.3.13 // indirect
	github.com/tklauser/numcpus v0.7.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.opentelemetry.io/contrib/instrumentation/host v0.46.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.46.1 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk v1.24.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.24.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)
