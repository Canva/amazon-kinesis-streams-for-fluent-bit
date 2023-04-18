package eks

import (
	"testing"
	"time"

	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/mappings"
	"github.com/stretchr/testify/assert"
)

func TestValidNewEnricher(t *testing.T) {
	var cases = []struct {
		Name     string
		Env      map[string]string
		Expected *Enricher
	}{
		{
			Name: "Gets AccountId",
			Env: map[string]string{
				mappings.ENV_ACCOUNT_ID:    "1234567890",
				mappings.ENV_ACCOUNT_GROUP: DummyAccountFunction,
			},
			Expected: &Enricher{
				AccountId:            "1234567890",
				CanvaAccountFunction: DummyAccountFunction,
			},
		},
		{
			Name: "Gets Account Group",
			Env: map[string]string{
				mappings.ENV_ACCOUNT_ID:    DummyAccountId,
				mappings.ENV_ACCOUNT_GROUP: "PII",
			},
			Expected: &Enricher{
				AccountId:            DummyAccountId,
				CanvaAccountFunction: "PII",
			},
		},
	}

	for _, v := range cases {
		t.Run(v.Name, func(tt *testing.T) {
			for k, v := range v.Env {
				tt.Setenv(k, v)
			}
			actual, err := NewEnricher()

			assert.NoError(tt, err)

			assert.Equal(tt, v.Expected, actual)

			tt.Cleanup(func() {})
		})
	}
}

func TestInvalidNewEnricher(t *testing.T) {
	enricher, err := NewEnricher()

	assert.Nil(t, enricher)
	assert.Error(t, err)
}

func TestEnrichRecords(t *testing.T) {
	var cases = []struct {
		Name     string
		Enricher Enricher
		Input    map[interface{}]interface{}
		Expected map[interface{}]interface{}
	}{
		{
			Name: "Adds Account Id",
			Enricher: Enricher{
				AccountId:            "1234567",
				CanvaAccountFunction: DummyAccountFunction,
			},
			Input: map[interface{}]interface{}{
				"log": "hello world",
				"kubernetes": map[interface{}]interface{}{
					"key": "value",
				},
			},
			Expected: map[interface{}]interface{}{
				mappings.OBSERVED_TIMESTAMP: ExpectedTime,
				"log":                       "hello world",
				"resource": map[interface{}]interface{}{
					mappings.RESOURCE_CLOUD_ACCOUNT_ID: "1234567",
					mappings.RESOURCE_ACCOUNT_GROUP:    DummyAccountFunction,
				},
				"kubernetes": map[interface{}]interface{}{
					"key": "value",
				},
			},
		},
		{
			Name: "Adds Account Group Function",
			Enricher: Enricher{
				AccountId:            DummyAccountId,
				CanvaAccountFunction: "PII",
			},
			Input: map[interface{}]interface{}{
				"log": "hello world",
				"kubernetes": map[interface{}]interface{}{
					"key": "value",
				},
			},
			Expected: map[interface{}]interface{}{
				mappings.OBSERVED_TIMESTAMP: ExpectedTime,
				"log":                       "hello world",
				"resource": map[interface{}]interface{}{
					mappings.RESOURCE_CLOUD_ACCOUNT_ID: DummyAccountId,
					mappings.RESOURCE_ACCOUNT_GROUP:    "PII",
				},
				"kubernetes": map[interface{}]interface{}{
					"key": "value",
				},
			},
		},
		{
			Name: "Drop Log If Log Field Is Empty",
			Enricher: Enricher{
				AccountId:            DummyAccountId,
				CanvaAccountFunction: DummyAccountFunction,
			},
			Input: map[interface{}]interface{}{
				"observedTimestamp": DummyTime,
				"kubernetes": map[interface{}]interface{}{
					"key": "value",
				},
			},
			Expected: nil,
		},
		{
			Name: "Enrich placeholder service name if kubernetes dict is empty",
			Enricher: Enricher{
				AccountId:            DummyAccountId,
				CanvaAccountFunction: DummyAccountFunction,
			},
			Input: map[interface{}]interface{}{
				"log": "hello world",
			},
			Expected: map[interface{}]interface{}{
				mappings.OBSERVED_TIMESTAMP: ExpectedTime,
				"log":                       "hello world",
				"resource": map[interface{}]interface{}{
					mappings.RESOURCE_CLOUD_ACCOUNT_ID: DummyAccountId,
					mappings.RESOURCE_ACCOUNT_GROUP:    DummyAccountFunction,
				},
				"kubernetes": map[interface{}]interface{}{
					mappings.KUBERNETES_CONTAINER_NAME: mappings.PLACEHOLDER_MISSING_KUBERNETES_METADATA,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(tt *testing.T) {
			actual := c.Enricher.EnrichRecord(c.Input, DummyTime)
			assert.Equal(tt, c.Expected, actual)
		})
	}
}

var (
	DummyAccountFunction = "general"
	DummyAccountId       = "Account_Id"
	DummyTime            = time.Date(2009, time.November, 10, 23, 7, 5, 432000000, time.UTC)
)

var (
	ExpectedTime = int64(1257894425432)
)
