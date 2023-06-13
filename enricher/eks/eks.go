package eks

import (
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher"
	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/mappings"
)

// Log type
const (
	TYPE_EMPTY = iota + 1
	TYPE_APPLICATION
	TYPE_HOST
)

type EnricherConfiguration func(*Enricher) error

type Enricher struct {
	metric *EnricherMetric

	CloudAccountId            string `env:"CLOUD_ACCOUNT_ID,required"`
	CloudAccountName          string `env:"CLOUD_ACCOUNT_NAME,required"`
	CloudRegion               string `env:"CLOUD_REGION,required"`
	K8sClusterName            string `env:"K8S_CLUSTER_NAME,required"`
	K8sNodeName               string `env:"K8S_NODE_NAME,required"`
	CloudPartition            string `env:"CLOUD_PARTITION,required"`
	CloudAccountGroupFunction string `env:"CLOUD_ACCOUNT_GROUP_FUNCTION,required"`
	Organization              string `env:"ORGANIZATION,required"`
	CloudProvider             string `env:"CLOUD_PROVIDER,required"`
	CloudPlatform             string `env:"CLOUD_PLATFORM,required"`
}

// NewEnricher returns a enricher with env vars being parsed.
// These env vars are derived from mappings.go.
func NewEnricher(cfgs ...EnricherConfiguration) (*Enricher, error) {
	enricher := Enricher{}
	if err := env.Parse(&enricher); err != nil {
		return nil, err
	}

	return &enricher, nil
}

var _ enricher.IEnricher = (*Enricher)(nil)

func (e *Enricher) EnrichRecord(r map[interface{}]interface{}, t time.Time) map[interface{}]interface{} {
	logType := e.inferType(r)

	// Drop log if "log" field and "message" field is empty
	if logType == TYPE_EMPTY {
		// add drop metric
		e.AddDropCount()
		return nil
	}

	// Add static attributes
	r[mappings.RESOURCE_FIELD_NAME] = map[interface{}]interface{}{
		mappings.RESOURCE_ACCOUNT_ID:             e.CloudAccountId,
		mappings.RESOURCE_ACCOUNT_NAME:           e.CloudAccountName,
		mappings.RESOURCE_ACCOUNT_GROUP_FUNCTION: e.CloudAccountGroupFunction,
		mappings.RESOURCE_PARTITION:              e.CloudPartition,
		mappings.RESOURCE_REGION:                 e.CloudRegion,
		mappings.RESOURCE_ORGANIZATION:           e.Organization,
		mappings.RESOURCE_PLATFORM:               e.CloudPlatform,
		mappings.RESOURCE_PROVIDER:               e.CloudProvider,
	}
	r[mappings.OBSERVED_TIMESTAMP] = t.UnixMilli()

	switch logType {
	case TYPE_HOST:
		r[mappings.RESOURCE_FIELD_NAME].(map[interface{}]interface{})[mappings.RESOURCE_SERVICE_NAME] = mappings.EKS_HOST_LOG_SERVICE_NAME
		r[mappings.KUBERNETES_RESOURCE_FIELD_NAME] = map[interface{}]interface{}{}
	case TYPE_APPLICATION:
		if _, ok := r[mappings.KUBERNETES_RESOURCE_FIELD_NAME]; !ok {
			// The pod has started up and potentially died too early for us to enrich the log
			// We insert a placeholder value for the kubernetes.container_name
			// https://docs.google.com/document/d/1vRCUKMeo6ypnAq34iwQN7LtDsXxmlj0aYEfRofwV7A4/edit
			r[mappings.KUBERNETES_RESOURCE_FIELD_NAME] = map[interface{}]interface{}{
				mappings.KUBERNETES_CONTAINER_NAME: mappings.PLACEHOLDER_MISSING_KUBERNETES_METADATA,
			}
		}
	default:
		break
	}

	// Add kubernetes static attributes
	r[mappings.KUBERNETES_RESOURCE_FIELD_NAME].(map[interface{}]interface{})[mappings.KUBERNETES_RESOURCE_CLUSTER_NAME] = e.K8sClusterName
	r[mappings.KUBERNETES_RESOURCE_FIELD_NAME].(map[interface{}]interface{})[mappings.KUBERNETES_RESOURCE_NODE_NAME] = e.K8sNodeName

	// add metrics
	e.AddRecordCount(r, logType)

	return r
}

func (e Enricher) inferType(record map[interface{}]interface{}) int {
	_, hasApplicationLog := record[mappings.LOG_FIELD_NAME]
	_, hasHostMsgOk := record[mappings.MESSAGE_FIELD_NAME]
	_, hostTransportOk := record[mappings.TRANSPORT_FIELD_NAME]

	if hasApplicationLog {
		return TYPE_APPLICATION
	}

	if hasHostMsgOk && hostTransportOk {
		return TYPE_HOST
	}

	return TYPE_EMPTY
}
