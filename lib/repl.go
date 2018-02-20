package fifth

// Run repl.
import (
	"fmt"

	"github.com/chzyer/readline"
)

const (
	COLOR_GREEN = "\033[1;32m"
	COLOR_BLUE  = "\033[1;34m"

	COLOR_OFF = "\033[m"
)

func (i *Interpreter) prompt() string {
	if i.IsCompile {
		return fmt.Sprintf("%sc>%s ", COLOR_BLUE, COLOR_OFF)
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

// gomi
func (i *Interpreter) completer() *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItemDynamic(i.words,
			readline.PcItemDynamic(i.words,
				readline.PcItemDynamic(i.words,
					readline.PcItemDynamic(i.words,
						readline.PcItemDynamic(i.words,
							readline.PcItemDynamic(i.words,
								readline.PcItemDynamic(i.words,
									readline.PcItemDynamic(i.words,
										readline.PcItemDynamic(i.words,
											readline.PcItemDynamic(i.words,
												readline.PcItemDynamic(i.words))))))))))),
	)
}

func (i *Interpreter) Repl() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "\033[31m»\033[0m ",
		AutoComplete: i.completer(),
		EOFPrompt:    "\nbye!",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		rl.SetPrompt(i.prompt())

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
