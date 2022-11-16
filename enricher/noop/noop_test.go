package noop

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Nothing to check, no panic is enough for this test.
	_ = New()
}

func TestNoop_EnrichRecord(t *testing.T) {
	type args struct {
		r   map[interface{}]interface{}
		in1 time.Time
	}
	tests := []struct {
		name string
		enr  Noop
		args args
		want map[interface{}]interface{}
	}{
		{
			"noop",
			Noop{},
			args{
				map[interface{}]interface{}{
					"k1": 1,
					2:    "v2",
				},
				time.Now(),
			},
			map[interface{}]interface{}{
				"k1": 1,
				2:    "v2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enr := Noop{}
			if got := enr.EnrichRecord(tt.args.r, tt.args.in1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Noop.EnrichRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}
