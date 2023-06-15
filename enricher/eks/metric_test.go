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
		},
		{
			Test: "Application Non String",
			Input: map[interface{}]interface{}{
				mappings.KUBERNETES_RESOURCE_FIELD_NAME: map[interface{}]interface{}{
					mappings.KUBERNETES_CONTAINER_NAME: "container name",
					mappings.KUBERNETES_LABELS_FIELD_NAME: map[interface{}]interface{}{
						mappings.KUBERNETES_LABELS_NAME: -123,
					},
				},
			},
			Expected: "-123",
		},
		{
			Test: "Container Name",
			Input: map[interface{}]interface{}{
				mappings.KUBERNETES_RESOURCE_FIELD_NAME: map[interface{}]interface{}{
					mappings.KUBERNETES_CONTAINER_NAME: "container name",
				},
			},
			Expected: "container name",
		},
		{
			Test: "Container Name Non String",
			Input: map[interface{}]interface{}{
				mappings.KUBERNETES_RESOURCE_FIELD_NAME: map[interface{}]interface{}{
					mappings.KUBERNETES_CONTAINER_NAME: -123,
				},
			},
			Expected: "-123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			actual := inferServiceName(tc.Input)

			assert.Equal(t, tc.Expected, actual)
		})
	}
}
