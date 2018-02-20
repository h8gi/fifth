package fifth

import "fmt"

func (i *Interpreter) LoadPrimitives() {
	i.SetPrimitive("+", i.add())
	i.SetPrimitive("-", i.sub())
	i.SetPrimitive("*", i.mult())
	i.SetPrimitive("/", i.div())
	i.SetPrimitive(".", i.dot)
	i.SetPrimitive("dup", i.dup)
	i.SetPrimitive(".s", i.show)
	i.SetPrimitive(".r", i.rShow)

	i.SetPrimitive(":", i.colonDefine)

	i.SetPrimitive(";", i.semicolon)
	i.Dictionary[";"].IsImmediate = true
	i.Dictionary[";"].IsCompileOnly = true

	i.SetPrimitive("see", i.see)
	// compilation behavior of if
	i.SetPrimitive("if", i.innerIf)
	i.Dictionary["if"].IsImmediate = true
	i.Dictionary["if"].IsCompileOnly = true

	i.SetPrimitive("else", i.innerElse)
	i.Dictionary["else"].IsImmediate = true
	i.Dictionary["else"].IsCompileOnly = true

	i.SetPrimitive("then", i.innerThen)
	i.Dictionary["then"].IsImmediate = true
	i.Dictionary["then"].IsCompileOnly = true

	i.SetPrimitive("immediate", i.immediate)

	i.SetPrimitive("literal", i.literal)
	i.Dictionary["literal"].IsImmediate = true
	i.Dictionary["literal"].IsCompileOnly = true

	i.SetPrimitive("'", i.tick)

	i.SetPrimitive("execute", i.exec)

	i.SetPrimitive("bye", i.bye)

}

func (i *Interpreter) add() PrimBody {
	return i.DS.MakeBinFunc(func(n, m int) int { return n + m })
}

func (i *Interpreter) sub() PrimBody {
	return i.DS.MakeBinFunc(func(n, m int) int { return m - n })
}

func (i *Interpreter) mult() PrimBody {
	return i.DS.MakeBinFunc(func(n, m int) int { return n * m })
}

func (i *Interpreter) div() PrimBody {
	return i.DS.MakeBinFunc(func(n, m int) int { return m / n })
}

func (i *Interpreter) dot() error {
	n, err := i.DS.Pop()
	if err != nil {
		return err
	}
	fmt.Println(n)
	return nil
}

func (i *Interpreter) dup() error {
	n, err := i.DS.Pop()
	if err != nil {
		return err
	}
	i.DS.Push(n)
	i.DS.Push(n)
	return nil

}

func (i *Interpreter) show() error {
	fmt.Println(i.DS.data)
	return nil

}

func (i *Interpreter) rShow() error {
	fmt.Println(i.RS.data)
	return nil

}

// : <name> word ... ;
func (i *Interpreter) colonDefine() error {
	i.IsCompile = true
	if !i.Scanner.Scan() {
		return EOFError
	}
	i.CWord = &Word{
		Name: i.Scanner.Text(),
	}
	return nil
}

func (i *Interpreter) semicolon() error {
	i.IsCompile = false
	i.Dictionary[i.CWord.Name] = i.CWord
	return nil
}

func (i *Interpreter) see() error {
	if !i.Scanner.Scan() {
		return EOFError
	}
	name := i.Scanner.Text()
	word, ok := i.Dictionary[name]
	if !ok {
		return UndefinedError(name)
	}
	fmt.Print(word)
	return nil
}

func (i *Interpreter) innerIf() error {
	cword := i.CWord
	ifRuntime := &Word{
		Name:        "if",
		IsPrimitive: true,
	}
	i.CWord.Compile(ifRuntime)
	// get `dest` from `then`
	resolveIf := func(dest int) {
		ifRuntime.PrimBody = func() error {
			// flag is stack top at run-time
			flag, err := i.DS.Pop()
			if err != nil {
				return err
			}
			if flag != 0 { // if stack top ≠ 0,
				// jump to `dest` (position of `then`)
				cword.pc = dest
			}
			return nil
		}
	}
	i.RS.Push(resolveIf)
	return nil
}

func (i *Interpreter) innerElse() error {
	// destination for if
	// location beyond `else`
	dest := len(i.CWord.Body)
	rs, err := i.RS.Pop()
	if err != nil {
		return UnstructuredError
	}
	resolveIf, ok := rs.(func(int))
	if !ok {
		return TypeError
	}
	resolveIf(dest)

	cword := i.CWord
	elseRuntime := &Word{
		Name:        "else",
		IsPrimitive: true,
	}
	i.CWord.Compile(elseRuntime)
	resolveElse := func(dest int) {
		elseRuntime.PrimBody = func() error { // run-time: jump to `then`
			cword.pc = dest
			return nil
		}
	}
	i.RS.Push(resolveElse)
	return nil
}

func (i *Interpreter) innerThen() error {
	// provide destination for if or else.
	// dest is location beyond `then`.
	dest := len(i.CWord.Body)
	rs, err := i.RS.Pop()
	if err != nil {
		return UnstructuredError
	}
	resolve, ok := rs.(func(int))
	if !ok {
		return TypeError
	}
	// resolve the branch originated by `if`, `else`...
	resolve(dest)
	i.CWord.Compile(&Word{
		Name:        "then",
		IsPrimitive: true,
		PrimBody: func() error { // run-time: do nothing.
			return nil
		},
	})
	return nil
}

func (i *Interpreter) immediate() error {
	if i.CWord == nil {
		return NoLastWordError
	}
	i.CWord.IsImmediate = true
	return nil
}

func (i *Interpreter) literal() error {
	n, err := i.DS.Pop()
	if err != nil {
		return err
	}
	num, ok := n.(int)
	if !ok {
		return TypeError
	}
	i.CompileNum(num)
	return nil
}

// ' <name>
// search the dictionary for name and leave its execution token on the stack.
func (i *Interpreter) tick() error {
	if !i.Scanner.Scan() {
		return EOFError
	}
	name := i.Scanner.Text()
	word, ok := i.Dictionary[name]
	if !ok {
		return UndefinedError(name)
	}
	i.DS.Push(word)
	return nil
}

// i*x xt - j*x
func (i *Interpreter) exec() error {
	xt, err := i.DS.Pop()
	if err != nil {
		return err
	}
	word, ok := xt.(*Word)
	if !ok {
		return TypeError
	}
	// ?????
	return i.Interpret(word)
}

func (i *Interpreter) bye() error {
	return QuitError
}
