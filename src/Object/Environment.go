package Object

type Environment struct {
	Store     map[string]Object
	Immutable map[string]bool
	Outer     *Environment
}

func NewEnv() *Environment {
	return &Environment{
		Store:     make(map[string]Object),
		Immutable: make(map[string]bool),
		Outer:     nil,
	}
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

func (E *Environment) SetConst(Name string, Val Object) Object {
	E.Store[Name] = Val
	E.Immutable[Name] = true
	return Val
}

func (E *Environment) IsImmutable(Name string) bool {
	if E.Immutable[Name] {
		return true
	}
	if E.Outer != nil {
		return E.Outer.IsImmutable(Name)
	}
	return false
}