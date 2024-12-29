package interpreter

import (
	"sync"
)

// execute runs the passed program, using the jump table and input/output channels
func (i *Interpreter) execute(program []byte, jumps map[int]int, input <-chan byte, output chan<- byte) error {
	var wg sync.WaitGroup
	// Execute
	for i.inPtr < len(program)-1 {
		b := program[i.inPtr]
		switch b {
		case '>':
			i.tape.Right()
		case '<':
			i.tape.Left()
		case '+':
			i.tape.Add()
		case '-':
			i.tape.Sub()
		case '.':
			output <- i.tape.Read()
		case ',':
			i.tape.Write(<-input)
		case '[':
			if i.tape.Zero() {
				i.inPtr = jumps[i.inPtr]
			}
		case ']':
			if !i.tape.Zero() {
				i.inPtr = jumps[i.inPtr]
			}
		case 'Y':
			wg.Add(1)
			child := i.Fork()
			go func() {
				defer wg.Done()
				child.execute(program, jumps, input, output)
			}()

			i.tape.Write(0)
		}

		i.inPtr++
	}

	// Wait for all forked children to finish
	wg.Wait()
	return nil
}
