package jsonl_test

import (
	"bytes"
	"io"
	"testing"

	"mohamed.attahri.com/jsonl"
)

var FixtureString = `{"name":"Gilbert","wins":[["straight","7♣"],["one pair","10♥"]]}
{"name":"Alexa","wins":[["two pair","4♠"],["two pair","9♠"]]}
{"name":"May","wins":[]}
{"name":"Deloise","wins":[["three of a kind","5♣"]]}`

var FixturesStringComments = `{"name":"Gilbert","wins":[["straight","7♣"],["one pair","10♥"]]}
// first comment
# second comment

{"name":"Alexa","wins":[["two pair","4♠"],["two pair","9♠"]]}
// third comment
{"name":"May","wins":[]}
# fourth comment

{"name":"Deloise","wins":[["three of a kind","5♣"]]}`

type Player struct {
	Name string      `json:"name"`
	Wins [][2]string `json:"wins"`
}

func (p *Player) Equal(other *Player) bool {
	if p.Name != other.Name {
		return false
	}
	if len(p.Wins) != len(other.Wins) {
		return false
	}
	for i, w := range p.Wins {
		if w != other.Wins[i] {
			return false
		}
	}
	return true
}

var Fixtures = []*Player{
	{
		Name: "Gilbert",
		Wins: [][2]string{
			{"straight", "7♣"},
			{"one pair", "10♥"}},
	},
	{
		Name: "Alexa",
		Wins: [][2]string{
			{"two pair", "4♠"},
			{"two pair", "9♠"},
		},
	},
	{
		Name: "May",
		Wins: [][2]string{},
	},
	{
		Name: "Deloise",
		Wins: [][2]string{
			{"three of a kind", "5♣"},
		},
	},
}

func TestWriter(t *testing.T) {
	out := &bytes.Buffer{}

	w := jsonl.NewWriter[*Player](out)
	n, err := w.Write(Fixtures...)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("unexpected n value")
	}
	if n != w.Written() {
		t.Fatalf("bytes written mismatch: %d vs %d", n, w.Written())
	}

	if out.String() != FixtureString {
		t.Fatal("not matching")
	}
}

func TestScanner(t *testing.T) {
	rd := bytes.NewBufferString(FixtureString)
	scan := jsonl.NewScanner(rd)
	for scan.Next() {
		line, err := scan.Line()
		if err != nil {
			t.Fatal(err)
		}
		p := new(Player)
		if err := line.Scan(p); err != nil {
			t.Fatal(err)
		}
	}
}

func testReadAll(t *testing.T, src io.Reader) {
	players, err := jsonl.ReadAll[*Player](src)
	if err != nil {
		t.Fatal(err)
	}
	if len(players) != len(Fixtures) {
		t.Fatal("count does not match")
	}

	for i, player := range players {
		if !player.Equal(Fixtures[i]) {
			t.Fatalf("Player %d does not match", i)
		}
	}
}

func TestReadAll(t *testing.T) {
	t.Run("no comments", func(t *testing.T) {
		testReadAll(t, bytes.NewBufferString(FixtureString))
	})
	t.Run("with comments", func(t *testing.T) {
		testReadAll(t, bytes.NewBufferString(FixturesStringComments))
	})
}

func TestScannerError(t *testing.T) {
	rd := bytes.NewBufferString(FixturesStringComments)
	scan := jsonl.NewScanner(rd)
	for scan.Next() {
		continue
	}
	// blank lines should cause an error
	if err := scan.Err(); err == nil {
		t.Fatal("expected an error")
	}
}
