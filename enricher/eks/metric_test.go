package eks

import (
	"testing"

	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/mappings"
	"github.com/stretchr/testify/assert"
)

func Test_InferServiceName(t *testing.T) {
	type TestCase struct {
		Test     string
		Input    map[interface{}]interface{}
		Expected string
		LogType  LogType
	}

	testCases := []TestCase{
		{
			Test: "Application",
			Input: map[interface{}]interface{}{
				mappings.KUBERNETES_RESOURCE_FIELD_NAME: map[interface{}]interface{}{
					mappings.KUBERNETES_CONTAINER_NAME: "container name",
					mappings.KUBERNETES_LABELS_FIELD_NAME: map[interface{}]interface{}{
						mappings.KUBERNETES_LABELS_NAME: "rpc",
					},
				},
			},
			Expected: "rpc",
			LogType:  TYPE_APPLICATION,
		},
		{
			Test: "Application Non String",
			Input: map[interface{}]interface{}{
				mappings.KUBERNETES_RESOURCE_FIELD_NAME: map[interface{}]interface{}{
					mappings.KUBERNETES_CONTAINER_NAME: "container name",
					mappings.KUBERNETES_LABELS_FIELD_NAME: map[interface{}]interface{}{
						mappings.KUBERNETES_LABELS_NAME: []byte("ABC"),
					},
				},
			},
			Expected: "ABC",
			LogType:  TYPE_APPLICATION,
		},
		{
			Test: "Container Name",
			Input: map[interface{}]interface{}{
				mappings.KUBERNETES_RESOURCE_FIELD_NAME: map[interface{}]interface{}{
					mappings.KUBERNETES_CONTAINER_NAME: "container name",
				},
			},
			Expected: "container name",
			LogType:  TYPE_EMPTY,
		},
		{
			Test: "Container Name Non String",
			Input: map[interface{}]interface{}{
				mappings.KUBERNETES_RESOURCE_FIELD_NAME: map[interface{}]interface{}{
					mappings.KUBERNETES_CONTAINER_NAME: []byte("ABC"),
				},
			},
			Expected: "ABC",
			LogType:  TYPE_EMPTY,
		},
		{
			Test: "Host",
			Input: map[interface{}]interface{}{
				mappings.RESOURCE_FIELD_NAME: map[interface{}]interface{}{
					mappings.RESOURCE_SERVICE_NAME: mappings.EKS_HOST_LOG_SERVICE_NAME,
				},
			},
			Expected: mappings.EKS_HOST_LOG_SERVICE_NAME,
			LogType:  TYPE_HOST,
		},
		{
			Test:     "Invalid",
			Input:    map[interface{}]interface{}{},
			Expected: "_missing",
			LogType:  TYPE_APPLICATION,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			actual := inferServiceName(tc.Input, tc.LogType)

			assert.Equal(t, tc.Expected, actual)
		})
	}
}
