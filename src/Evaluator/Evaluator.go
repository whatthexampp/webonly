package Evaluator

import (
	"fmt"
	"os"
	"path/filepath"
	"webonly/src/Ast"
	"webonly/src/Lexer"
	"webonly/src/Object"
	"webonly/src/Parser"
)

var (
	TrueObj  = &Object.Bool{Value: true}
	FalseObj = &Object.Bool{Value: false}
	NullObj  = &Object.Null{}
)

func Eval(Node Ast.Node, Env *Object.Environment) Object.Object {
	switch N := Node.(type) {
	case *Ast.Program:
		return EvalProg(N, Env)
	case *Ast.HtmlStmt:
		return EvalHtml(N, Env)
	case *Ast.ClassStmt:
		return EvalClass(N, Env)
	case *Ast.WhileStmt:
		return EvalWhile(N, Env)
	case *Ast.ConstStmt:
		return EvalConst(N, Env)
	case *Ast.EnumStmt:
		return EvalEnum(N, Env)
	case *Ast.PublicStmt:
		return EvalPublic(N, Env)
	case *Ast.ExprStmt:
		return Eval(N.Expr, Env)
	case *Ast.BlockStmt:
		return EvalBlock(N, Env)
	case *Ast.RetStmt:
		Val := Eval(N.Value, Env)
		if IsErr(Val) {
			return Val
		}
		return &Object.Ret{Value: Val}
	case *Ast.AssignExpr:
		return EvalAssign(N, Env)
	case *Ast.NumLit:
		return &Object.Num{Value: N.Value}
	case *Ast.StrLit:
		return &Object.Str{Value: N.Value}
	case *Ast.NullLit:
		return NullObj
	case *Ast.ArrayLit:
		return EvalArray(N, Env)
	case *Ast.Bool:
		return NatToBool(N.Value)
	case *Ast.PrefixExpr:
		Right := Eval(N.Right, Env)
		if IsErr(Right) {
			return Right
		}
		return EvalPre(N.Op, Right, N.Token.Line)
	case *Ast.InfixExpr:
		Left := Eval(N.Left, Env)
		if IsErr(Left) {
			return Left
		}
		if N.Op == "&&" {
			if !IsTruthy(Left) {
				return Left
			}
			Right := Eval(N.Right, Env)
			if IsErr(Right) {
				return Right
			}
			return Right
		}
		if N.Op == "||" {
			if IsTruthy(Left) {
				return Left
			}
			Right := Eval(N.Right, Env)
			if IsErr(Right) {
				return Right
			}
			return Right
		}
		Right := Eval(N.Right, Env)
		if IsErr(Right) {
			return Right
		}
		return EvalInfix(N.Op, Left, Right, N.Token.Line)
	case *Ast.IndexExpr:
		Left := Eval(N.Left, Env)
		if IsErr(Left) {
			return Left
		}
		Index := Eval(N.Index, Env)
		if IsErr(Index) {
			return Index
		}
		return EvalIndex(Left, Index, N.Token.Line)
	case *Ast.DotExpr:
		Left := Eval(N.Left, Env)
		if IsErr(Left) {
			return Left
		}
		return EvalDot(Left, N.Right.Value, N.Token.Line)
	case *Ast.NewExpr:
		return EvalNew(N, Env)
	case *Ast.IfExpr:
		return EvalIf(N, Env)
	case *Ast.Ident:
		return EvalIdent(N, Env)
	case *Ast.FuncLit:
		FnObj := &Object.Func{Params: N.Params, Env: Env, Body: N.Body}
		if N.Name != "" {
			Env.SetConst(N.Name, FnObj)
		}
		return FnObj
	case *Ast.CallExpr:
		Func := Eval(N.Func, Env)
		if IsErr(Func) {
			return Func
		}
		Args := EvalExprs(N.Args, Env)
		if len(Args) == 1 && IsErr(Args[0]) {
			return Args[0]
		}
		return Apply(Func, Args, N.Token.Line)
	case *Ast.ImportExpr:
		return EvalImport(N, Env)
	}
	return nil
}

func EvalProg(P *Ast.Program, Env *Object.Environment) Object.Object {
	var Res Object.Object
	for _, Stmt := range P.Stmts {
		Res = Eval(Stmt, Env)
		switch R := Res.(type) {
		case *Object.Ret:
			return R.Value
		case *Object.Err:
			return R
		}
	}
	return Res
}

