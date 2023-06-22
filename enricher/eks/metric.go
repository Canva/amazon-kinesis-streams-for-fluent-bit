package eks

import (
	"context"
	"fmt"

	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/mappings"
	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/metricserver"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type EnricherMetric struct {
	meter metric.Meter

	outputRecordCount  metric.Int64Counter
	outputDroppedCount metric.Int64Counter
}

func WithMetricServer(ms *metricserver.MetricServer) EnricherConfiguration {
	return func(e *Enricher) error {
		meter := ms.GetMeter("github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/eks")
		outputRecordCount, err := meter.Int64Counter("fluentbit.output.record", metric.WithDescription("output record counter"))

		if err != nil {
			return err
		}

		outputDroppedCount, err := meter.Int64Counter("fluentbit.output.dropped.body", metric.WithDescription("output dropped record counter where log or message does not exist"))

		if err != nil {
			return err
		}

		e.metric = &EnricherMetric{
			meter:              ms.GetMeter("github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/eks"),
			outputRecordCount:  outputRecordCount,
			outputDroppedCount: outputDroppedCount,
		}

		return nil
	}
}

func (e *Enricher) AddRecordCount(record map[interface{}]interface{}, logType LogType) {
	if e.metric == nil {
		return
	}

	var serviceName = inferServiceName(record, logType)

	e.metric.outputRecordCount.Add(context.Background(), 1, metric.WithAttributes(attribute.Key("source").String(serviceName)))
}

func (e *Enricher) AddDropCount() {
	if e.metric == nil {
		return
	}

	e.metric.outputDroppedCount.Add(context.Background(), 1)
}

func inferServiceName(record map[interface{}]interface{}, logType LogType) string {
	// fallback in case unable to find service name.
	var serviceName interface{} = "_missing_service_name"

	k8sPayload, ok := record[mappings.KUBERNETES_RESOURCE_FIELD_NAME].(map[interface{}]interface{})
	// k8s field exists, set default to container name.
	if ok {
		serviceName = k8sPayload[mappings.KUBERNETES_CONTAINER_NAME]
	}

	switch logType {
	case TYPE_APPLICATION:
		if labels, ok := k8sPayload[mappings.KUBERNETES_LABELS_FIELD_NAME].(map[interface{}]interface{}); ok {
			if val, ok := labels[mappings.KUBERNETES_LABELS_NAME]; ok {
				serviceName = val
			}
		}
	case TYPE_HOST:
		if resource, ok := record[mappings.RESOURCE_FIELD_NAME].(map[interface{}]interface{}); ok {
			if val, ok := resource[mappings.RESOURCE_SERVICE_NAME]; ok {
				serviceName = val
			}
		}
	}

	return fmt.Sprintf("%s", serviceName)
}
