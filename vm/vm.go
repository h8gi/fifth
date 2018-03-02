package vm

type Cell int

const (
	OpLit Cell = iota
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

	dsSize     = 128
	rsSize     = 2048
	memorySize = 524288 * 4
)

type VM struct {
	ds     *stack
	rs     *stack
	pc     Cell
	memory []Cell
}

func NewVM() *VM {
	return &VM{
		ds:     newStack(dsSize),
		rs:     newStack(rsSize),
		pc:     Cell(0),
		memory: make([]Cell, memorySize),
	}
}

func (vm *VM) InstLit() error {
	vm.pc++
	return vm.ds.push(vm.memory[vm.pc])
}

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