func EvalBlock(B *Ast.BlockStmt, Env *Object.Environment) Object.Object {
	var Res Object.Object
	for _, Stmt := range B.Stmts {
		Res = Eval(Stmt, Env)
		if Res != nil {
			RT := Res.Type()
			if RT == Object.RetObj || RT == Object.ErrObj {
				return Res
			}
		}
	}
	return Res
}

func EvalWhile(W *Ast.WhileStmt, Env *Object.Environment) Object.Object {
	var Res Object.Object
	for {
		Cond := Eval(W.Cond, Env)
		if IsErr(Cond) {
			return Cond
		}
		if !IsTruthy(Cond) {
			break
		}
		Res = Eval(W.Body, Env)
		if Res != nil {
			RT := Res.Type()
			if RT == Object.RetObj || RT == Object.ErrObj {
				return Res
			}
		}
	}
	if Res == nil {
		return NullObj
	}
	return Res
}

func EvalClass(C *Ast.ClassStmt, Env *Object.Environment) Object.Object {
	var ParentCls *Object.Class
	if C.Parent != nil {
		ParentObj, Ok := Env.Get(C.Parent.Value)
		if !Ok {
			return NewErr(C.Token.Line, "Parent class not found: %s", C.Parent.Value)
		}
		if Cls, IsCls := ParentObj.(*Object.Class); IsCls {
			ParentCls = Cls
		} else {
			return NewErr(C.Token.Line, "Not a class: %s", C.Parent.Value)
		}
	}

	Methods := make(map[string]*Object.Func)
	for _, M := range C.Methods {
		Methods[M.Name] = &Object.Func{Params: M.Params, Env: Env, Body: M.Body}
	}
	Cls := &Object.Class{Name: C.Name.Value, Parent: ParentCls, Methods: Methods}
	Env.SetConst(C.Name.Value, Cls)
	return Cls
}

func EvalConst(C *Ast.ConstStmt, Env *Object.Environment) Object.Object {
	Val := Eval(C.Value, Env)
	if IsErr(Val) {
		return Val
	}
	Env.SetConst(C.Name.Value, Val)
	return NullObj
}

func EvalEnum(E *Ast.EnumStmt, Env *Object.Environment) Object.Object {
	Cases := make(map[string]Object.Object)
	for Idx, Case := range E.Cases {
		Cases[Case.Value] = &Object.Num{Value: float64(Idx)}
	}
	EnumObj := &Object.Enum{Name: E.Name.Value, Cases: Cases}
	Env.SetConst(E.Name.Value, EnumObj)
	return NullObj
}

func EvalPublic(P *Ast.PublicStmt, Env *Object.Environment) Object.Object {
	Res := Eval(P.Stmt, Env)
	if IsErr(Res) {
		return Res
	}
	Names := ExtractPublicNames(P.Stmt)
	for _, Name := range Names {
		Env.SetPublic(Name)
	}
	return Res
}

func ExtractPublicNames(Node Ast.Node) []string {
	switch N := Node.(type) {
	case *Ast.ExprStmt:
		if Asgn, Ok := N.Expr.(*Ast.AssignExpr); Ok {
			if Id, Ok := Asgn.Left.(*Ast.Ident); Ok {
				return []string{Id.Value}
			}
		}
		if Fn, Ok := N.Expr.(*Ast.FuncLit); Ok && Fn.Name != "" {
			return []string{Fn.Name}
		}
	case *Ast.ClassStmt:
		return []string{N.Name.Value}
	case *Ast.ConstStmt:
		return []string{N.Name.Value}
	case *Ast.EnumStmt:
		return []string{N.Name.Value}
	}
	return nil
}

func EvalAssign(A *Ast.AssignExpr, Env *Object.Environment) Object.Object {
	Val := Eval(A.Value, Env)
	if IsErr(Val) {
		return Val
	}
	switch L := A.Left.(type) {
	case *Ast.Ident:
		if Env.IsImmutable(L.Value) {
			return NewErr(A.Token.Line, "Cannot assign to constant: %s", L.Value)
		}
		Env.Set(L.Value, Val)
		return Val
	case *Ast.IndexExpr:
		ArrObj := Eval(L.Left, Env)
		IdxObj := Eval(L.Index, Env)
		if Arr, Ok := ArrObj.(*Object.Array); Ok {
			if Idx, Ok := IdxObj.(*Object.Num); Ok {
				I := int(Idx.Value)
				if I >= 0 && I < len(Arr.Elems) {
					Arr.Elems[I] = Val
					return Val
				} else if I == len(Arr.Elems) {
					Arr.Elems = append(Arr.Elems, Val)
					return Val
				}
			}
		}
		return NewErr(A.Token.Line, "Invalid array assignment")
	case *Ast.DotExpr:
		Obj := Eval(L.Left, Env)
		if Inst, Ok := Obj.(*Object.Instance); Ok {
			Inst.Fields[L.Right.Value] = Val
			return Val
		}
		return NewErr(A.Token.Line, "Invalid property assignment")
	}
	return NewErr(A.Token.Line, "Invalid assignment left-hand side")
}

