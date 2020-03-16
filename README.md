# fifth

Simple forth interpreter written in go

## Usage

You can use fifth as a repl or a library.

### as a repl

```shell
go get github.com/h8gi/fifth
./fifth
```

### as a library

```go
package main

import (
	fifth "github.com/h8gi/fifth/lib"
)

func main() {
	i := fifth.NewInterpreter()
	i.SetReader(os.Stdin)
	i.SetWriter(os.Stdout)
	i.Repl()
}
```

## Primitive words

see [primitives.go](https://github.com/h8gi/fifth/blob/master/lib/primitives.go)

## Todo

- Get values from the interpreter
- Send values to the interpreter
- Eval text
