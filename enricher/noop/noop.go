// Package noop is an implementation of Enricher interface that doesn't modify records.
package noop

import "time"

// Noop is an implementation of Enricher interface that doesn't modify records.
type Noop struct{}

// New creates a new instance of Noop.
func New() Noop {
	return Noop{}
}

// EnrichRecord will return record without any modifying it.
func (enr Noop) EnrichRecord(r map[interface{}]interface{}, _ time.Time) map[interface{}]interface{} {
	return r
}
