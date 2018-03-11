package vm

import "fmt"

type Cell int

func (c Cell) Bool() bool {
	return c != 0
}

const (
	OpNop Cell = iota
	OpLit
	OpStore
	OpFetch

	OpAdd
	OpSub

	OpAnd
	OpOr
	OpXor

	OpDup
	OpDrop
	OpOver
	OpSwap

	OpDsToRs // >R
	OpRsToDs // R>

	OpIf
	OpCall
	OpExit

	OpEnd

	dsSize     = 128
	rsSize     = 2048
	memorySize = 524288 * 4
)

type VM struct {
	ds      *stack
	rs      *stack
	pc      Cell
	memory  []Cell
	opMap   map[Cell]func() error
	Verbose bool
}

func NewVM() *VM {
	vm := &VM{
		ds:     newStack(dsSize),
		rs:     newStack(rsSize),
		pc:     Cell(0),
		memory: make([]Cell, memorySize),
		opMap:  make(map[Cell]func() error),
	}
	vm.InstallOps()
	return vm
}

func (vm *VM) Run() error {
	for vm.pc < memorySize {
		if vm.Verbose {
			fmt.Println(vm.ds)
		}
		opcode := vm.memory[vm.pc]
		op, ok := vm.opMap[opcode]
		if !ok {
			return fmt.Errorf("illegal opcode: %v", opcode)
		}
		if err := op(); err != nil {
			return err
		}
		vm.pc++
	}
	return nil
}

func (vm *VM) InstNop() error {
	return nil
}

// InstLit push the next memory content to data stack.
// -- N1
func (vm *VM) InstLit() error {
	vm.pc++
	return vm.ds.push(vm.memory[vm.pc])
}

// N1 ADDR --
func (vm *VM) InstStore() error {
	addr, err := vm.ds.pop()
	if err != nil {
		return err
	}
	n1, err := vm.ds.pop()
	if err != nil {
		return err
	}
	vm.memory[addr] = n1
	return nil
}

// ADDR -- N1
func (vm *VM) InstFetch() error {
	addr, err := vm.ds.pop()
	if err != nil {
		return err
	}
	n1 := vm.memory[addr]
	return vm.ds.push(n1)
}

// N1 N2 -- N
func (vm *VM) InstAdd() error {
	n2, err := vm.ds.pop()
	if err != nil {
		return err
	}
	n1, err := vm.ds.pop()
	if err != nil {
		return err
	}
	return vm.ds.push(n1 + n2)
}

// N1 N2 -- N
func (vm *VM) InstSub() error {
	n2, err := vm.ds.pop()
	if err != nil {
		return err
	}
	n1, err := vm.ds.pop()
	if err != nil {
		return err
	}
	return vm.ds.push(n2 - n1)
}

// N1 N2 -- N
func (vm *VM) InstAnd() error {
	n2, err := vm.ds.pop()
	if err != nil {
		return err
	}
	n1, err := vm.ds.pop()
	if err != nil {
		return err
	}
	return vm.ds.push(n1 & n2)
}

// N1 N2 -- N
func (vm *VM) InstOr() error {
	n2, err := vm.ds.pop()
	if err != nil {
		return err
	}
	n1, err := vm.ds.pop()
	if err != nil {
		return err
	}
	return vm.ds.push(n1 | n2)
}

func (vm *VM) InstXor() error {
	n2, err := vm.ds.pop()
	if err != nil {
		return err
	}
	n1, err := vm.ds.pop()
	if err != nil {
		return err
	}
	return vm.ds.push(n1 ^ n2)
}

func (vm *VM) InstEnd() error {
	vm.pc = memorySize
	return nil
}

// n -- n n
func (vm *VM) InstDup() error {
	n := vm.ds.tos()
	return vm.ds.push(n)
}

func (vm *VM) InstDrop() error {
	_, err := vm.ds.pop()
	return err
}

// n1 n2 -- n1 n2 n1
func (vm *VM) InstOver() error {
	n2, err := vm.ds.pop()
	if err != nil {
		return err
	}
	n1 := vm.ds.tos()
	err = vm.ds.push(n2)
	if err != nil {
		return err
	}
	return vm.ds.push(n1)
}

// n1 n2 -- n2 n1
func (vm *VM) InstSwap() error {
	n2, err := vm.ds.pop()
	if err != nil {
		return err
	}
	n1, err := vm.ds.pop()
	if err != nil {
		return err
	}
	err = vm.ds.push(n2)
	if err != nil {
		return err
	}
	return vm.ds.push(n1)
}

func (vm *VM) InstDsToRs() error {
	n, err := vm.ds.pop()
	if err != nil {
		return err
	}
	return vm.rs.push(n)
}

func (vm *VM) InstRsToDs() error {
	n, err := vm.rs.pop()
	if err != nil {
		return err
	}
	return vm.ds.push(n)
}

// n --
// Next memory cell is jump address.
// Jump address should be actual address - 1
func (vm *VM) InstIf() error {
	n, err := vm.ds.pop()
	if err != nil {
		return err
	}

	if n.Bool() {
		// do nothing
		return nil
	}
	vm.pc++
	vm.pc = vm.memory[vm.pc]
	return nil
}

func (vm *VM) InstCall() error {
	// save return address
	if err := vm.rs.push(vm.pc); err != nil {
		return err
	}
	vm.pc++
	// -1 address
	vm.pc = vm.memory[vm.pc]
	return nil
}

func (vm *VM) InstExit() error {
	addr, err := vm.rs.pop()
	if err != nil {
		return err
	}
	vm.pc = addr
	return nil
}

func (vm *VM) InstallOps() {
	vm.opMap[OpNop] = vm.InstNop
	vm.opMap[OpLit] = vm.InstLit
	vm.opMap[OpStore] = vm.InstStore
	vm.opMap[OpFetch] = vm.InstFetch
	vm.opMap[OpAdd] = vm.InstAdd
	vm.opMap[OpSub] = vm.InstSub
	vm.opMap[OpAnd] = vm.InstAnd
	vm.opMap[OpOr] = vm.InstOr
	vm.opMap[OpXor] = vm.InstXor
	vm.opMap[OpDup] = vm.InstDup
	vm.opMap[OpOver] = vm.InstOver
	vm.opMap[OpSwap] = vm.InstSwap
	vm.opMap[OpDsToRs] = vm.InstDsToRs
	vm.opMap[OpRsToDs] = vm.InstRsToDs
	vm.opMap[OpIf] = vm.InstIf
	vm.opMap[OpCall] = vm.InstCall
	vm.opMap[OpExit] = vm.InstExit
	vm.opMap[OpEnd] = vm.InstEnd
}
