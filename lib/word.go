package fifth

type Word struct {
	Name        string
	Immediate   bool
	IsPrimitive bool
	PrimBody    func() error
	Body        []*Word
}
