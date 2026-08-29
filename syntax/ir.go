package syntax

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
)

var programIRMagic = []byte{'R', '2', 'G', 'O', 'I', 'R', 1}

func init() {
	for _, value := range []any{&Literal{}, &Symbol{}, &Call{}, &Block{}, &Function{}, &If{}, &While{}, &For{}, &Repeat{}} {
		gob.Register(value)
	}
}

// EncodeProgramIR serializes the typed syntax graph, including source maps,
// into a versioned compressed build artifact. It never stores R source text.
func EncodeProgramIR(program *Program) ([]byte, error) {
	var payload bytes.Buffer
	payload.Write(programIRMagic)
	zw, err := gzip.NewWriterLevel(&payload, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	if err = gob.NewEncoder(zw).Encode(program); err != nil {
		_ = zw.Close()
		return nil, err
	}
	if err = zw.Close(); err != nil {
		return nil, err
	}
	return payload.Bytes(), nil
}

func DecodeProgramIR(data []byte) (*Program, error) {
	if len(data) < len(programIRMagic) || !bytes.Equal(data[:len(programIRMagic)], programIRMagic) {
		return nil, fmt.Errorf("invalid or unsupported R2Go IR")
	}
	zr, err := gzip.NewReader(bytes.NewReader(data[len(programIRMagic):]))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	program := &Program{}
	if err := gob.NewDecoder(io.LimitReader(zr, 512<<20)).Decode(program); err != nil {
		return nil, err
	}
	return program, nil
}
