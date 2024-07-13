package trace

import (
	"github.com/huangchunlong818/go-rabbitmq/rabbitmq/config"
	traces "github.com/huangchunlong818/go-trace/trace"
)

// 获取trace
func GetTraces(config config.TraceConfig) (*traces.TracerSpan, error) {
	traceObject, err := traces.GetNewTracerSpan(traces.TraceConfig{
		Version:          config.Version,
		Endpoint:         config.Endpoint,
		Project:          config.Project,
		Instance:         config.Instance,
		AccessKeyId:      config.AccessKeyId,
		AccessKeySecret:  config.AccessKeySecret,
		ServiceName:      config.ServiceName,
		ServiceNamespace: config.ServiceNamespace,
	})
	if err != nil {
		return nil, err
	}
	return traceObject, nil
}
