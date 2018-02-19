package fifth

type Dictionary map[string]*Word

func (d *Dictionary) Set(name string, w *Word) {
	(*d)[name] = w
}

func (d *Dictionary) Get(name string) (*Word, bool) {
	word, ok := (*d)[name]
	return word, ok
}
