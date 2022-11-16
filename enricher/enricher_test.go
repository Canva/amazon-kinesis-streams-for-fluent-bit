package enricher

import (
	"reflect"
	"testing"
	"time"

	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/ecs"
	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/eks"
	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/enricher/noop"
)

func IsType[T Enricher](c Enricher) bool {
	_, ok := c.(T)
	return ok
}

func TestNew(t *testing.T) {
	type args struct {
		conf *Config
	}
	tests := []struct {
		name          string
		args          args
		typeCheckFunc func(Enricher) bool
		wantErr       bool
	}{
		{
			"noop",
			args{
				&Config{
					Environment: "noop",
				},
			},
			IsType[noop.Noop],
			false,
		},
		{
			"gzip",
			args{
				&Config{
					Environment: "ecs",
				},
			},
			IsType[*ecs.ECS],
			false,
		},
		{
			"zstd",
			args{
				&Config{
					Environment: "eks",
				},
			},
			IsType[*eks.EKS],
			false,
		},
		{
			"invalid",
			args{
				&Config{
					Environment: "invalid",
				},
			},
			IsType[Enricher],
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if got != nil {
					t.Errorf("New() = %v, want nil", got)
				}
			} else {
				if !tt.typeCheckFunc(got) {
					t.Errorf("wrong New().(type) = %T", got)
				}
			}
		})
	}
}

func TestInit(t *testing.T) {
	defaultEnricher = nil
	Init(&Config{
		Environment: EnvironmentNoop,
	})
	if defaultEnricher == nil {
		t.Error("Init(); defaultEnricher = nil")
	}
}

func TestEnrichRecord(t *testing.T) {
	defaultEnricher = noop.Noop{}
	want := map[interface{}]interface{}{
		"key_1": "value_1",
		"key_2": "value_2",
	}
	got := EnrichRecord(
		map[interface{}]interface{}{
			"key_1": "value_1",
			"key_2": "value_2",
		},
		time.Date(2009, time.November, 10, 23, 7, 5, 432000000, time.UTC),
	)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Noop.EnrichRecord() = %v, want %v", got, want)
	}
}
