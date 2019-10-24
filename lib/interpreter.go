package fifth

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Interpreter struct {
	Scanner    bufio.Scanner
	DS         Stack // Data stack
	RS         Stack // Return stack??
	Dictionary Dictionary
	IsCompile  bool
	CWord      *Word // the pointer to the word being compiled
	Error      error
}

func (i *Interpreter) SetPrimitive(name string, pbody PrimBody) {
	i.Dictionary.Set(name, &Word{
		Name:        name,
		IsPrimitive: true,
		PrimBody:    pbody,
	})
}

// Initialize interpreter's dictionary

func NewInterpreter() *Interpreter {
	i := new(Interpreter)
	i.Scanner = *bufio.NewScanner(os.Stdin)
	i.Scanner.Split(bufio.ScanWords)
	i.Dictionary = make(Dictionary)
	i.LoadPrimitives()
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
	// if the word is immediate, execute normal behavior
	if word.IsImmediate {
		return i.Interpret(word)
	}
	// execute compilation behavior of word
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

// eval string and execute it
func (i *Interpreter) EvalToken(token string) error {
	// Look token up in the dictionary.
	if word, ok := i.Dictionary.Get(token); ok {
		if i.IsCompile {
			return i.Compile(word)
		} else {
			if word.IsCompileOnly {
				return CompileOnlyError(word.Name)
			}
			return i.Interpret(word)
		}
	}
	// read token as number.
	num, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return UndefinedError(token)
	}

	if i.IsCompile {
		i.CompileNum(int(num))
	} else {
		i.DS.Push(int(num))
	}
	return nil
}

// Reset the stacks and interpreter.
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
			return nil
		}

		token := i.Scanner.Text()
		err := i.EvalToken(token)
		if err != nil {
			i.Abort()
			// fmt.Println(err.Error())
			return err
		}
	}
}

// Set input source to string
func (i *Interpreter) SetString(s string) {
	i.Scanner = *bufio.NewScanner(strings.NewReader(s))
	i.Scanner.Split(bufio.ScanWords)
}
