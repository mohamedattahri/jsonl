package jsonl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// Line represents a JSON line
type Line json.RawMessage

// Scan decodes l into v using
// [json.Unmarshal].
func (l Line) Scan(v interface{}) error {
	return json.Unmarshal(l, v)
}

// String returns a string representation of
// l.
func (l Line) String() string {
	return string(l)
}

// Scanner can read JSON lines from an
// [io.Reader], one at a time.
type Scanner struct {
	// Line separator (End-Of-Line),
	// defaults to '\n'.
	EOL byte

	// Skip blank lines, or
	// trigger an error.
	SkipBlank bool

	// Lines starting with these
	// prefixes are ignored as comments.
	// Defaults to [].
	SkipComments []string

	rd   *bufio.Reader
	line Line
	loc  int
	err  error
}

// Line returns the latest line read.
func (s *Scanner) Line() (Line, error) {
	if s.loc == 0 {
		return nil, fmt.Errorf("jsonl: Next() must be called first")
	}
	if s.err != nil && s.err != io.EOF {
		return nil, s.err
	}
	return s.line, nil
}

// Next returns true if there's a line
// to read using [Line].
func (s *Scanner) Next() bool {
	if s.err != nil {
		return false
	}

	var (
		err error
		raw []byte
	)
	defer func() {
		if err != nil {
			s.err = err
		}
	}()

	for {
		raw, err = s.rd.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return false
		}
		s.loc++

		raw = bytes.TrimRightFunc(raw, unicode.IsSpace)
		if len(raw) == 0 && s.SkipBlank {
			continue
		}
		if hasPrefixAny(string(raw), s.SkipComments) {
			continue
		}

		s.line = Line(raw)
		break
	}

	return true
}

// NewScanner returns a way to read
// the JSON lines in src one at a time.
func NewScanner(src io.Reader) *Scanner {
	return &Scanner{
		EOL: '\n',
		rd:  bufio.NewReader(src),
	}
}

// hasPrefixAny returns true if raw is prefixed with
// one of the given prefixes
func hasPrefixAny(raw string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(raw, prefix) {
			return true
		}
	}
	return false
}

// ReadAll returns all the items found in src,
// skipping blanks, and comment lines prefixed with
// "//" and "#".
func ReadAll[T any](src io.Reader) ([]T, error) {
	s := NewScanner(src)
	s.SkipBlank = true
	s.SkipComments = []string{"//", "#"}

	result := make([]T, 0)
	for s.Next() {
		line, err := s.Line()
		if err != nil {
			return nil, fmt.Errorf("jsonl: reading error (line: %d): %v", s.loc, err)
		}

		v := new(T)
		if err := line.Scan(v); err != nil {
			return nil, fmt.Errorf("jsonl: scanning error (line: %d): %v", s.loc, err)
		}
		result = append(result, *v)
	}

	return result, nil
}
