package fifth

// Run repl.
import (
	"fmt"

	"github.com/chzyer/readline"
)

const (
	COLOR_RED   = "\033[0;31m"
	COLOR_GREEN = "\033[1;32m"
	COLOR_BLUE  = "\033[1;34m"

	COLOR_OFF = "\033[m"
)

func (i *Interpreter) prompt() string {
	if i.IsCompile {
		return fmt.Sprintf("%scompile:%s ", COLOR_BLUE, COLOR_OFF)
	} else {
		return fmt.Sprintf("%s%d>%s ", COLOR_GREEN, len(i.DS.data), COLOR_OFF)
	}
}

func (i *Interpreter) words(string) []string {
	names := make([]string, 0)
	for name, _ := range i.Dictionary {
		names = append(names, name)
	}
	return names
}

func (i *Interpreter) Repl() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          i.prompt(),
		EOFPrompt:       "\nbye!",
		InterruptPrompt: fmt.Sprintf("\n%sinterrupt%s", COLOR_RED, COLOR_OFF),
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		rl.SetPrompt(i.prompt())

		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			i.IsCompile = false
			continue
		} else if err != nil {
			break
		}

		i.SetString(line)
		err = i.Run()
		if err == QuitError {
			return
		}
		if err != nil {
			fmt.Printf("%s%s%s\n", COLOR_RED, err.Error(), COLOR_OFF)
		}
	}
}
