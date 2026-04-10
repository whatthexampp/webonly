package Ast

import (
	"fmt"
	"strings"
	"webonly/src/Lexer"
)

type Node interface {
	TokenLit() string
	String() string
}

type Statement interface {
	Node
	StmtNode()
}

type Expression interface {
	Node
	ExprNode()
}

type Program struct {
	Stmts []Statement
}

func (P *Program) TokenLit() string {
	if len(P.Stmts) > 0 {
		return P.Stmts[0].TokenLit()
	}
	return ""
}

func (P *Program) String() string {
	Out := ""
	for _, S := range P.Stmts {
		Out += S.String()
	}
	return Out
}

type HtmlStmt struct {
	Token Lexer.Token
	Value string
}

func (H *HtmlStmt) StmtNode() {
}

func (H *HtmlStmt) TokenLit() string {
	return H.Token.Lit
}

func (H *HtmlStmt) String() string {
	return H.Value
}

type ClassStmt struct {
	Token   Lexer.Token
	Name    *Ident
	Methods []*FuncLit
}

func (C *ClassStmt) StmtNode() {
}

func (C *ClassStmt) TokenLit() string {
	return C.Token.Lit
}

func (C *ClassStmt) String() string {
	MStr := ""
	for _, M := range C.Methods {
		MStr += M.String()
	}
	return "class " + C.Name.String() + ":\n" + MStr + "\nend;"
}

type RetStmt struct {
	Token Lexer.Token
	Value Expression
}

func (R *RetStmt) StmtNode() {
}

func (R *RetStmt) TokenLit() string {
	return R.Token.Lit
}

func (R *RetStmt) String() string {
	if R.Value != nil {
		return R.TokenLit() + " " + R.Value.String() + ";"
	}
	return R.TokenLit() + ";"
}

type ExprStmt struct {
	Token Lexer.Token
	Expr  Expression
}

func (E *ExprStmt) StmtNode() {
}

func (E *ExprStmt) TokenLit() string {
	return E.Token.Lit
}

func (E *ExprStmt) String() string {
	if E.Expr != nil {
		return E.Expr.String()
	}
	return ""
}

type BlockStmt struct {
	Token Lexer.Token
	Stmts []Statement
}

func (B *BlockStmt) StmtNode() {
}

func (B *BlockStmt) TokenLit() string {
	return B.Token.Lit
}

func (B *BlockStmt) String() string {
	Out := ""
	for _, S := range B.Stmts {
		Out += S.String()
	}
	return Out
}

type WhileStmt struct {
	Token Lexer.Token
	Cond  Expression
	Body  *BlockStmt
}

func (W *WhileStmt) StmtNode() {
}

func (W *WhileStmt) TokenLit() string {
	return W.Token.Lit
}

func (W *WhileStmt) String() string {
	return "while [" + W.Cond.String() + "]:\n" + W.Body.String() + "\nend;"
}

type Ident struct {
	Token Lexer.Token
	Value string
}

func (I *Ident) ExprNode() {
}

func (I *Ident) TokenLit() string {
	return I.Token.Lit
}

func (I *Ident) String() string {
	return I.Value
}

type NumLit struct {
	Token Lexer.Token
	Value float64
}

func (N *NumLit) ExprNode() {
}

func (N *NumLit) TokenLit() string {
	return N.Token.Lit
}

func (N *NumLit) String() string {
	return fmt.Sprintf("%g", N.Value)
}

type StrLit struct {
	Token Lexer.Token
	Value string
}

func (S *StrLit) ExprNode() {
}

func (S *StrLit) TokenLit() string {
	return S.Token.Lit
}

func (S *StrLit) String() string {
	return S.Token.Lit
}

type ArrayLit struct {
	Token Lexer.Token
	Elems []Expression
}

func (A *ArrayLit) ExprNode() {
}

func (A *ArrayLit) TokenLit() string {
	return A.Token.Lit
}

func (A *ArrayLit) String() string {
	Str := []string{}
	for _, E := range A.Elems {
		Str = append(Str, E.String())
	}
	return "[" + strings.Join(Str, ", ") + "]"
}

type PrefixExpr struct {
	Token Lexer.Token
	Op    string
	Right Expression
}

