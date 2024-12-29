package interpreter_test

import (
	"os"
	"strings"
	"testing"

	"github.com/j0hax/brainfork/interpreter"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestFork(t *testing.T) {
	p, err := os.ReadFile("./hacker.bf")
	if err != nil {
		t.Error(err)
	}

	var out strings.Builder
	var in strings.Reader

	bf := interpreter.NewInterpreter(&in, &out)
	err = bf.Run(p)

	// Because this program is multithreaded, outputs depend on scheduling and
	// are non-deterministic. As a result, we check to make sure that there are
	// twice as of each character than expected, regardless of order

	expected := "Just another brainfuck hacker."
	want := len(expected) * 2
	have := len(out.String())
	if have != want {
		t.Fatalf(`Have output of length %d, want output with %d`, have, want)
	}

	actual := out.String()
	chars := make(map[rune]int)
	for _, c := range expected {
		chars[c]++
	}

	for k, v := range chars {
		// Character in output that is not expected
		if !strings.ContainsRune(actual, k) {
			t.Fatalf("Unexpected rune %c", k)
		}

		// A character in the output is not twice as many times as the input
		charCount := strings.Count(actual, string(k))
		if charCount != v*2 {
			t.Fatalf("Rune %c appears %d times instead of %d", k, charCount, chars[k])
		}
	}
}
