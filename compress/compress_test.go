package compress

import (
	"reflect"
	"testing"

	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/compress/gzip"
	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/compress/noop"
	"github.com/canva/amazon-kinesis-streams-for-fluent-bit/compress/zstd"
)

func IsType[T Compression](c Compression) bool {
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
		typeCheckFunc func(Compression) bool
		wantErr       bool
	}{
		{
			"noop",
			args{
				&Config{
					Format: "noop",
				},
			},
			IsType[noop.Noop],
			false,
		},
		{
			"gzip",
			args{
				&Config{
					Format: "gzip",
					Level:  1,
				},
			},
			IsType[*gzip.GZip],
			false,
		},
		{
			"zstd",
			args{
				&Config{
					Format: "zstd",
					Level:  1,
				},
			},
			IsType[*zstd.ZSTD],
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.typeCheckFunc(got) {
				t.Errorf("wrong New().(type) = %T", got)
			}
		})
	}
}

func TestInit(t *testing.T) {
	defaultCompression = nil
	Init(&Config{
		Format: FormatNoop,
	})
	if defaultCompression == nil {
		t.Error("Init(); defaultCompression = nil")
	}
}

func TestCompress(t *testing.T) {
	defaultCompression = noop.Noop{}

	want := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	got, err := Compress([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	if err != nil {
		t.Errorf("Noop.Compress() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Noop.Compress() = %v, want %v", got, want)
	}
}