func EvalArray(A *Ast.ArrayLit, Env *Object.Environment) Object.Object {
	Elems := EvalExprs(A.Elems, Env)
	if len(Elems) == 1 && IsErr(Elems[0]) {
		return Elems[0]
	}
	return &Object.Array{Elems: Elems}
}

func EvalIndex(Left Object.Object, Index Object.Object, Line int) Object.Object {
	if Arr, Ok := Left.(*Object.Array); Ok {
		if Idx, Ok := Index.(*Object.Num); Ok {
			I := int(Idx.Value)
			if I >= 0 && I < len(Arr.Elems) {
				return Arr.Elems[I]
			}
			return NullObj
		}
	}
	return NewErr(Line, "Index operator not supported for %s", Left.Type())
}

func EvalDot(Left Object.Object, Prop string, Line int) Object.Object {
	if Mod, Ok := Left.(*Object.Module); Ok {
		if Val, Ok := Mod.Exports[Prop]; Ok {
			return Val
		}
		return NewErr(Line, "Export not found: %s", Prop)
	}
	if Inst, Ok := Left.(*Object.Instance); Ok {
		if Val, Ok := Inst.Fields[Prop]; Ok {
			return Val
		}
		Cls := Inst.Cls
		for Cls != nil {
			if Method, Ok := Cls.Methods[Prop]; Ok {
				return &Object.BoundMethod{Self: Inst, Method: Method}
			}
			Cls = Cls.Parent
		}
		return NullObj
	}
	if Enm, Ok := Left.(*Object.Enum); Ok {
		if Val, Ok := Enm.Cases[Prop]; Ok {
			return Val
		}
		return NewErr(Line, "Enum case not found: %s", Prop)
	}
	return NewErr(Line, "Property access not supported for %s", Left.Type())
}

func EvalNew(N *Ast.NewExpr, Env *Object.Environment) Object.Object {
	ClsObj := Eval(N.Class, Env)
	if IsErr(ClsObj) {
		return ClsObj
	}
	Cls, Ok := ClsObj.(*Object.Class)
	if !Ok {
		return NewErr(N.Token.Line, "Not a class: %s", N.Class.String())
	}
	Inst := &Object.Instance{Cls: Cls, Fields: make(map[string]Object.Object)}

	ConstructCls := Cls
	var ConstructMethod *Object.Func
	for ConstructCls != nil {
		if Method, Ok := ConstructCls.Methods["Construct"]; Ok {
			ConstructMethod = Method
			break
		}
		ConstructCls = ConstructCls.Parent
	}

	if ConstructMethod != nil {
		Args := EvalExprs(N.Args, Env)
		if len(Args) == 1 && IsErr(Args[0]) {
			return Args[0]
		}
		Apply(&Object.BoundMethod{Self: Inst, Method: ConstructMethod}, Args, N.Token.Line)
	}
	return Inst
}

func EvalHtml(H *Ast.HtmlStmt, Env *Object.Environment) Object.Object {
	if OutObj, Ok := Env.Get("__webonly_html_out"); Ok {
		if BuiltinOut, IsBuiltin := OutObj.(*Object.Builtin); IsBuiltin {
			BuiltinOut.Fn(&Object.Str{Value: H.Value})
		}
	}
	return &Object.Str{Value: H.Value}
}

