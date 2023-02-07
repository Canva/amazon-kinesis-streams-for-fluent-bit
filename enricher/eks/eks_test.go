package eks_test

import (
	"testing"
	"time"

	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/eks"
	"github.com/maxatome/go-testdeep/td"
)

func TestEnrichRecordsWithAccountId(t *testing.T) {
	var cases = []struct {
		Name      string
		AccountId int64
		Input     map[interface{}]interface{}
		Expected  map[interface{}]interface{}
	}{
		{
			Name:      "Adds Account Id",
			AccountId: 1234567,
			Input: map[interface{}]interface{}{
				"log": "hello world",
			},
			Expected: map[interface{}]interface{}{
				"log": "hello world",
				"resource": map[interface{}]interface{}{
					"account_id": int64(1234567),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			enr := eks.NewEnricher(c.AccountId)

			actual := enr.EnrichRecord(c.Input, DummyTime)
			td.Cmp(t, actual, c.Expected)
		})
	}
}

var DummyTime = time.Now()
