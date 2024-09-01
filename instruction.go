package main

import (
	"strconv"
	"strings"
)

// Instruction struct to represent an instruction
type Instruction struct {
	Opcode    string
	Sources   []string
	Targets   []string
	NumRegs   int
	LineNumber int
	LiveRegs  map[string]struct{}
}

// NewInstruction initializes a new Instruction
func NewInstruction(opcode string, sources []string, targets []string, numRegs int) *Instruction {
	return &Instruction{
		Opcode:    opcode,
		Sources:   sources,
		Targets:   targets,
		NumRegs:   numRegs,
		LiveRegs:  make(map[string]struct{}),
	}
}

// String returns a string representation of the Instruction
func (instr *Instruction) String() string {
	var str strings.Builder
	str.WriteString("line number: ")
	str.WriteString(strconv.Itoa(instr.LineNumber))
	str.WriteString("\topcode: ")
	str.WriteString(instr.Opcode)
	str.WriteString("\tnum regs: ")
	str.WriteString(strconv.Itoa(instr.NumRegs))
	str.WriteString("\tsources: ")

	if len(instr.Sources) == 2 {
		str.WriteString(instr.Sources[0] + " " + instr.Sources[1])
	} else if len(instr.Sources) == 1 {
		str.WriteString(instr.Sources[0])
	}

	if len(instr.Targets) > 0 {
		str.WriteString("\ttargets: ")
		if len(instr.Targets) == 2 {
			str.WriteString(instr.Targets[0] + " " + instr.Targets[1])
		} else {
			str.WriteString(instr.Targets[0])
		}
	}

	str.WriteString("\t live regs: ")
	for s := range instr.LiveRegs {
		str.WriteString(s + "\t")
	}
	return str.String()
}

// ParseInstruction parses a string into an Instruction
func ParseInstruction(instruction string) *Instruction {
	instruction = strings.ReplaceAll(instruction, ",", "")
	split := strings.Fields(instruction)

	op := split[1]
	var src, tar []string

	numRegs := 0
	flag := 0

	for i := 2; i < len(split); i++ {
		if split[i] == "=>" {
			flag = 1
			continue
		} else if flag == 0 {
			src = append(src, split[i])
		} else {
			tar = append(tar, split[i])
		}
	}

	sources := src
	targets := tar

	for i := 0; i < len(sources); i++ {
		if strings.HasPrefix(sources[i], "r") {
			numRegs++
		}
	}
	for i := 0; i < len(targets); i++ {
		if strings.HasPrefix(targets[i], "r") {
			numRegs++
		}
	}

	return NewInstruction(op, sources, targets, numRegs)
}

// SetLineNumber sets the line number for the instruction
func (instr *Instruction) SetLineNumber(i int) {
	instr.LineNumber = i
}
