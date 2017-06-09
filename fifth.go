package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Stack struct {
	data []int64
}

func (s *Stack) Push(elm int64) (int64, error) {
	s.data = append(s.data, elm)
	return elm, nil
}

func (s *Stack) Pop() (int64, error) {
	if len(s.data) == 0 {
		return 0, fmt.Errorf("stack underflow")
	}

	last := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return last, nil
}

func (s *Stack) Clear() {
	s.data = []int64{}
}

func (s *Stack) RunBinop(op func(n, m int64) int64) error {
	n, err := s.Pop()
	if err != nil {
		return err
	}
	m, err := s.Pop()
	if err != nil {
		return err
	}
	_, err = s.Push(op(n, m))
	return err
}

type Word struct {
	Name        string
	Immediate   bool
	IsPrimitive bool
	PrimBody    func() error
	Body        []*Word
}

type Dictionary map[string]*Word

func (i *Interpreter) SetPrimitive(name string, pbody func() error) {
	i.Dictionary[name] = &Word{
		Name:        name,
		IsPrimitive: true,
		PrimBody:    pbody,
	}
}

func (i *Interpreter) InitDictionary() {
	i.SetPrimitive("+", func() error {
		return i.Stack.RunBinop(func(n, m int64) int64 { return n + m })
	})
	i.SetPrimitive("-", func() error {
		return i.Stack.RunBinop(func(n, m int64) int64 { return m - n })
	})
	i.SetPrimitive("*", func() error {
		return i.Stack.RunBinop(func(n, m int64) int64 { return n * m })
	})
	i.SetPrimitive("/", func() error {
		return i.Stack.RunBinop(func(n, m int64) int64 { return m / n })
	})
	i.SetPrimitive(".", func() error {
		n, err := i.Stack.Pop()
		if err != nil {
			return err
		}
		fmt.Println(n)
		return nil
	})
	i.SetPrimitive("dup", func() error {
		n, err := i.Stack.Pop()
		if err != nil {
			return err
		}
		i.Stack.Push(n)
		i.Stack.Push(n)
		return nil
	})
	i.SetPrimitive(".s", func() error {
		fmt.Printf("<%d> ", len(i.Stack.data))
		for _, num := range i.Stack.data {
			fmt.Print(num, " ")
		}
		fmt.Println()
		return nil
	})
	i.SetPrimitive(":", func() error {
		i.IsCompile = true
		if !i.Scanner.Scan() {
			return i.Scanner.Err()
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
			return i.Scanner.Err()
		}
		name := i.Scanner.Text()
		word, ok := i.Dictionary[name]
		if !ok {
			return fmt.Errorf("Undefined word")
		}
		fmt.Println(":", name)
		fmt.Print(" ")
		for _, w := range word.Body {
			fmt.Printf(" %s", w.Name)
		}
		fmt.Println(" ;")
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

type Interpreter struct {
	Scanner    bufio.Scanner
	Stack      Stack
	Dictionary Dictionary
	IsCompile  bool
	CWord      *Word
	Error      error
}

func NewInterpreter() *Interpreter {
	i := new(Interpreter)
	i.Scanner = *bufio.NewScanner(os.Stdin)
	i.Scanner.Split(bufio.ScanWords)
	i.Dictionary = make(Dictionary)
	i.InitDictionary()
	return i
}

func (i *Interpreter) Execute(word *Word) error {
	if word.IsPrimitive {
		return word.PrimBody()
	}
	// compound word
	for _, w := range word.Body {
		if err := i.Execute(w); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Compile(word *Word) error {
	if word.Immediate {
		return i.Execute(word)
	}
	i.CWord.Body = append(i.CWord.Body, word)
	return nil
}

func (i *Interpreter) ReadToken() error {
	// can't read.
	if !i.Scanner.Scan() {
		return i.Scanner.Err()
	}
	text := i.Scanner.Text()
	if word, ok := i.Dictionary[text]; ok {
		if i.IsCompile {
			return i.Compile(word)
		} else {
			return i.Execute(word)
		}
	}
	num, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return err
	}
	if i.IsCompile {
		i.CWord.Body = append(i.CWord.Body, &Word{
			Name:        strconv.Itoa(int(num)),
			IsPrimitive: true,
			PrimBody: func() error {
				_, err := i.Stack.Push(num)
				return err
			},
		})
		return nil
	} else {
		_, err = i.Stack.Push(num)
		return err
	}
}

func main() {
	fmt.Println("こんにちは!")
	i := NewInterpreter()
	for {
		if err := i.ReadToken(); err != nil {
			fmt.Println("error:", err.Error())
		}
	}
}
