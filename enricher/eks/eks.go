// Package eks is an implementation of Enricher interface using EKS environment.
package eks

import (
	"os"
	"time"
)

// EKS implements Enricher interface for EKS environment.
type EKS struct {
	canvaAWSAccount string
}

// New returns an EKS enricher.
func New() *EKS {
	return &EKS{
		canvaAWSAccount: os.Getenv("CANVA_AWS_ACCOUNT"),
	}
}

// EnrichRecord modifies existing record.
func (enr *EKS) EnrichRecord(r map[interface{}]interface{}, t time.Time) map[interface{}]interface{} {
	res := newResource(enr.canvaAWSAccount)

	var (
		ok     bool
		strKey string
		body   interface{}
	)

	for k, v := range r {
		strKey, ok = k.(string)
		if ok {
			switch strKey {
			case "body":
				body = v
			case "kubernetes":
				res.extractFromKubeTags(v)
			case "labels":
				res.extractFromLabelsTags(v)
			default:
				// Skip
				// "stream": "stderr"
				// "time":   "2022-11-15T05:30:18.573672299Z"
			}
		}
	}

	return map[interface{}]interface{}{
		"resource":          map[interface{}]interface{}(res),
		"body":              body,
		"timestamp":         extractTimestampFromBody(body),
		"observedTimestamp": t.UnixMilli(),
	}
}

// extractTimestampFromBody extracts and removes timestamp from log body.
func extractTimestampFromBody(body interface{}) interface{} {
	bodyMap, ok := body.(map[interface{}]interface{})
	if !ok {
		return nil
	}

	var strKey string
	for k, v := range bodyMap {
		strKey, ok = k.(string)
		if ok && strKey == "timestamp" {
			delete(bodyMap, k)
			return v
		}
	}

	return nil
}
