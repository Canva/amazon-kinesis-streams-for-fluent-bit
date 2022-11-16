package eks

// resource holds and operates over resource values.
type resource map[interface{}]interface{}

// newResource creates a new resource with initial values.
func newResource(awsAccount string) resource {
	return map[interface{}]interface{}{
		"cloud.platform":   "aws_eks",
		"cloud.account.id": awsAccount,
	}
}

// extractFromKubeTags extracts kubernetes tags values and add them to resource.
func (res resource) extractFromKubeTags(tags interface{}) {
	tagsMap, ok := tags.(map[interface{}]interface{})
	if !ok {
		return
	}

	var strKey string
	for k, v := range tagsMap {
		strKey, ok = k.(string)
		if ok {
			switch strKey {
			case "host":
				res["host.id"] = v
			case "docker_id":
				res["container.id"] = v
			case "container_name":
				res["container.name"] = v
			case "namespace_name":
				res["k8s.namespace.name"] = v
			case "pod_name":
				res["k8s.pod.name"] = v
			case "pod_id":
				res["k8s.pod.uid"] = v
			default:
				// Skip
				// "container_hash":  "k8s.gcr.io/etcd@sha256:13f53ed1d91e2e11aac476ee9a0269fdda6cc4874eba903efd40daf50c55eee5"
				// "container_image": "k8s.gcr.io/etcd:3.5.3-0"
			}
		}
	}
}

// extractFromLabelsTags extracts labels tags values and add them to resource.
func (res resource) extractFromLabelsTags(tags interface{}) {
	tagsMap, ok := tags.(map[interface{}]interface{})
	if !ok {
		return
	}

	var strKey string
	for k, v := range tagsMap {
		strKey, ok = k.(string)
		if ok {
			switch strKey {
			case "canva.component":
				res["service.name"] = v
			default:
				// Skip
			}
		}
	}
}
