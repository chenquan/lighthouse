package binary

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestReadBool(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"true",
			args{r: bytes.NewReader([]byte{1})},
			true,
			false,
		}, {
			"false",
			args{r: bytes.NewReader([]byte{0})},
			false,
			false,
		}, {
			"error",
			args{r: bytes.NewReader([]byte{})},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadBool(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadBool() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type limitWrite struct {
}

func (l *limitWrite) Write(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func TestWriteUint16(t *testing.T) {
	b := &bytes.Buffer{}

	err := WriteUint16(b, 1)
	assert.NoError(t, err)
	assert.EqualValues(t, bytes.NewBuffer([]byte{0, 1}), b)

	err = WriteUint16(&limitWrite{}, 1)
	assert.Error(t, err)

}

func TestReadUint16(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{
			"correct", args{r: bytes.NewReader([]byte{0, 1})}, 1, false,
		}, {
			"error", args{r: bytes.NewReader([]byte{1})}, 0, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadUint16(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint16() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteBool(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			"1",
			args{b: true},
			bytes.NewBuffer([]byte{1}).String(), false,
		}, {
			"1",
			args{b: false},
			bytes.NewBuffer([]byte{0}).String(), false,
		}, {
			"1",
			args{b: false},
			bytes.NewBuffer([]byte{0}).String(), false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := WriteBool(w, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("WriteBool() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestReadUint32(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		wantErr bool
	}{
		{
			"correct",
			args{r: bytes.NewReader([]byte{0, 0, 0, 1})},
			1, false,
		}, {
			"error",
			args{r: bytes.NewReader([]byte{0, 0, 1})},
			0, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadUint32(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint32() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteUint32(t *testing.T) {
	err := WriteUint32(&limitWrite{}, 1)
	assert.Error(t, err)

	buffer := &bytes.Buffer{}
	err = WriteUint32(buffer, 1)
	assert.NoError(t, err)
	assert.EqualValues(t, bytes.NewBuffer([]byte{0, 0, 0, 1}), buffer)

}

func TestWriteString(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := WriteString(buffer, []byte("1"))
	assert.NoError(t, err)
	assert.EqualValues(t, bytes.NewBuffer([]byte{0, 1, '1'}), buffer)
	err = WriteString(&limitWrite{}, []byte(" "))
	assert.Error(t, err)
}

func TestReadString(t *testing.T) {
	readString, err := ReadString(bytes.NewBuffer([]byte{0, 1, '1'}))
	assert.NoError(t, err)
	assert.EqualValues(t, "1", readString)

	readString, err = ReadString(bytes.NewBuffer([]byte{0, 2, '1'}))
	assert.Error(t, err)
	readString, err = ReadString(bytes.NewBuffer([]byte{0}))
	assert.Error(t, err)

}
