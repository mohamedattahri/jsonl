[![GoDoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/mohamed.attahri.com/jsonl)

`jsonl` is a Go library to encode and decode [JSON Lines](https://jsonlines.org). It uses the standard library's JSON package.

Go 1.20+ is required, because it relies on generics for type-safety.

See [documentation](https://pkg.go.dev/mohamed.attahri.com/jsonl) for more details.

## Encode

Use `Writer` to encode items into JSON
before writing them to an underlying `io.Writer` instance.

```golang
w := jsonl.NewWriter[*Item](dest)
if _, err := w.Write(item1, item2); err != nil {
  log.Fatal(err)
}
```

## Decode

Create a scanner with `NewScanner` to iterate over a source one line at a time, or call the `ReadAll` function to read all the lines at once.

Here's a example that uses `Scanner`:

```golang
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

// Check if an error occurred while scanning.
if err := s.Err(); err != nil {
  log.Fatalf("an error happened: %v", err)
}
```
