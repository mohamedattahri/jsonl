package jsonl

import (
	"encoding/json"
	"io"
)

// Writer returns a new JSONLines
// writer.
type Writer[T any] struct {
	// Line separator,
	// Defaults to '\n'
	EOL string

	dest    io.Writer
	loc     int
	written int
}

// Written returns the total number of bytes
// written
func (w *Writer[T]) Written() int {
	return w.written
}

func (w *Writer[T]) write(v T) (int, error) {
	var (
		raw []byte
		n   int
		nn  int
		err error
	)

	defer func() {
		w.written += nn
		if err == nil {
			w.loc++
		}
	}()

	raw, err = json.Marshal(v)
	if err != nil {
		return n, err
	}

	if w.loc > 0 {
		n, err = io.WriteString(w.dest, w.EOL)
		nn += n
		if err != nil {
			return n, err
		}
	}

	n, err = w.dest.Write(raw)
	nn += n
	if err != nil {
		return nn, err
	}

	return nn, nil
}

// Write writes the JSON representation of every item
// in a new line.
func (w *Writer[T]) Write(v ...T) (n int, err error) {
	var nn int
	for _, item := range v {
		nn, err = w.write(item)
		n += nn
		if err != nil {
			return
		}
	}
	return
}

// NewWriter returns a new writer.
func NewWriter[T any](dest io.Writer) *Writer[T] {
	return &Writer[T]{
		dest: dest,
		EOL:  "\n",
	}
}