func EvalImport(I *Ast.ImportExpr, Env *Object.Environment) Object.Object {
	PathObj := Eval(I.Path, Env)
	if IsErr(PathObj) {
		return PathObj
	}
	if PathObj.Type() != Object.StrObj {
		return NewErr(I.Token.Line, "Import path must be a string")
	}

	RelPath := PathObj.(*Object.Str).Value

	var CurrentFile string
	if FileObj, Ok := Env.Get("__FILE__"); Ok && FileObj.Type() == Object.StrObj {
		CurrentFile = FileObj.(*Object.Str).Value
	}

	var TargetPath string
	if filepath.IsAbs(RelPath) {
		TargetPath = RelPath
	} else if CurrentFile != "" {
		TargetPath = filepath.Join(filepath.Dir(CurrentFile), RelPath)
	} else {
		TargetPath = RelPath
	}

	TargetPath = filepath.Clean(TargetPath)

	GlobalEnv := Env
	for GlobalEnv.Outer != nil {
		GlobalEnv = GlobalEnv.Outer
	}

	if Cached, Ok := GlobalEnv.Modules[TargetPath]; Ok {
		return Cached
	}

	Src, Err := os.ReadFile(TargetPath)
	if Err != nil {
		return NewErr(I.Token.Line, "Failed to import file: %s", TargetPath)
	}

	Lex := Lexer.Create(string(Src))
	Par := Parser.Create(Lex)
	Prog := Par.ParseProg()

	if len(Par.Errs) > 0 {
		return NewErr(I.Token.Line, "Parse error in imported file %s: %s", TargetPath, Par.Errs[0])
	}

	GlobalEnv.Modules[TargetPath] = NullObj

	ModEnv := Object.NewEncEnv(GlobalEnv)
	ModEnv.SetConst("__FILE__", &Object.Str{Value: TargetPath})

	Res := EvalProg(Prog, ModEnv)

	if IsErr(Res) {
		delete(GlobalEnv.Modules, TargetPath)
		return Res
	}

	Exports := make(map[string]Object.Object)
	for Name := range ModEnv.Publics {
		if Val, Ok := ModEnv.Get(Name); Ok {
			Exports[Name] = Val
		}
	}
	Mod := &Object.Module{Name: TargetPath, Exports: Exports}

	GlobalEnv.Modules[TargetPath] = Mod
	return Mod
}

func NatToBool(In bool) *Object.Bool {
	if In {
		return TrueObj
	}
	return FalseObj
}

func EvalPre(Op string, Right Object.Object, Line int) Object.Object {
	switch Op {
	case "!":
		return EvalBang(Right)
	case "-":
		return EvalMinus(Right, Line)
	default:
		return NewErr(Line, "Unknown op: %s%s", Op, Right.Type())
	}
}

func EvalBang(Right Object.Object) Object.Object {
	switch Right {
	case TrueObj:
		return FalseObj
	case FalseObj:
		return TrueObj
	case NullObj:
		return TrueObj
	default:
		return FalseObj
	}
}

func EvalMinus(Right Object.Object, Line int) Object.Object {
	if Right.Type() != Object.NumObj {
		return NewErr(Line, "Unknown op: -%s", Right.Type())
	}
	return &Object.Num{Value: -Right.(*Object.Num).Value}
}

func EvalInfix(Op string, Left Object.Object, Right Object.Object, Line int) Object.Object {
	if Op == "==" {
		if Left.Type() != Right.Type() {
			return FalseObj
		}
		if Left.Type() == Object.NumObj {
			return NatToBool(Left.(*Object.Num).Value == Right.(*Object.Num).Value)
		}
		if Left.Type() == Object.StrObj {
			return NatToBool(Left.(*Object.Str).Value == Right.(*Object.Str).Value)
		}
		return NatToBool(Left == Right)
	}
	if Op == "!=" {
		if Left.Type() != Right.Type() {
			return TrueObj
		}
		if Left.Type() == Object.NumObj {
			return NatToBool(Left.(*Object.Num).Value != Right.(*Object.Num).Value)
		}
		if Left.Type() == Object.StrObj {
			return NatToBool(Left.(*Object.Str).Value != Right.(*Object.Str).Value)
		}
		return NatToBool(Left != Right)
	}
	if Left.Type() == Object.NumObj && Right.Type() == Object.NumObj {
		return EvalNumInfix(Op, Left, Right, Line)
	}
	if Left.Type() == Object.StrObj && Right.Type() == Object.StrObj {
		return EvalStrInfix(Op, Left, Right, Line)
	}
	if Op == "+" && (Left.Type() == Object.StrObj || Right.Type() == Object.StrObj) {
		return &Object.Str{Value: Left.Inspect() + Right.Inspect()}
	}
	if Left.Type() != Right.Type() {
		return NewErr(Line, "Type mismatch: %s %s %s", Left.Type(), Op, Right.Type())
	}
	return NewErr(Line, "Unknown op: %s %s %s", Left.Type(), Op, Right.Type())
}

