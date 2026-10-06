package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"
)

func decompress(t *testing.T, data []byte) []byte {
	t.Helper()
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("gzip.NewReader() error = %v", err)
	}
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}
	return got
}

func TestNew(t *testing.T) {
	z, err := New(5)

	if err != nil {
		t.Errorf("New() error = %v", err)
	}
	if z == nil {
		t.Errorf("New() = nil")
	}

	if z.level != 5 {
		t.Errorf("New().level = %v, want %v", z.level, 5)
	}

	if z.writer == nil {
		t.Errorf("New().writer = nil")
	}
}

func TestGZip_Compress(t *testing.T) {
	buf := bytes.Buffer{}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		z       *GZip
		args    args
		lock    bool
		wantErr bool
	}{
		{
			"not_locked",
			&GZip{
				writer: gzip.NewWriter(&buf),
			},
			args{[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
			false,
			false,
		},
		{
			"locked",
			&GZip{
				writer: gzip.NewWriter(&buf),
				level:  gzip.DefaultCompression,
			},
			args{[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.lock {
				tt.z.mutex.Lock()
			}
			got, err := tt.z.Compress(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GZip.Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got := decompress(t, got); !bytes.Equal(got, tt.args.data) {
				t.Errorf("GZip.Compress() decompressed = %v, want %v", got, tt.args.data)
			}
		})
	}
}

func TestBufferReuse(t *testing.T) {
	data1 := []byte{1, 2, 3}
	data2 := []byte{4, 5, 6, 7, 8}
	data3 := []byte{9, 10}

	z, err := New(1)
	if err != nil {
		t.Errorf("New() error = %v", err)
	}

	compressed1, err := z.Compress(data1)
	if err != nil {
		t.Errorf("New() error = %v", err)
	}
	compressed2, err := z.Compress(data2)
	if err != nil {
		t.Errorf("New() error = %v", err)
	}
	compressed3, err := z.Compress(data3)
	if err != nil {
		t.Errorf("New() error = %v", err)
	}

	if got := decompress(t, compressed1); !bytes.Equal(got, data1) {
		t.Errorf("compressed1 decompressed to %v, want %v", got, data1)
	}
	if got := decompress(t, compressed2); !bytes.Equal(got, data2) {
		t.Errorf("compressed2 decompressed to %v, want %v", got, data2)
	}
	if got := decompress(t, compressed3); !bytes.Equal(got, data3) {
		t.Errorf("compressed3 decompressed to %v, want %v", got, data3)
	}
}
