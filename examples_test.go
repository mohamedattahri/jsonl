package jsonl_test

import (
	"io"
	"log"

	"mohamed.attahri.com/jsonl"
)

type Item = string

var (
	src                 io.Reader
	dest                io.Writer
	item1, item2, item3 *Item
)

func ExampleScanner() {
	s := jsonl.NewScanner(src)
	for s.Next() {
		line, err := s.Line()
		if err != nil {
			log.Fatal(err)
		}

		var item Item
		if err := line.Scan(&item); err != nil {
			log.Fatal(err)
		}

		log.Println(item)
	}
}

func ExampleWriter() {
	w := jsonl.NewWriter[*Item](dest)
	if _, err := w.Write(item1, item2, item3); err != nil {
		log.Fatal(err)
	}
}

func ExampleReadAll() {
	items, err := jsonl.ReadAll[*Item](src)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(items))
}
