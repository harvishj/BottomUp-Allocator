package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	input         []Instruction
	virtualRegs   = make(map[string]*Register)
	phyVirMap     = make(map[string]string)
	BASE_ADDRESS  = 4096
	CURR_OFFSET   = 0
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: <registers> <inputFile> <outputFile>")
		return
	}

	registers, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("First argument passed is not an integer.")
		return
	}

	inputFile := os.Args[2]
	outputFile := os.Args[3]

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Println("Input file does not exist.")
		return
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Error creating output file.")
			return
		}
		file.Close()
	}

	for i := 0; i < registers; i++ {
		phyVirMap["r"+string(rune(i+97))] = "NULL"
	}

	getInputInstructions(inputFile)
	getRegistersAndNextUse()
	getLiveRegs()

	err = allocateRegs(outputFile)
	if err != nil {
		fmt.Println("Error during register allocation:", err)
	}
}

func getInputInstructions(inputFile string) {
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
			inst := ParseInstruction(line)
			inst.SetLineNumber(lineNumber)
			input = append(input, *inst)
			lineNumber++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input file:", err)
	}
}

func getRegistersAndNextUse() {
	for _, inst := range input {
		// Check sources
		for _, src := range inst.Sources {
			if strings.HasPrefix(src, "r") && virtualRegs[src] == nil {
				virtualRegs[src] = NewRegisterWithUsage(src, inst.LineNumber, inst.LineNumber)
			} else if virtualRegs[src] != nil {
				reg := virtualRegs[src]
				if reg.NextUse == 0 {
					reg.SetNextUse(inst.LineNumber)
				}
				reg.SetLastUse(inst.LineNumber)
			}
		}

		// Check targets
		for _, targ := range inst.Targets {
			if strings.HasPrefix(targ, "r") && virtualRegs[targ] == nil {
				virtualRegs[targ] = NewRegisterWithUsage(targ, inst.LineNumber, inst.LineNumber)
			} else if virtualRegs[targ] != nil {
				reg := virtualRegs[targ]
				if reg.NextUse == 0 {
					reg.SetNextUse(inst.LineNumber)
				}
				reg.SetLastUse(inst.LineNumber)
			}
		}
	}

	for _, reg := range virtualRegs {
		fmt.Println(reg)
	}
}

func getLiveRegs() {
	for i, inst := range input {
		if i != 0 {
			inst.LiveRegs = input[i-1].LiveRegs
		}

		for _, src := range inst.Sources {
			if strings.HasPrefix(src, "r") {
				inst.LiveRegs[src] = struct{}{}
			}
		}

		for _, targ := range inst.Targets {
			if strings.HasPrefix(targ, "r") {
				inst.LiveRegs[targ] = struct{}{}
			}
		}

		for reg := range inst.LiveRegs {
			if virtualRegs[reg] != nil && virtualRegs[reg].LastUse == inst.LineNumber {
				delete(inst.LiveRegs, reg)
			}
		}
		input[i] = inst
	}

	for _, inst := range input {
		fmt.Println(inst.String())
	}
}