func (P *PrefixExpr) ExprNode() {
}

func (P *PrefixExpr) TokenLit() string {
	return P.Token.Lit
}

func (P *PrefixExpr) String() string {
	return "(" + P.Op + P.Right.String() + ")"
}

type InfixExpr struct {
	Token Lexer.Token
	Left  Expression
	Op    string
	Right Expression
}

func (I *InfixExpr) ExprNode() {
}

func (I *InfixExpr) TokenLit() string {
	return I.Token.Lit
}

func (I *InfixExpr) String() string {
	return "(" + I.Left.String() + " " + I.Op + " " + I.Right.String() + ")"
}

type AssignExpr struct {
	Token Lexer.Token
	Left  Expression
	Value Expression
}

func (A *AssignExpr) ExprNode() {
}

func (A *AssignExpr) TokenLit() string {
	return A.Token.Lit
}

func (A *AssignExpr) String() string {
	return "(" + A.Left.String() + " = " + A.Value.String() + ")"
}

type IndexExpr struct {
	Token Lexer.Token
	Left  Expression
	Index Expression
}

func (I *IndexExpr) ExprNode() {
}

func (I *IndexExpr) TokenLit() string {
	return I.Token.Lit
}

func (I *IndexExpr) String() string {
	return I.Left.String() + "[" + I.Index.String() + "]"
}

type DotExpr struct {
	Token Lexer.Token
	Left  Expression
	Right *Ident
}

func (D *DotExpr) ExprNode() {
}

func (D *DotExpr) TokenLit() string {
	return D.Token.Lit
}

func (D *DotExpr) String() string {
	return D.Left.String() + "." + D.Right.String()
}

type Bool struct {
	Token Lexer.Token
	Value bool
}

func (B *Bool) ExprNode() {
}

func (B *Bool) TokenLit() string {
	return B.Token.Lit
}

func (B *Bool) String() string {
	return B.Token.Lit
}

type ElseifBlock struct {
	Cond Expression
	Cons *BlockStmt
}

type IfExpr struct {
	Token Lexer.Token
	Cond  Expression
	Cons  *BlockStmt
	Elifs []*ElseifBlock
	Alt   *BlockStmt
}

func (I *IfExpr) ExprNode() {
}

func (I *IfExpr) TokenLit() string {
	return I.Token.Lit
}

func (I *IfExpr) String() string {
	Out := "if [" + I.Cond.String() + "]:\n" + I.Cons.String()
	for _, E := range I.Elifs {
		Out += "elseif [" + E.Cond.String() + "]:\n" + E.Cons.String()
	}
	if I.Alt != nil {
		Out += "else:\n" + I.Alt.String()
	}
	Out += "end;"
	return Out
}

type FuncLit struct {
	Token  Lexer.Token
	Name   string
	Params []*Ident
	Body   *BlockStmt
}

func (F *FuncLit) ExprNode() {
}

func (F *FuncLit) TokenLit() string {
	return F.Token.Lit
}

func (F *FuncLit) String() string {
	P := []string{}
	for _, Param := range F.Params {
		P = append(P, Param.String())
	}
	NStr := ""
	if F.Name != "" {
		NStr = " " + F.Name
	}
	return "function" + NStr + "(" + strings.Join(P, ", ") + "):\n" + F.Body.String() + "\nend;"
}

type CallExpr struct {
	Token Lexer.Token
	Func  Expression
	Args  []Expression
}

func (C *CallExpr) ExprNode() {
}

func (C *CallExpr) TokenLit() string {
	return C.Token.Lit
}

func (C *CallExpr) String() string {
	A := []string{}
	for _, Arg := range C.Args {
		A = append(A, Arg.String())
	}
	return C.Func.String() + "(" + strings.Join(A, ", ") + ")"
}

type NewExpr struct {
	Token Lexer.Token
	Class *Ident
	Args  []Expression
}

func (N *NewExpr) ExprNode() {
}

func (N *NewExpr) TokenLit() string {
	return N.Token.Lit
}

func (N *NewExpr) String() string {
	A := []string{}
	for _, Arg := range N.Args {
		A = append(A, Arg.String())
	}
	return "new " + N.Class.String() + "(" + strings.Join(A, ", ") + ")"
}