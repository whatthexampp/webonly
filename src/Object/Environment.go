package Object

type Environment struct {
	Store map[string]Object
	Outer *Environment
}

func NewEnv() *Environment {
	return &Environment{Store: make(map[string]Object), Outer: nil}
}

func NewEncEnv(Outer *Environment) *Environment {
	E := NewEnv()
	E.Outer = Outer
	return E
}

func (E *Environment) Get(Name string) (Object, bool) {
	Obj, Ok := E.Store[Name]
	if !Ok && E.Outer != nil {
		Obj, Ok = E.Outer.Get(Name)
	}
	return Obj, Ok
}

func (E *Environment) Set(Name string, Val Object) Object {
	E.Store[Name] = Val
	return Val
}