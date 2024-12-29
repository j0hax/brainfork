package interpreter

import "maps"

type Tape struct {
	mem     map[int]byte
	pointer int
}

// Left decreases the tape position by one
func (t *Tape) Left() {
	t.pointer--
}

// Right increases the tape position by one
func (t *Tape) Right() {
	t.pointer++
}

// Add increases the value of the current cell by one
func (t *Tape) Add() {
	t.mem[t.pointer]++
}

// Sub decreases the value of the current cell by one
func (t *Tape) Sub() {
	t.mem[t.pointer]--
}

// Read reads one byte from the current tape position
func (t *Tape) Read() byte {
	return t.mem[t.pointer]
}

func (t *Tape) Zero() bool {
	return t.mem[t.pointer] == 0
}

// Write writes one byte to the current tape position
func (t *Tape) Write(b byte) {
	t.mem[t.pointer] = b
}

func NewTape() *Tape {
	return &Tape{
		mem:     make(map[int]byte),
		pointer: 0,
	}
}

func (t *Tape) Clone() *Tape {
	return &Tape{
		mem:     maps.Clone(t.mem),
		pointer: t.pointer,
	}
}
