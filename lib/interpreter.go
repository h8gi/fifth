package fifth

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

type Dictionary map[string]*Word

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
	i.Dictionary[name] = &Word{
		Name:        name,
		IsPrimitive: true,
		PrimBody:    pbody,
	}
}

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
		orig := len(i.CWord.Body)
		word := i.CWord
		i.RS.Push(orig)
		// compilation behavior
		i.CWord.Compile(&Word{
			Name:        "if",
			IsPrimitive: true,
			// Update runtime behavior by called from then
			PrimBody: func() error {
				dest, err := i.RS.Pop() // `then` position
				if err != nil {
					return err
				}
				// runtime is the interpratation behavior of if
				runtime := func() error {
					flag, err := i.DS.Pop()
					if err != nil {
						return err
					}
					if flag != 0 { // if not (stack top neq zero)
						// jump to destination
						pc, ok := dest.(int)
						if !ok {
							return TypeError
						}
						word.pc = pc
					}
					return nil
				}
				// update code to `runtime`
				i.CWord.Body[orig].PrimBody = runtime
				return nil
			},
		})
		return nil
	})
	i.Dictionary["if"].IsImmediate = true
	i.Dictionary["if"].IsCompileOnly = true

	i.SetPrimitive("then", func() error {
		o, err := i.RS.Pop()
		if err != nil {
			return UnstructuredError
		}
		orig, ok := o.(int)
		if !ok {
			return TypeError
		}

		dest := len(i.CWord.Body)
		i.RS.Push(dest)
		i.CWord.Compile(&Word{
			Name:        "then",
			IsPrimitive: true,
			PrimBody: func() error {
				return nil
			},
		})
		branch := i.CWord.Body[orig]
		return branch.PrimBody()
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
	if word, ok := i.Dictionary[text]; ok {
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
func (i *Interpreter) Setstring(s string) {
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
		i.Setstring(line)
		if err := i.Run(); err == QuitError {
			return
		}
	}
}