func EvalNumInfix(Op string, Left Object.Object, Right Object.Object, Line int) Object.Object {
	LV := Left.(*Object.Num).Value
	RV := Right.(*Object.Num).Value
	switch Op {
	case "+":
		return &Object.Num{Value: LV + RV}
	case "-":
		return &Object.Num{Value: LV - RV}
	case "*":
		return &Object.Num{Value: LV * RV}
	case "/":
		if RV == 0 {
			return NewErr(Line, "Division by zero")
		}
		return &Object.Num{Value: LV / RV}
	case "%":
		if RV == 0 {
			return NewErr(Line, "Modulo by zero")
		}
		return &Object.Num{Value: float64(int(LV) % int(RV))}
	case "<":
		return NatToBool(LV < RV)
	case ">":
		return NatToBool(LV > RV)
	default:
		return NewErr(Line, "Unknown op: %s %s %s", Left.Type(), Op, Right.Type())
	}
}

func EvalStrInfix(Op string, Left Object.Object, Right Object.Object, Line int) Object.Object {
	if Op != "+" {
		return NewErr(Line, "Unknown op: %s %s %s", Left.Type(), Op, Right.Type())
	}
	return &Object.Str{Value: Left.(*Object.Str).Value + Right.(*Object.Str).Value}
}

func EvalIf(I *Ast.IfExpr, Env *Object.Environment) Object.Object {
	Cond := Eval(I.Cond, Env)
	if IsErr(Cond) {
		return Cond
	}
	if IsTruthy(Cond) {
		return Eval(I.Cons, Env)
	}
	for _, Elif := range I.Elifs {
		ECond := Eval(Elif.Cond, Env)
		if IsErr(ECond) {
			return ECond
		}
		if IsTruthy(ECond) {
			return Eval(Elif.Cons, Env)
		}
	}
	if I.Alt != nil {
		return Eval(I.Alt, Env)
	}
	return NullObj
}

func IsTruthy(Obj Object.Object) bool {
	switch Obj {
	case NullObj:
		return false
	case TrueObj:
		return true
	case FalseObj:
		return false
	default:
		return true
	}
}

func EvalIdent(I *Ast.Ident, Env *Object.Environment) Object.Object {
	if Val, Ok := Env.Get(I.Value); Ok {
		return Val
	}
	return NewErr(I.Token.Line, "Ident not found: "+I.Value)
}

func EvalExprs(Exprs []Ast.Expression, Env *Object.Environment) []Object.Object {
	var Res []Object.Object
	for _, E := range Exprs {
		Ev := Eval(E, Env)
		if IsErr(Ev) {
			return []Object.Object{Ev}
		}
		Res = append(Res, Ev)
	}
	return Res
}

func Apply(Fn Object.Object, Args []Object.Object, Line int) Object.Object {
	switch F := Fn.(type) {
	case *Object.Func:
		ExtEnv := ExtFuncEnv(F, Args)
		Ev := Eval(F.Body, ExtEnv)
		return UnwrapRet(Ev)
	case *Object.BoundMethod:
		ExtEnv := ExtFuncEnv(F.Method, Args)
		ExtEnv.Set("$this", F.Self)
		Ev := Eval(F.Method.Body, ExtEnv)
		return UnwrapRet(Ev)
	case *Object.Builtin:
		if len(Args) > 0 {
			for _, A := range Args {
				if A != nil && A.Type() == Object.ErrObj {
					if ErrObj, Ok := A.(*Object.Err); Ok {
						ErrObj.Msg = fmt.Sprintf("Line %d: ", Line) + ErrObj.Msg
					}
					return A
				}
			}
		}
		Res := F.Fn(Args...)
		if Res != nil && Res.Type() == Object.ErrObj {
			if ErrObj, Ok := Res.(*Object.Err); Ok {
				ErrObj.Msg = fmt.Sprintf("Line %d: ", Line) + ErrObj.Msg
			}
		}
		return Res
	default:
		return NewErr(Line, "Not a func: %s", Fn.Type())
	}
}

func ExtFuncEnv(F *Object.Func, Args []Object.Object) *Object.Environment {
	Env := Object.NewEncEnv(F.Env)
	for Idx, P := range F.Params {
		if Idx < len(Args) {
			Env.Set(P.Value, Args[Idx])
		}
	}
	return Env
}

func UnwrapRet(Obj Object.Object) Object.Object {
	if R, Ok := Obj.(*Object.Ret); Ok {
		return R.Value
	}
	return Obj
}

func IsErr(Obj Object.Object) bool {
	if Obj != nil {
		return Obj.Type() == Object.ErrObj
	}
	return false
}

func NewErr(Line int, Format string, Args ...interface{}) *Object.Err {
	return &Object.Err{Msg: fmt.Sprintf("Line %d: ", Line) + fmt.Sprintf(Format, Args...)}
}