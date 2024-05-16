package mappings

const (
	OBSERVED_TIMESTAMP = "observedTimestamp"
)

const (
	LOG_FIELD_NAME            = "log"
	MESSAGE_FIELD_NAME        = "message"
	TRANSPORT_FIELD_NAME      = "transport"
	EKS_HOST_LOG_SERVICE_NAME = "eks_host_log"
)

const (
	KUBERNETES_RESOURCE_FIELD_NAME          = "kubernetes"
	KUBERNETES_LABELS_FIELD_NAME            = "labels"
	KUBERNETES_LABELS_NAME                  = "app.kubernetes.io/name"
	KUBERNETES_RESOURCE_CLUSTER_NAME        = "cluster.name"
	KUBERNETES_CONTAINER_NAME               = "container_name"
	KUBERNETES_RESOURCE_NODE_NAME           = "node.name"
	PLACEHOLDER_MISSING_KUBERNETES_METADATA = "_missing_metadata"
)

const (
	RESOURCE_APPLICATION_ID         = "application_id"
	RESOURCE_FIELD_NAME             = "resource"
	RESOURCE_SERVICE_NAME           = "service.name"
	RESOURCE_PARTITION              = "cloud.partition"
	RESOURCE_ACCOUNT_ID             = "cloud.account.id"
	RESOURCE_ACCOUNT_NAME           = "cloud.account.name"
	RESOURCE_REGION                 = "cloud.region"
	RESOURCE_ACCOUNT_GROUP_FUNCTION = "cloud.account.function"
	RESOURCE_ORGANIZATION           = "organization"
	RESOURCE_PLATFORM               = "cloud.platform"
	RESOURCE_PROVIDER               = "cloud.provider"
)

const (
	ENV_ACCOUNT_ID             = "CLOUD_ACCOUNT_ID"
	ENV_ACCOUNT_NAME           = "CLOUD_ACCOUNT_NAME"
	ENV_REGION                 = "CLOUD_REGION"
	ENV_ACCOUNT_GROUP_FUNCTION = "CLOUD_ACCOUNT_GROUP_FUNCTION"
	ENV_CLUSTER_NAME           = "K8S_CLUSTER_NAME"
	ENV_NODE_NAME              = "K8S_NODE_NAME"
	ENV_PARTITION              = "CLOUD_PARTITION"
	ENV_ORGANISATION           = "ORGANIZATION"
	ENV_PLATFORM               = "CLOUD_PLATFORM"
	ENV_PROVIDER               = "CLOUD_PROVIDER"
)
