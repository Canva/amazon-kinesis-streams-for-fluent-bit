package ecs

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	want := &ECS{
		canvaAWSAccount: "test_account",
		canvaAppName:    "test_app_name",
		logGroup:        "test_log_group",
		ecsTaskFamily:   "test_ecs_task_family",
		ecsTaskRevision: 10001,
	}
	t.Setenv("CANVA_AWS_ACCOUNT", "test_account")
	t.Setenv("CANVA_APP_NAME", "test_app_name")
	t.Setenv("LOG_GROUP", "test_log_group")
	t.Setenv("ECS_TASK_DEFINITION", "test_ecs_task_family:10001")
	if got := New(); !reflect.DeepEqual(got, want) {
		t.Errorf("New() = %v, want %v", got, want)
	}
}

func TestECS_EnrichRecord(t *testing.T) {
	type fields struct {
		canvaAWSAccount string
		canvaAppName    string
		logGroup        string
		ecsTaskFamily   string
		ecsTaskRevision int
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
			"all_values",
			fields{
				"canva_aws_account_val",
				"canva_app_name_val",
				"log_group_val",
				"ecs_task_family_val",
				10001,
			},
			args{
				map[interface{}]interface{}{
					"ec2_instance_id":     "ec2_instance_id_val",
					"ecs_cluster":         "ecs_cluster_val",
					"ecs_task_arn":        "ecs_task_arn_val",
					"container_id":        "container_id_val",
					"container_name":      "container_name_val",
					"other_key_1":         "other_value_1",
					"other_key_2":         "other_value_2",
					"other_key_3":         "other_value_3",
					"timestamp":           "1234567890",
					"ecs_task_definition": "ecs_task_definition_val",
				},
				time.Date(2009, time.November, 10, 23, 7, 5, 432000000, time.UTC),
			},
			map[interface{}]interface{}{
				"resource": map[interface{}]interface{}{
					"cloud.account.id":      "canva_aws_account_val",
					"service.name":          "canva_app_name_val",
					"cloud.platform":        "aws_ecs",
					"aws.ecs.launchtype":    "EC2",
					"aws.ecs.task.family":   "ecs_task_family_val",
					"aws.ecs.task.revision": 10001,
					"aws.log.group.names":   "log_group_val",
					"host.id":               "ec2_instance_id_val",
					"aws.ecs.cluster.name":  "ecs_cluster_val",
					"aws.ecs.task.arn":      "ecs_task_arn_val",
					"container.id":          "container_id_val",
					"container.name":        "container_name_val",
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
			enr := &ECS{
				canvaAWSAccount: tt.fields.canvaAWSAccount,
				canvaAppName:    tt.fields.canvaAppName,
				logGroup:        tt.fields.logGroup,
				ecsTaskFamily:   tt.fields.ecsTaskFamily,
				ecsTaskRevision: tt.fields.ecsTaskRevision,
			}
			if got := enr.EnrichRecord(tt.args.r, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ECS.EnrichRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}
