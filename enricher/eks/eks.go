package eks

import (
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher"
	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/mappings"
)

type Enricher struct {
	accountId    string `env:"CANVA_AWS_ACCOUNT"`
	accountGroup string `env:"CANVA_AWS_ACCOUNT_GROUP"`
}

func NewEnricher() (*Enricher, error) {
	enricher := Enricher{}
	err := env.Parse(&enricher)

	return &enricher, err
}

var _ enricher.IEnricher = (*Enricher)(nil)

func (e Enricher) EnrichRecord(r map[interface{}]interface{}, _ time.Time) map[interface{}]interface{} {
	// add resource attributes
	r["resource"] = map[interface{}]interface{}{
		mappings.RESOURCE_CLOUD_ACCOUNT_ID: e.accountId,
		mappings.RESOURCE_ACCOUNT_GROUP:    e.accountGroup,
	}

	return r
}
