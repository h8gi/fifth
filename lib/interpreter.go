package fifth

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

type Interpreter struct {
	Scanner    bufio.Scanner
	DS         Stack // Data stack
	RS         Stack // Return stack??
	Dictionary Dictionary
	IsCompile  bool
	CWord      *Word // current word
	Error      error
}

func (i *Interpreter) SetPrimitive(name string, pbody func() error) {
	i.Dictionary.Set(name, &Word{
		Name:        name,
		IsPrimitive: true,
		PrimBody:    pbody,
	})
}

// Initialize interpreter's dictionary
func (i *Interpreter) InitDictionary() {
	i.SetPrimitive("+",
		i.DS.MakeBinFunc(func(n, m int) int { return n + m }))
	i.SetPrimitive("-",
		i.DS.MakeBinFunc(func(n, m int) int { return m - n }))
	i.SetPrimitive("*",
		i.DS.MakeBinFunc(func(n, m int) int { return n * m }))
	i.SetPrimitive("/",
		i.DS.MakeBinFunc(func(n, m int) int { return m / n }))
	i.SetPrimitive(".", func() error {
		n, err := i.DS.Pop()
		if err != nil {
			return err
		}
		fmt.Println(n)
		return nil
	})
	i.SetPrimitive("dup", func() error {
		n, err := i.DS.Pop()
		if err != nil {
			return err
		}
		i.DS.Push(n)
		i.DS.Push(n)
		return nil
	})
	i.SetPrimitive(".s", func() error {
		fmt.Println(i.DS.data)
		return nil
	})
	i.SetPrimitive(".r", func() error {
		fmt.Println(i.RS.data)
		return nil
	})
	i.SetPrimitive(":", func() error {
		i.IsCompile = true
		if !i.Scanner.Scan() {
			return EOFError
		}
		i.CWord = &Word{
			Name: i.Scanner.Text(),
		}
		return nil
	})
	i.SetPrimitive(";", func() error {
		i.IsCompile = false
		i.Dictionary[i.CWord.Name] = i.CWord
		return nil
	})
	i.Dictionary[";"].IsImmediate = true
	i.Dictionary[";"].IsCompileOnly = true

	i.SetPrimitive("see", func() error {
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
	})
	// compilation behavior of if
	i.SetPrimitive("if", func() error {
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
	})
	i.Dictionary["if"].IsImmediate = true
	i.Dictionary["if"].IsCompileOnly = true

	i.SetPrimitive("else", func() error {
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
	})
	i.Dictionary["else"].IsImmediate = true
	i.Dictionary["else"].IsCompileOnly = true

	i.SetPrimitive("then", func() error {
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
	})
	i.Dictionary["then"].IsImmediate = true
	i.Dictionary["then"].IsCompileOnly = true

	i.SetPrimitive("immediate", func() error {
		if i.CWord == nil {
			return NoLastWordError
		}
		i.CWord.IsImmediate = true
		return nil
	})

	i.SetPrimitive("literal", func() error {
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
	})
	i.Dictionary["literal"].IsImmediate = true
	i.Dictionary["literal"].IsCompileOnly = true

	i.SetPrimitive("'", func() error {
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
	})

	i.SetPrimitive("bye", func() error {
		return QuitError
	})

}

func NewInterpreter() *Interpreter {
	i := new(Interpreter)
	i.Scanner = *bufio.NewScanner(os.Stdin)
	i.Scanner.Split(bufio.ScanWords)
	i.Dictionary = make(Dictionary)
	i.InitDictionary()
	return i
}

// Execute interpretation behavior fo word
func (i *Interpreter) Interpret(word *Word) error {
	if word.IsPrimitive {
		return word.PrimBody()
	}

	// compound word
	for word.pc = 0; word.pc < len(word.Body); word.pc += 1 {
		w := word.Body[word.pc]
		if err := i.Interpret(w); err != nil {
			return err
		}
	}
	return nil
}

// Execute compilation behavior of word
func (i *Interpreter) Compile(word *Word) error {
	if word.IsImmediate {
		return i.Interpret(word)
	}
	i.CWord.Compile(word)
	return nil
}

// Compile literal
func (i *Interpreter) CompileNum(num int) {
	w := &Word{
		Name:        strconv.Itoa(int(num)),
		IsPrimitive: true,
		PrimBody: func() error {
			i.DS.Push(num)
			return nil
		},
	}
	i.CWord.Compile(w)
}

func (i *Interpreter) Execute(text string) error {
	// Look text up in the dictionary.
	if word, ok := i.Dictionary.Get(text); ok {
		if i.IsCompile {
			return i.Compile(word)
		} else {
			if word.IsCompileOnly {
				return CompileOnlyError(word.Name)
			}
			return i.Interpret(word)
		}
	}
	// read text as number.
	num, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return UndefinedError(text)
	}

	if i.IsCompile {
		i.CompileNum(int(num))
	} else {
		i.DS.Push(int(num))
	}
	return nil
}

func (i *Interpreter) Abort() {
	i.DS.Clear()
	i.RS.Clear()
	i.IsCompile = false
	i.CWord = nil
}

// Execute until EOF
func (i *Interpreter) Run() error {
	for {
		// EOF
		if !i.Scanner.Scan() {
			fmt.Println(" ok")
			return nil
		}

		text := i.Scanner.Text()
		err := i.Execute(text)
		if err != nil {
			i.Abort()
			fmt.Println(err.Error())
			return err
		}
	}
}

// Set input source to string
func (i *Interpreter) SetString(s string) {
	i.Scanner = *bufio.NewScanner(strings.NewReader(s))
	i.Scanner.Split(bufio.ScanWords)
}

func (i *Interpreter) Repl() {
	rl, err := readline.New(fmt.Sprintf("%d> ", len(i.DS.data)))
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	for {
		if i.IsCompile {
			rl.SetPrompt("compile> ")
		} else {
			rl.SetPrompt(fmt.Sprintf("%d> ", len(i.DS.data)))
		}

		line, err := rl.Readline()
		if err != nil {
			break
		}
		i.SetString(line)
		if err := i.Run(); err == QuitError {
			return
		}
	}
}
