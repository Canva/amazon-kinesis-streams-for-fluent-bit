package eks

import (
	"time"

	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher"
)

type Enricher struct {
	accountId int64
}

func NewEnricher(accountId int64) enricher.IEnricher {
	return &Enricher{
		accountId: accountId,
	}
}

var _ enricher.IEnricher = (*Enricher)(nil)

func (e Enricher) EnrichRecord(r map[interface{}]interface{}, _ time.Time) map[interface{}]interface{} {
	// add resource attributes
	r["resource"] = map[interface{}]interface{}{
		"account_id": e.accountId,
	}

	return r
}
