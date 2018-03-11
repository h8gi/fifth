package vm

import "testing"

func TestOps(t *testing.T) {
	vm := NewVM()
	vm.memory[0] = OpLit
	vm.memory[1] = Cell(200)
	vm.memory[2] = OpLit
	vm.memory[3] = Cell(100)
	vm.memory[4] = OpAdd
	vm.memory[5] = OpEnd
	vm.Run()
	if vm.ds.tos() != Cell(300) {
		t.Error("add")
	}
}

func TestIf(t *testing.T) {
	vm := NewVM()
	vm.memory[0] = OpLit
	vm.memory[1] = Cell(0) // push false
	vm.memory[2] = OpIf
	vm.memory[3] = Cell(5) // jump address (5 mean 6)
	vm.memory[4] = OpLit
	vm.memory[5] = Cell(42) // true value
	vm.memory[6] = OpLit
	vm.memory[7] = Cell(24) // false value
	vm.memory[8] = OpEnd
	vm.Verbose = true
	if err := vm.Run(); err != nil {
		t.Error(err)
	}
	actual, err := vm.ds.pop()
	if err != nil {
		t.Error(err)
	}
	if actual != Cell(24) {
		t.Error("aaaa")
	}
}

func TestStore(t *testing.T) {
	addr := Cell(0)
	value := Cell(42)

	vm := NewVM()
	if err := vm.ds.push(value); err != nil {
		t.Error(err)
	}
	if err := vm.ds.push(addr); err != nil {
		t.Error(err)
	}
	if err := vm.InstStore(); err != nil {
		t.Error(err)
	}

	if vm.memory[addr] != value {
		t.Errorf("expected %v, actual %v", value, vm.memory[addr])
	}
}

func TestFetch(t *testing.T) {
	addr := Cell(0)
	value := Cell(42)

	vm := NewVM()
	// 42
	if err := vm.ds.push(value); err != nil {
		t.Error(err)
	}
	// 42 0
	if err := vm.ds.push(addr); err != nil {
		t.Error(err)
	}
	//
	if err := vm.InstStore(); err != nil {
		t.Error(err)
	}
	// 0
	if err := vm.ds.push(addr); err != nil {
		t.Error(err)
	}
	// 42

	if err := vm.InstFetch(); err != nil {
		t.Error(err)
	}

	if vm.ds.tos() != value {

		t.Errorf("expect %v, actual %v", value, vm.ds.tos())
	}
}
