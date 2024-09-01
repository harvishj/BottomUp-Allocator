package main

import (
	"math"
	"strconv"
)

// Register struct to represent a register
type Register struct {
	Name     string
	Allocated int
	FirstUse  int
	NextUse   int
	LastUse   int
	Offset    int
}

// NewRegister initializes a new Register with default values
func NewRegister(name string) *Register {
	return &Register{
		Name:      name,
		Allocated: 0,
		FirstUse:  0,
		NextUse:   0,
		LastUse:   math.MaxInt32,
		Offset:    math.MaxInt32,
	}
}

// NewRegisterWithUsage initializes a new Register with usage information
func NewRegisterWithUsage(name string, firstUse int, lastUse int) *Register {
	return &Register{
		Name:      name,
		FirstUse:  firstUse,
		NextUse:   0,
		Allocated: 0,
		LastUse:   lastUse,
		Offset:    math.MaxInt32,
	}
}

// String returns a string representation of the Register
func (reg *Register) String() string {
	return "name: " + reg.Name + "\t allocated: " + strconv.Itoa(reg.Allocated) + "\t first use: " + strconv.Itoa(reg.FirstUse) + "\t next use: " + strconv.Itoa(reg.NextUse) + "\t last use: " + strconv.Itoa(reg.LastUse) + "\t offset: " + strconv.Itoa(reg.Offset)
}

// SetNextUse sets the next use line number for the register
func (reg *Register) SetNextUse(lineNumber int) {
	reg.NextUse = lineNumber
}

// SetLastUse sets the last use line number for the register
func (reg *Register) SetLastUse(lineNumber int) {
	reg.LastUse = lineNumber
}
