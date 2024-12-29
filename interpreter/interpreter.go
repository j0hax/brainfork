package interpreter

import (
	"bufio"
	"io"
	"sync"
)

type Interpreter struct {
	input  *bufio.Reader
	output *bufio.Writer
	tape   *Tape
	inPtr  int
}

// NewInterpreter creates a new Brainfork interpreter, which reads input from r
// and writes output to w. Its memory is set to the original 30 000 bytes.
func NewInterpreter(r io.Reader, w io.Writer) *Interpreter {
	return &Interpreter{
		input:  bufio.NewReader(r),
		output: bufio.NewWriter(w),
		tape:   NewTape(),
		inPtr:  0,
	}
}

// compileJumpTable parses a program, returning a map of jump indices for the
// '[' and ']' operators.
func compileJumpTable(code []byte) map[int]int {
	var stack []int
	jumps := make(map[int]int)
	for i, c := range code {
		if c == '[' {
			stack = append(stack, i)
		} else if c == ']' {
			jumps[i], stack = stack[len(stack)-1], stack[:len(stack)-1]
			jumps[jumps[i]] = i
		}
	}

	return jumps
}

// Fork copies the interpreter's memory, increments the data and instruction
// pointers, and returns a new Interpreter, which can then be executed
// concurrently.
func (i *Interpreter) Fork() *Interpreter {
	newMem := i.tape.Clone()
	newMem.Right()
	newMem.Write(1)

	return &Interpreter{
		input:  i.input,
		output: i.output,
		tape:   newMem,
		inPtr:  i.inPtr + 1,
	}
}

// Run Executes a program in the interpreter,
func (i *Interpreter) Run(program []byte) error {
	var wg sync.WaitGroup

	jumps := compileJumpTable(program)

	in := make(chan byte)
	out := make(chan byte)

	// Reads and Writes are not thread safe;
	// as a result, we start concurrent closures which read and write
	// out-of-order data
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer i.output.Flush()
		for o := range out {
			i.output.WriteByte(o)
		}
	}()

	// Run goroutine to direct input to channels.
	// As a rule of thumb, writers close channels, hence the defer
	go func() {
		defer close(in)
		for {
			b, err := i.input.ReadByte()
			if err != nil {
				break
			}
			in <- b
		}
	}()

	err := i.execute(program, jumps, in, out)
	close(out)
	wg.Wait()

	if err != nil {
		return err
	}

	return nil
}
