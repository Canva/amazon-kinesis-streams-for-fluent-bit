package eks

import (
	"testing"
	"time"

	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/mappings"
	"github.com/maxatome/go-testdeep/td"
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
				mappings.ENV_ACCOUNT_GROUP: DummyAccountGroup,
			},
			Expected: &Enricher{
				accountId:    "1234567890",
				accountGroup: DummyAccountGroup,
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

func TestEnrichRecordsWithAccountId(t *testing.T) {
	var cases = []struct {
		Name     string
		Enricher Enricher
		Input    map[interface{}]interface{}
		Expected map[interface{}]interface{}
	}{
		{
			Name: "Adds Account Id",
			Enricher: Enricher{
				accountId:    "1234567",
				accountGroup: DummyAccountGroup,
			},
			Input: map[interface{}]interface{}{
				"log": "hello world",
			},
			Expected: map[interface{}]interface{}{
				"log": "hello world",
				"resource": map[interface{}]interface{}{
					mappings.RESOURCE_CLOUD_ACCOUNT_ID: "1234567",
					mappings.RESOURCE_ACCOUNT_GROUP:    DummyAccountGroup,
				},
			},
		},
		{
			Name: "Adds Account Group",
			Enricher: Enricher{
				accountId:    DummyAccountId,
				accountGroup: "PII",
			},
			Input: map[interface{}]interface{}{
				"log": "hello world",
			},
			Expected: map[interface{}]interface{}{
				"log": "hello world",
				"resource": map[interface{}]interface{}{
					mappings.RESOURCE_CLOUD_ACCOUNT_ID: DummyAccountId,
					mappings.RESOURCE_ACCOUNT_GROUP:    "PII",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(tt *testing.T) {
			actual := c.Enricher.EnrichRecord(c.Input, DummyTime)
			td.Cmp(tt, actual, c.Expected)
		})
	}
}

var (
	DummyTime         = time.Now()
	DummyAccountGroup = "general"
	DummyAccountId    = "Account_Id"
)
