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
			i.dtPtr++
		case '<':
			i.dtPtr--
		case '+':
			i.mem[i.dtPtr]++
		case '-':
			i.mem[i.dtPtr]--
		case '.':
			output <- i.mem[i.dtPtr]
		case ',':
			i.mem[i.dtPtr] = <-input
		case '[':
			if i.mem[i.dtPtr] == 0 {
				i.inPtr = jumps[i.inPtr]
			}
		case ']':
			if i.mem[i.dtPtr] != 0 {
				i.inPtr = jumps[i.inPtr]
			}
		case 'Y':
			wg.Add(1)
			child := i.Fork()
			go func() {
				defer wg.Done()
				child.execute(program, jumps, input, output)
			}()

			i.mem[i.dtPtr] = 0
		}

		i.inPtr++
	}

	// Wait for all forked children to finish
	wg.Wait()
	return nil
}