func allocateRegs(outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("could not create output file: %v", err)
	}
	defer file.Close()

	bw := bufio.NewWriter(file)
	defer bw.Flush()

	for _, inst := range input {
		updateRegisterNextUse(inst.LineNumber - 1)
		updateRegMap(inst.LiveRegs)

		regs := make(map[string]struct{})
		for _, src := range inst.Sources {
			if strings.HasPrefix(src, "r") {
				regs[src] = struct{}{}
			}
		}

		for _, targ := range inst.Targets {
			if strings.HasPrefix(targ, "r") {
				regs[targ] = struct{}{}
			}
		}

		for register := range regs {
			curr := virtualRegs[register]
			if !containsValue(phyVirMap, register) && curr.LastUse >= inst.LineNumber {
				virtualReg := virtualRegs[register]
				if virtualReg.Offset != math.MaxInt32 {
					address := BASE_ADDRESS + virtualReg.Offset
					emptyPhyReg := findEmptyPhyReg()
					if emptyPhyReg != "NULL" {
						virtualReg.Allocated = 1
						phyVirMap[emptyPhyReg] = virtualReg.Name
					} else {
						emptyPhyReg = getMaxNextUseAndSpill(bw)
					}
					loadSpill := fmt.Sprintf("load\t%d => %s", address, emptyPhyReg)
					fmt.Println("write 1:", loadSpill)
					bw.WriteString(loadSpill + "\n")
				} else {
					flag := 0
					emptyPhyReg := findEmptyPhyReg()
					if emptyPhyReg != "NULL" {
						virReg := virtualRegs[register]
						virReg.Allocated = 1
						phyVirMap[emptyPhyReg] = virReg.Name
						fmt.Println()
						flag = 1
					}

					if flag == 0 {
						emptiedReg := getMaxNextUseAndSpill(bw)
						virReg := virtualRegs[register]
						virReg.Allocated = 1
						phyVirMap[emptiedReg] = register
						virtualRegs[register] = virReg
					}
				}
			}
		}

		// Safely construct the instruction string
		instruction := inst.Opcode + "\t"
		if len(inst.Sources) > 0 {
			if len(inst.Sources) == 2 {
				instruction += getMappedPhyReg(inst.Sources[0]) + ", " + getMappedPhyReg(inst.Sources[1])
			} else {
				instruction += getMappedPhyReg(inst.Sources[0])
			}
		}

		if len(inst.Targets) > 0 {
			instruction += " => "
			if len(inst.Targets) == 2 {
				instruction += getMappedPhyReg(inst.Targets[0]) + ", " + getMappedPhyReg(inst.Targets[1])
			} else {
				instruction += getMappedPhyReg(inst.Targets[0])
			}
		}

		fmt.Println("write 2:", instruction)
		bw.WriteString(instruction + "\n")
	}

	return nil
}


func updateRegMap(liveRegs map[string]struct{}) {
	for reg := range phyVirMap {
		if _, ok := liveRegs[phyVirMap[reg]]; !ok {
			phyVirMap[reg] = "NULL"
		}
	}
}

func updateRegisterNextUse(current int) {
	for i := current; i < len(input); i++ {
		inst := input[i]
		for _, src := range inst.Sources {
			if strings.HasPrefix(src, "r") {
				temp := virtualRegs[src]
				if temp.NextUse < current+1 {
					temp.NextUse = i + 1
					virtualRegs[temp.Name] = temp
					fmt.Println("inst:", current+1, "updated next use:", temp)
				}
			}
		}

		for _, targ := range inst.Targets {
			if strings.HasPrefix(targ, "r") {
				temp := virtualRegs[targ]
				if temp.NextUse < current+1 {
					temp.NextUse = i + 1
					virtualRegs[temp.Name] = temp
					fmt.Println("inst:", current+1, "updated next use:", temp)
				}
			}
		}
	}
}

func getMappedPhyReg(virtualReg string) string {
	if strings.HasPrefix(virtualReg, "r") {
		for phy, vir := range phyVirMap {
			if vir == virtualReg {
				return phy
			}
		}
	} else {
		return virtualReg
	}
	return ""
}

func getMaxNextUseAndSpill(bw *bufio.Writer) string {
	maxNextUse := 0
	var regToSpill *Register
	regName := ""
	for phy, vir := range phyVirMap {
		currReg := virtualRegs[vir]
		if currReg != nil {
			if currReg.NextUse > maxNextUse {
				maxNextUse = currReg.NextUse
				regToSpill = currReg
				regName = phy
			}
		}
	}

	if regToSpill != nil {
		regToSpill.Allocated = 0
		regToSpill.Offset = CURR_OFFSET
		address := BASE_ADDRESS + CURR_OFFSET
		CURR_OFFSET += 4
		spillCode := fmt.Sprintf("store\t%s => %d", regName, address)
		fmt.Println("write 3:", spillCode)
		bw.WriteString(spillCode + "\n")
		phyVirMap[regName] = "NULL"
		virtualRegs[regToSpill.Name] = regToSpill
	}

	return regName
}

func findEmptyPhyReg() string {
	for phy, vir := range phyVirMap {
		if vir == "NULL" {
			return phy
		}
	}
	return "NULL"
}

func containsValue(m map[string]string, value string) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}
