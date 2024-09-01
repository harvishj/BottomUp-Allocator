# Register Allocation Program in Go

This program implements a register allocation algorithm in Go. The program reads instructions from an input file, allocates registers, handles spills, and writes the resulting instructions to an output file.

## Table of Contents

- [Register Allocation Program in Go](#register-allocation-program-in-go)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Requirements](#requirements)
  - [Installation](#installation)
    - [Arguments](#arguments)
  - [Example](#example)
  - [Files](#files)

## Introduction

The purpose of this program is to demonstrate a bottom-up register allocation strategy. The program reads a list of instructions from an input file, allocates registers, handles spills when there are not enough registers available, and writes the modified instructions to an output file.

## Requirements

- Go 1.16 or higher
- A terminal or command prompt to run the program
- An input file (`input.i`) containing assembly-like instructions to be processed

## Installation

1. **Clone the repository or download the source files**:

   ```bash
   git clone https://github.com/yourusername/register-allocator-go.git
   cd register-allocator-go

2. Ensure all source files are in the same directory:

    alloc.go
    instruction.go
    register.go

## Usage
To run the program, use the go run command with the necessary arguments:

```bash
go run *.go <number_of_physical_registers> <input_file> <output_file>
```

### Arguments
<number_of_physical_registers>: The number of physical registers available for allocation.

<input_file>: The path to the input file containing instructions.

<output_file>: The path to the output file where the modified instructions will be written.

## Example

Assuming you have input.i as your input file and want to use output.i as your output file with 3 physical registers:

```bash
go run *.go 3 input.i output.i
```


## Files
* alloc.go: The main file that contains the entry point and the main register allocation logic.
* instruction.go: Defines the Instruction struct and methods for parsing and handling instructions.
* register.go: Defines the Register struct and methods for handling register usage, allocation, and spills.

