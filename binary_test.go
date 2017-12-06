package hessian

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestBinary(t *testing.T) {
	buf := new(bytes.Buffer)

	var nilBytes []byte = nil
	if _, err := WriteBytes(buf, nilBytes); err != nil {
		t.Fatal(err)
	}
	if b, err := ReadBytes(buf); err != nil {
		t.Fatal(err)
	} else if b != nil {
		t.Fatal("encode nil bytes failed!")
	}

	compactBytes := []byte("hello world!")

	if _, err := WriteBytes(buf, compactBytes); err != nil {
		t.Fatal(err)
	}
	if b, err := ReadBytes(buf); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(compactBytes, b) {
		t.Fatal("encode short bytes failed!")
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))

	shortBytes := make([]byte, r.Intn(maxShortBinSize+1))
	for i := 0; i < len(shortBytes); i++ {
		shortBytes[i] = byte(r.Intn(0xFF + 1))
	}
	if _, err := WriteBytes(buf, shortBytes); err != nil {
		t.Fatal(err)
	}
	if b, err := ReadBytes(buf); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(shortBytes, b) {
		t.Fatal("encode short size bytes failed!")
	}

	chunkSizeBytes := make([]byte, maxBinChunkSize)
	for i := 0; i < len(chunkSizeBytes); i++ {
		chunkSizeBytes[i] = byte(r.Intn(0xFF + 1))
	}
	if _, err := WriteBytes(buf, chunkSizeBytes); err != nil {
		t.Fatal(err)
	}
	if b, err := ReadBytes(buf); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(chunkSizeBytes, b) {
		t.Fatal("encode chunk size bytes failed!")
	}

	ordBytes := make([]byte, maxBinChunkSize*r.Intn(10)+r.Intn(maxBinChunkSize))
	for i := 0; i < len(ordBytes); i++ {
		ordBytes[i] = byte(r.Intn(0xFF + 1))
	}
	if _, err := WriteBytes(buf, ordBytes); err != nil {
		t.Fatal(err)
	}
	if b, err := ReadBytes(buf); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(ordBytes, b) {
		t.Fatal("encode ordinary size bytes failed!")
	}

}
