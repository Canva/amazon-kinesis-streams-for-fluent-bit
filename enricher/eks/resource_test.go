package eks

import (
	"reflect"
	"testing"
)

func Test_newResource(t *testing.T) {
	const accountID = "testAccountID"
	want := resource{
		"cloud.platform":   "aws_eks",
		"cloud.account.id": accountID,
	}
	got := newResource(accountID)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("newResource() = %v, want %v", got, want)
	}
}

func Test_resource_extractFromKubeTags(t *testing.T) {
	type args struct {
		tags interface{}
	}
	tests := []struct {
		name         string
		res          resource
		args         args
		wantResource resource
	}{
		{
			"service_pod",
			resource{},
			args{
				map[interface{}]interface{}{
					"host":            "test_host_id",
					"docker_id":       "test_container_id",
					"container_name":  "test_container_name",
					"namespace_name":  "test_namespace_name",
					"pod_name":        "test_pod_name",
					"pod_id":          "test_pod_id",
					"container_hash":  "test_container_hash",
					"container_image": "test_container_image",
				},
			},
			resource{
				"host.id":            "test_host_id",
				"container.id":       "test_container_id",
				"container.name":     "test_container_name",
				"k8s.namespace.name": "test_namespace_name",
				"k8s.pod.name":       "test_pod_name",
				"k8s.pod.uid":        "test_pod_id",
			},
		},
		{
			"invalid_type",
			resource{
				"key_1": "val_1",
			},
			args{
				"string",
			},
			resource{
				"key_1": "val_1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.res.extractFromKubeTags(tt.args.tags)
			if !reflect.DeepEqual(tt.res, tt.wantResource) {
				t.Errorf("EKS.extractFromKubeTags().resource = %v, want %v", tt.res, tt.wantResource)
			}
		})
	}
}

func Test_resource_extractFromLabelsTags(t *testing.T) {
	type args struct {
		tags interface{}
	}
	tests := []struct {
		name         string
		res          resource
		args         args
		wantResource resource
	}{
		{
			"service_pod",
			resource{},
			args{
				map[interface{}]interface{}{
					"canva.component":      "test_canva_component",
					"app":                  "deploy-exploration-worker",
					"canva.flavor":         "dev",
					"canva.platform":       "eks",
					"canva.application-id": "deploy-exploration",
					"name":                 "deploy-exploration-worker",
					"deploymentstack":      "new",
				},
			},
			resource{
				"service.name": "test_canva_component",
			},
		},
		{
			"invalid_type",
			resource{
				"key_1": "val_1",
			},
			args{
				"string",
			},
			resource{
				"key_1": "val_1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.res.extractFromLabelsTags(tt.args.tags)
			if !reflect.DeepEqual(tt.res, tt.wantResource) {
				t.Errorf("EKS.extractFromLabelsTags().resource = %v, want %v", tt.res, tt.wantResource)
			}
		})
	}
}
