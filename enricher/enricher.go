package enricher

import (
	"errors"
	"time"

	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/ecs"
	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/eks"
	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/noop"
)

// Enricher modifies records by enriching them.
type Enricher interface {
	EnrichRecord(r map[interface{}]interface{}, t time.Time) map[interface{}]interface{}
}

// Environment is the environment that enricher works in.
type Environment string

// Environment is the environment that enricher works in.
const (
	EnvironmentNoop = Environment("noop")
	EnvironmentECS  = Environment("ecs")
	EnvironmentEKS  = Environment("eks")
)

// Config holds configurations need for creating a new Enricher instance.
type Config struct {
	Environment Environment
}

// New creates an instance of Enricher.
// Based on the Environment value, it creates corresponding instance.
func New(conf *Config) (Enricher, error) {
	switch conf.Environment {
	default:
		return nil, errors.New("enricher environment not found")
	case "", EnvironmentNoop:
		return noop.New(), nil
	case EnvironmentECS:
		return ecs.New(), nil
	case EnvironmentEKS:
		return eks.New(), nil
	}
}

// defaultEnricher holds global Enricher instance.
var defaultEnricher Enricher

// Init will initialise global Enricher instance.
func Init(conf *Config) error {
	var err error

	defaultEnricher, err = New(conf)

	return err
}

// EnrichRecord modifies existing record.
func EnrichRecord(r map[interface{}]interface{}, t time.Time) map[interface{}]interface{} {
	return defaultEnricher.EnrichRecord(r, t)
}
