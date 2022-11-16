package eks

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	want := &EKS{
		canvaAWSAccount: "test_account",
	}
	t.Setenv("CANVA_AWS_ACCOUNT", "test_account")
	if got := New(); !reflect.DeepEqual(got, want) {
		t.Errorf("New() = %v, want %v", got, want)
	}
}

func TestEKS_EnrichRecord(t *testing.T) {
	type fields struct {
		canvaAWSAccount string
	}
	type args struct {
		r map[interface{}]interface{}
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[interface{}]interface{}
	}{
		{
			"service_pod",
			fields{
				"test_account",
			},
			args{
				map[interface{}]interface{}{
					"body": map[interface{}]interface{}{
						"other_key_1": "other_value_1",
						"other_key_2": "other_value_2",
						"other_key_3": "other_value_3",
						"timestamp":   "1234567890",
					},
					"kubernetes": map[interface{}]interface{}{
						"host":            "test_host_id",
						"docker_id":       "test_container_id",
						"container_name":  "test_container_name",
						"namespace_name":  "test_namespace_name",
						"pod_name":        "test_pod_name",
						"pod_id":          "test_pod_id",
						"container_hash":  "test_container_hash",
						"container_image": "test_container_image",
					},
					"labels": map[interface{}]interface{}{
						"canva.component":      "test_canva_component",
						"app":                  "deploy-exploration-worker",
						"canva.flavor":         "dev",
						"canva.platform":       "eks",
						"canva.application-id": "deploy-exploration",
						"name":                 "deploy-exploration-worker",
						"deploymentstack":      "new",
					},
					"stream": "stderr",
					"time":   "2022-11-15T05:30:18.573672299Z",
				},
				time.Date(2009, time.November, 10, 23, 7, 5, 432000000, time.UTC),
			},
			map[interface{}]interface{}{
				"resource": map[interface{}]interface{}{
					"cloud.platform":     "aws_eks",
					"cloud.account.id":   "test_account",
					"host.id":            "test_host_id",
					"container.id":       "test_container_id",
					"container.name":     "test_container_name",
					"k8s.namespace.name": "test_namespace_name",
					"k8s.pod.name":       "test_pod_name",
					"k8s.pod.uid":        "test_pod_id",
					"service.name":       "test_canva_component",
				},
				"body": map[interface{}]interface{}{
					"other_key_1": "other_value_1",
					"other_key_2": "other_value_2",
					"other_key_3": "other_value_3",
				},
				"timestamp":         "1234567890",
				"observedTimestamp": int64(1257894425432),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enr := &EKS{
				canvaAWSAccount: tt.fields.canvaAWSAccount,
			}
			if got := enr.EnrichRecord(tt.args.r, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EKS.EnrichRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractTimestampFromBody(t *testing.T) {
	type args struct {
		body interface{}
	}
	tests := []struct {
		name     string
		args     args
		want     interface{}
		wantBody interface{}
	}{
		{
			"not_found",
			args{
				map[interface{}]interface{}{
					"key_1": "val_1",
					"key_2": "val_2",
				},
			},
			nil,
			map[interface{}]interface{}{
				"key_1": "val_1",
				"key_2": "val_2",
			},
		},
		{
			"found",
			args{
				map[interface{}]interface{}{
					"key_1":     "val_1",
					"key_2":     "val_2",
					"timestamp": "1234567890",
				},
			},
			"1234567890",
			map[interface{}]interface{}{
				"key_1": "val_1",
				"key_2": "val_2",
			},
		},
		{
			"invalid_type",
			args{
				"string",
			},
			nil,
			"string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractTimestampFromBody(tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractTimestampFromBody() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.args.body, tt.wantBody) {
				t.Errorf("extractTimestampFromBody().body = %v, want %v", tt.args.body, tt.wantBody)
			}
		})
	}
}
