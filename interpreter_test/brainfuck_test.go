package interpreter_test

import (
	"os"
	"strings"
	"testing"

	"github.com/j0hax/brainfork/interpreter"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloWorld(t *testing.T) {
	p, err := os.ReadFile("./helloworld.bf")
	if err != nil {
		t.Error(err)
	}

	var out strings.Builder
	var in strings.Reader

	bf := interpreter.NewInterpreter(&in, &out)
	err = bf.Run(p)

	want := "Hello World!\n"
	if out.String() != want {
		t.Fatalf(`%q, %v, want match for %#q, nil`, out.String(), err, want)
	}
}
