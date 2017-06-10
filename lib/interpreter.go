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
	S          Stack // Data stack
	R          Stack // Return stack??
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
		i.S.MakeBinFunc(func(n, m int64) int64 { return n + m }))
	i.SetPrimitive("-",
		i.S.MakeBinFunc(func(n, m int64) int64 { return m - n }))
	i.SetPrimitive("*",
		i.S.MakeBinFunc(func(n, m int64) int64 { return n * m }))
	i.SetPrimitive("/",
		i.S.MakeBinFunc(func(n, m int64) int64 { return m / n }))
	i.SetPrimitive(".", func() error {
		n, err := i.S.Pop()
		if err != nil {
			return err
		}
		fmt.Println(n)
		return nil
	})
	i.SetPrimitive("dup", func() error {
		n, err := i.S.Pop()
		if err != nil {
			return err
		}
		i.S.Push(n)
		i.S.Push(n)
		return nil
	})
	i.SetPrimitive(".s", func() error {
		for _, num := range i.S.data {
			fmt.Print(num, " ")
		}
		return nil
	})
	i.SetPrimitive(":", func() error {
		i.IsCompile = true
		if !i.Scanner.Scan() {
			return EOFError
		}
		name := i.Scanner.Text()
		word := Word{
			Name: name,
		}
		i.CWord = &word
		return nil
	})
	i.SetPrimitive(";", func() error {
		i.IsCompile = false
		i.Dictionary[i.CWord.Name] = i.CWord
		return nil
	})
	i.Dictionary[";"].Immediate = true

	i.SetPrimitive("see", func() error {
		if !i.Scanner.Scan() {
			return EOFError
		}
		name := i.Scanner.Text()
		word, ok := i.Dictionary[name]
		if !ok {
			return UndefinedError(name)
		}
		if word.IsPrimitive {
			fmt.Printf("primitive word: %q ", name)
			return nil
		}

		fmt.Println(":", name)
		fmt.Print(" ")
		for _, w := range word.Body {
			fmt.Printf(" %s", w.Name)
		}
		fmt.Print(" ; ")
		if word.Immediate {
			fmt.Print("immediate ")
		}
		return nil
	})

	i.SetPrimitive("immediate", func() error {
		i.CWord.Immediate = true
		return nil
	})

	i.Dictionary["square"] = &Word{
		Body: []*Word{i.Dictionary["dup"], i.Dictionary["*"]},
	}
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
	for _, w := range word.Body {
		if err := i.Interpret(w); err != nil {
			return err
		}
	}
	return nil
}

// Execute compilation behavior of word
func (i *Interpreter) Compile(word *Word) error {
	if word.Immediate {
		return i.Interpret(word)
	}
	i.CWord.Body = append(i.CWord.Body, word)
	return nil
}

// Compile literal
func (i *Interpreter) CompileNum(num int64) {
	i.CWord.Body = append(i.CWord.Body, &Word{
		Name:        strconv.Itoa(int(num)),
		IsPrimitive: true,
		PrimBody: func() error {
			i.S.Push(num)
			return nil
		},
	})
}

func (i *Interpreter) Execute(text string) error {
	// Look text up in the dictionary.
	if word, ok := i.Dictionary[text]; ok {
		if i.IsCompile {
			return i.Compile(word)
		} else {
			return i.Interpret(word)
		}
	}
	// read text as number.
	num, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return UndefinedError(text)
	}

	if i.IsCompile {
		i.CompileNum(num)
	} else {
		i.S.Push(num)
	}
	return nil
}

func (i *Interpreter) Abort() {
	i.S.Clear()
	i.R.Clear()
	i.IsCompile = false
	i.CWord = nil
}

// Execute until EOF
func (i *Interpreter) Run() {
	for {
		// EOF
		if !i.Scanner.Scan() {
			fmt.Println("ok")
			return
		}

		text := i.Scanner.Text()
		err := i.Execute(text)
		if err != nil {
			i.Abort()
			fmt.Println(err.Error())
			return
		}
	}
}

// Set input source to string
func (i *Interpreter) Setstring(s string) {
	i.Scanner = *bufio.NewScanner(strings.NewReader(s))
	i.Scanner.Split(bufio.ScanWords)
}

func (i *Interpreter) Repl() {
	rl, err := readline.New(fmt.Sprintf("<%d> ", len(i.S.data)))
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	for {
		rl.SetPrompt(fmt.Sprintf("<%d> ", len(i.S.data)))
		line, err := rl.Readline()
		if err != nil {
			break
		}
		i.Setstring(line)
		i.Run()
	}
}
