package Object

import (
	"fmt"
	"strconv"
	"strings"
	"webonly/src/Ast"
)

type ObjType string

const (
	NumObj   ObjType = "Num"
	StrObj   ObjType = "Str"
	BoolObj  ObjType = "Bool"
	NullObj  ObjType = "Null"
	RetObj   ObjType = "Ret"
	ErrObj   ObjType = "Err"
	FuncObj  ObjType = "Func"
	ArrObj   ObjType = "Array"
	ClassObj ObjType = "Class"
	InstObj  ObjType = "Instance"
	BoundObj ObjType = "BoundMethod"
	BltObj   ObjType = "Builtin"
	EnumObj  ObjType = "Enum"
	ModObj   ObjType = "Module"
)

type Object interface {
	Type() ObjType
	Inspect() string
}

type Num struct{ Value float64 }

func (N *Num) Type() ObjType {
	return NumObj
}

func (N *Num) Inspect() string {
	return strconv.FormatFloat(N.Value, 'f', -1, 64)
}

type Str struct{ Value string }

func (S *Str) Type() ObjType {
	return StrObj
}

func (S *Str) Inspect() string {
	return S.Value
}

type Bool struct{ Value bool }

func (B *Bool) Type() ObjType {
	return BoolObj
}

func (B *Bool) Inspect() string {
	return fmt.Sprintf("%t", B.Value)
}

type Null struct{}

func (N *Null) Type() ObjType {
	return NullObj
}

func (N *Null) Inspect() string {
	return "Null"
}

type Ret struct{ Value Object }

func (R *Ret) Type() ObjType {
	return RetObj
}

func (R *Ret) Inspect() string {
	return R.Value.Inspect()
}

type Err struct{ Msg string }

func (E *Err) Type() ObjType {
	return ErrObj
}

func (E *Err) Inspect() string {
	return "ERROR: " + E.Msg
}

type Func struct {
	Params []*Ast.Ident
	Body   *Ast.BlockStmt
	Env    *Environment
}

func (F *Func) Type() ObjType {
	return FuncObj
}

func (F *Func) Inspect() string {
	P := []string{}
	for _, Param := range F.Params {
		P = append(P, Param.String())
	}
	return "function(" + strings.Join(P, ", ") + "): ... end;"
}

type Array struct{ Elems []Object }

func (A *Array) Type() ObjType {
	return ArrObj
}

func (A *Array) Inspect() string {
	E := []string{}
	for _, Val := range A.Elems {
		E = append(E, Val.Inspect())
	}
	return "[" + strings.Join(E, ", ") + "]"
}

type Class struct {
	Name    string
	Parent  *Class
	Methods map[string]*Func
}

func (C *Class) Type() ObjType {
	return ClassObj
}

func (C *Class) Inspect() string {
	return "class " + C.Name
}

type Enum struct {
	Name  string
	Cases map[string]Object
}

func (E *Enum) Type() ObjType {
	return EnumObj
}

func (E *Enum) Inspect() string {
	return "enum " + E.Name
}

type Instance struct {
	Cls    *Class
	Fields map[string]Object
}

func (I *Instance) Type() ObjType {
	return InstObj
}

func (I *Instance) Inspect() string {
	return "Instance of " + I.Cls.Name
}

type BoundMethod struct {
	Self   *Instance
	Method *Func
}

func (B *BoundMethod) Type() ObjType {
	return BoundObj
}

func (B *BoundMethod) Inspect() string {
	return "Bound Method of " + B.Self.Cls.Name
}

type BuiltinFn func(Args ...Object) Object

type Builtin struct{ Fn BuiltinFn }

func (B *Builtin) Type() ObjType {
	return BltObj
}

func (B *Builtin) Inspect() string {
	return "Builtin Function"
}

type Module struct {
	Name    string
	Exports map[string]Object
}

func (M *Module) Type() ObjType {
	return ModObj
}

func (M *Module) Inspect() string {
	return "Module(" + M.Name + ")"
}