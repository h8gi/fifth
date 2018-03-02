package vm

import "testing"

func TestStore(t *testing.T) {
	addr := Cell(0)
	value := Cell(10)

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
