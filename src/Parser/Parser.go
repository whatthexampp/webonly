package Parser

import (
	"fmt"
	"strconv"
	"strings"
	"webonly/src/Ast"
	"webonly/src/Lexer"
)

const (
	Lowest = iota
	AssignPrec
	LogicOr
	LogicAnd
	Equals
	LessGt
	Sum
	Prod
	Prefix
	Call
	IndexPrec
	DotPrec
)

var Precs = map[Lexer.TokenType]int{
	Lexer.Assign:   AssignPrec,
	Lexer.Or:       LogicOr,
	Lexer.And:      LogicAnd,
	Lexer.Eq:       Equals,
	Lexer.Neq:      Equals,
	Lexer.Lt:       LessGt,
	Lexer.Gt:       LessGt,
	Lexer.Plus:     Sum,
	Lexer.Minus:    Sum,
	Lexer.Slash:    Prod,
	Lexer.Ast:      Prod,
	Lexer.Modulo:   Prod,
	Lexer.Lparen:   Call,
	Lexer.Lbracket: IndexPrec,
	Lexer.Dot:      DotPrec,
}

type PrefixFn func() Ast.Expression
type InfixFn func(Ast.Expression) Ast.Expression

type Parser struct {
	Lex    *Lexer.Lexer
	Cur    Lexer.Token
	Peek   Lexer.Token
	Errs   []string
	PreFns map[Lexer.TokenType]PrefixFn
	InFns  map[Lexer.TokenType]InfixFn
}

func Create(L *Lexer.Lexer) *Parser {
	P := &Parser{Lex: L, Errs: []string{}}
	P.PreFns = make(map[Lexer.TokenType]PrefixFn)
	P.RegPre(Lexer.Ident, P.ParseIdent)
	P.RegPre(Lexer.Num, P.ParseNum)
	P.RegPre(Lexer.Str, P.ParseStr)
	P.RegPre(Lexer.Null, P.ParseNull)
	P.RegPre(Lexer.Bang, P.ParsePre)
	P.RegPre(Lexer.Minus, P.ParsePre)
	P.RegPre(Lexer.True, P.ParseBool)
	P.RegPre(Lexer.False, P.ParseBool)
	P.RegPre(Lexer.Lparen, P.ParseGroup)
	P.RegPre(Lexer.Lbracket, P.ParseArrayLit)
	P.RegPre(Lexer.If, P.ParseIf)
	P.RegPre(Lexer.Func, P.ParseFunc)
	P.RegPre(Lexer.New, P.ParseNew)
	P.RegPre(Lexer.Import, P.ParseImport)
	P.InFns = make(map[Lexer.TokenType]InfixFn)
	P.RegIn(Lexer.Plus, P.ParseInfix)
	P.RegIn(Lexer.Minus, P.ParseInfix)
	P.RegIn(Lexer.Slash, P.ParseInfix)
	P.RegIn(Lexer.Ast, P.ParseInfix)
	P.RegIn(Lexer.Modulo, P.ParseInfix)
	P.RegIn(Lexer.Eq, P.ParseInfix)
	P.RegIn(Lexer.Neq, P.ParseInfix)
	P.RegIn(Lexer.Lt, P.ParseInfix)
	P.RegIn(Lexer.Gt, P.ParseInfix)
	P.RegIn(Lexer.And, P.ParseInfix)
	P.RegIn(Lexer.Or, P.ParseInfix)
	P.RegIn(Lexer.Lparen, P.ParseCall)
	P.RegIn(Lexer.Lbracket, P.ParseIndex)
	P.RegIn(Lexer.Dot, P.ParseDot)
	P.RegIn(Lexer.Assign, P.ParseAssignExpr)
	P.Next()
	P.Next()
	return P
}

func (P *Parser) Next() {
	P.Cur = P.Peek
	P.Peek = P.Lex.Next()
}

func (P *Parser) ParseProg() *Ast.Program {
	Prog := &Ast.Program{Stmts: []Ast.Statement{}}
	for P.Cur.Type != Lexer.Eof {
		if S := P.ParseStmt(); S != nil {
			Prog.Stmts = append(Prog.Stmts, S)
		}
		P.Next()
	}
	return Prog
}

func (P *Parser) ParseStmt() Ast.Statement {
	switch P.Cur.Type {
	case Lexer.Html:
		return P.ParseHtml()
	case Lexer.Class:
		return P.ParseClass()
	case Lexer.While:
		return P.ParseWhile()
	case Lexer.Ret:
		return P.ParseRet()
	case Lexer.Const:
		return P.ParseConst()
	case Lexer.Enum:
		return P.ParseEnum()
	case Lexer.Public:
		return P.ParsePublic()
	default:
		return P.ParseExprStmt()
	}
}

func (P *Parser) ParseHtml() *Ast.HtmlStmt {
	Stmt := &Ast.HtmlStmt{Token: P.Cur, Value: P.Cur.Lit}
	return Stmt
}

func (P *Parser) ParseClass() *Ast.ClassStmt {
	Stmt := &Ast.ClassStmt{Token: P.Cur}
	if !P.Exp(Lexer.Ident) {
		return nil
	}
	Stmt.Name = &Ast.Ident{Token: P.Cur, Value: P.Cur.Lit}

	if P.Peek.Type == Lexer.Extends {
		P.Next()
		if !P.Exp(Lexer.Ident) {
			return nil
		}
		Stmt.Parent = &Ast.Ident{Token: P.Cur, Value: P.Cur.Lit}
	}

	if !P.Exp(Lexer.Colon) {
		return nil
	}
	P.Next()
	Stmt.Methods = []*Ast.FuncLit{}
	for !P.IsAtEnd(Lexer.End) && P.Cur.Type != Lexer.Eof {
		if P.Cur.Type == Lexer.Func {
			M := P.ParseFuncDef()
			if M != nil {
				Stmt.Methods = append(Stmt.Methods, M)
			}
		} else {
			P.Next()
		}
	}
	if P.Cur.Type != Lexer.End {
		P.Errs = append(P.Errs, fmt.Sprintf("Expected end, got %s at line %d", P.Cur.Type, P.Cur.Line))
	}
	if P.Peek.Type == Lexer.Semi {
		P.Next()
	}
	return Stmt
}

func (P *Parser) ParseConst() *Ast.ConstStmt {
	Stmt := &Ast.ConstStmt{Token: P.Cur}
	if !P.Exp(Lexer.Ident) {
		return nil
	}
	if !strings.HasPrefix(P.Cur.Lit, "$") {
		P.Errs = append(P.Errs, fmt.Sprintf("Constant %s must start with '$' at line %d", P.Cur.Lit, P.Cur.Line))
	}
	Stmt.Name = &Ast.Ident{Token: P.Cur, Value: P.Cur.Lit}
	if !P.Exp(Lexer.Assign) {
		return nil
	}
	P.Next()
	Stmt.Value = P.ParseExpr(Lowest)
	if P.Peek.Type == Lexer.Semi {
		P.Next()
	}
	return Stmt
}

func (P *Parser) ParseEnum() *Ast.EnumStmt {
	Stmt := &Ast.EnumStmt{Token: P.Cur}
	if !P.Exp(Lexer.Ident) {
		return nil
	}
	Stmt.Name = &Ast.Ident{Token: P.Cur, Value: P.Cur.Lit}
	if !P.Exp(Lexer.Colon) {
		return nil
	}
	P.Next()
	Stmt.Cases = []*Ast.Ident{}
	for !P.IsAtEnd(Lexer.End) && P.Cur.Type != Lexer.Eof {
		if P.Cur.Type == Lexer.Ident {
			Stmt.Cases = append(Stmt.Cases, &Ast.Ident{Token: P.Cur, Value: P.Cur.Lit})
		} else if P.Cur.Type != Lexer.Comma {
			P.Errs = append(P.Errs, fmt.Sprintf("Unexpected token in enum: %s at line %d", P.Cur.Type, P.Cur.Line))
		}
		P.Next()
	}
	if P.Cur.Type != Lexer.End {
		P.Errs = append(P.Errs, fmt.Sprintf("Expected end, got %s at line %d", P.Cur.Type, P.Cur.Line))
	}
	if P.Peek.Type == Lexer.Semi {
		P.Next()
	}
	return Stmt
}

func (P *Parser) ParsePublic() *Ast.PublicStmt {
	PubToken := P.Cur
	P.Next()
	Stmt := P.ParseStmt()
	if Stmt == nil {
		return nil
	}
	return &Ast.PublicStmt{Token: PubToken, Stmt: Stmt}
}

func (P *Parser) ParseWhile() *Ast.WhileStmt {
	Stmt := &Ast.WhileStmt{Token: P.Cur}
	if !P.Exp(Lexer.Lparen) {
		return nil
	}
	P.Next()
	Stmt.Cond = P.ParseExpr(Lowest)
	if !P.Exp(Lexer.Rparen) || !P.Exp(Lexer.Colon) {
		return nil
	}
	P.Next()
	Stmt.Body = P.ParseBlock(Lexer.End)
	if P.Cur.Type != Lexer.End {
		P.Errs = append(P.Errs, fmt.Sprintf("Expected end, got %s at line %d", P.Cur.Type, P.Cur.Line))
	}
	P.Next()
	if P.Peek.Type == Lexer.Semi {
		P.Next()
	}
	return Stmt
}

func (P *Parser) ParseRet() *Ast.RetStmt {
	Stmt := &Ast.RetStmt{Token: P.Cur}
	P.Next()
	Stmt.Value = P.ParseExpr(Lowest)
	if P.Peek.Type == Lexer.Semi {
		P.Next()
	}
	return Stmt
}

func (P *Parser) ParseExprStmt() *Ast.ExprStmt {
	Stmt := &Ast.ExprStmt{Token: P.Cur, Expr: P.ParseExpr(Lowest)}
	if P.Peek.Type == Lexer.Semi {
		P.Next()
	}
	return Stmt
}

func (P *Parser) ParseExpr(Prec int) Ast.Expression {
	PreFn := P.PreFns[P.Cur.Type]
	if PreFn == nil {
		P.Errs = append(P.Errs, fmt.Sprintf("No prefix function for %s at line %d", P.Cur.Type, P.Cur.Line))
		return nil
	}
	Left := PreFn()
	for P.Peek.Type != Lexer.Semi && Prec < P.PeekPrec() {
		InFn := P.InFns[P.Peek.Type]
		if InFn == nil {
			return Left
		}
		P.Next()
		Left = InFn(Left)
	}
	return Left
}

func (P *Parser) ParseAssignExpr(Left Ast.Expression) Ast.Expression {
	E := &Ast.AssignExpr{Token: P.Cur, Left: Left}
	Prec := P.CurPrec()
	P.Next()
	E.Value = P.ParseExpr(Prec - 1)
	return E
}

func (P *Parser) ParseIdent() Ast.Expression {
	return &Ast.Ident{Token: P.Cur, Value: P.Cur.Lit}
}

func (P *Parser) ParseNum() Ast.Expression {
	Val, Err := strconv.ParseFloat(P.Cur.Lit, 64)
	if Err != nil {
		P.Errs = append(P.Errs, fmt.Sprintf("Parse fail %q as number at line %d", P.Cur.Lit, P.Cur.Line))
		return nil
	}
	return &Ast.NumLit{Token: P.Cur, Value: Val}
}

func (P *Parser) ParseStr() Ast.Expression {
	return &Ast.StrLit{Token: P.Cur, Value: P.Cur.Lit}
}

func (P *Parser) ParseNull() Ast.Expression {
	return &Ast.NullLit{Token: P.Cur}
}

func (P *Parser) ParseArrayLit() Ast.Expression {
	A := &Ast.ArrayLit{Token: P.Cur}
	A.Elems = P.ParseExprList(Lexer.Rbracket)
	return A
}

func (P *Parser) ParsePre() Ast.Expression {
	E := &Ast.PrefixExpr{Token: P.Cur, Op: P.Cur.Lit}
	P.Next()
	E.Right = P.ParseExpr(Prefix)
	return E
}

func (P *Parser) ParseInfix(Left Ast.Expression) Ast.Expression {
	E := &Ast.InfixExpr{Token: P.Cur, Op: P.Cur.Lit, Left: Left}
	Prec := P.CurPrec()
	P.Next()
	E.Right = P.ParseExpr(Prec)
	return E
}

func (P *Parser) ParseBool() Ast.Expression {
	return &Ast.Bool{Token: P.Cur, Value: P.Cur.Type == Lexer.True}
}

func (P *Parser) ParseGroup() Ast.Expression {
	P.Next()
	E := P.ParseExpr(Lowest)
	if !P.Exp(Lexer.Rparen) {
		return nil
	}
	return E
}

func (P *Parser) ParseIf() Ast.Expression {
	E := &Ast.IfExpr{Token: P.Cur, Elifs: []*Ast.ElseifBlock{}}
	HasParen := false
	if P.Peek.Type == Lexer.Lparen {
		P.Next()
		HasParen = true
	}
	P.Next()
	E.Cond = P.ParseExpr(Lowest)
	if HasParen {
		if !P.Exp(Lexer.Rparen) {
			return nil
		}
	}
	if !P.Exp(Lexer.Colon) {
		return nil
	}
	P.Next()
	E.Cons = P.ParseBlock(Lexer.Elseif, Lexer.Else, Lexer.End, Lexer.Endif)
	for P.Cur.Type == Lexer.Elseif {
		Elif := &Ast.ElseifBlock{}
		HasParenElif := false
		if P.Peek.Type == Lexer.Lparen {
			P.Next()
			HasParenElif = true
		}
		P.Next()
		Elif.Cond = P.ParseExpr(Lowest)
		if HasParenElif {
			if !P.Exp(Lexer.Rparen) {
				return nil
			}
		}
		if !P.Exp(Lexer.Colon) {
			return nil
		}
		P.Next()
		Elif.Cons = P.ParseBlock(Lexer.Elseif, Lexer.Else, Lexer.End, Lexer.Endif)
		E.Elifs = append(E.Elifs, Elif)
	}
	if P.Cur.Type == Lexer.Else {
		if !P.Exp(Lexer.Colon) {
			return nil
		}
		P.Next()
		E.Alt = P.ParseBlock(Lexer.End, Lexer.Endif)
	}
	if P.Cur.Type != Lexer.End && P.Cur.Type != Lexer.Endif {
		P.Errs = append(P.Errs, fmt.Sprintf("Expected end or endif, got %s at line %d", P.Cur.Type, P.Cur.Line))
	}
	return E
}

func (P *Parser) ParseBlock(EndTokens ...Lexer.TokenType) *Ast.BlockStmt {
	B := &Ast.BlockStmt{Token: P.Cur, Stmts: []Ast.Statement{}}
	for !P.IsAtEnd(EndTokens...) && P.Cur.Type != Lexer.Eof {
		if S := P.ParseStmt(); S != nil {
			B.Stmts = append(B.Stmts, S)
		}
		P.Next()
	}
	return B
}

func (P *Parser) IsAtEnd(EndTokens ...Lexer.TokenType) bool {
	for _, T := range EndTokens {
		if P.Cur.Type == T {
			return true
		}
	}
	return false
}

func (P *Parser) ParseFuncDef() *Ast.FuncLit {
	F := &Ast.FuncLit{Token: P.Cur}
	if !P.Exp(Lexer.Ident) {
		return nil
	}
	F.Name = P.Cur.Lit
	if !P.Exp(Lexer.Lparen) {
		return nil
	}
	F.Params = P.ParseParams()
	if !P.Exp(Lexer.Colon) {
		return nil
	}
	P.Next()
	F.Body = P.ParseBlock(Lexer.End)
	if P.Cur.Type != Lexer.End {
		P.Errs = append(P.Errs, fmt.Sprintf("Expected end, got %s at line %d", P.Cur.Type, P.Cur.Line))
	}
	P.Next()
	return F
}

func (P *Parser) ParseFunc() Ast.Expression {
	F := &Ast.FuncLit{Token: P.Cur}
	if P.Peek.Type == Lexer.Ident {
		P.Next()
		F.Name = P.Cur.Lit
	}
	if !P.Exp(Lexer.Lparen) {
		return nil
	}
	F.Params = P.ParseParams()
	if !P.Exp(Lexer.Colon) {
		return nil
	}
	P.Next()
	F.Body = P.ParseBlock(Lexer.End)
	if P.Cur.Type != Lexer.End {
		P.Errs = append(P.Errs, fmt.Sprintf("Expected end, got %s at line %d", P.Cur.Type, P.Cur.Line))
	}
	return F
}

func (P *Parser) ParseParams() []*Ast.Ident {
	Pms := []*Ast.Ident{}
	if P.Peek.Type == Lexer.Rparen {
		P.Next()
		return Pms
	}
	P.Next()
	Pms = append(Pms, &Ast.Ident{Token: P.Cur, Value: P.Cur.Lit})
	for P.Peek.Type == Lexer.Comma {
		P.Next()
		P.Next()
		Pms = append(Pms, &Ast.Ident{Token: P.Cur, Value: P.Cur.Lit})
	}
	if !P.Exp(Lexer.Rparen) {
		return nil
	}
	return Pms
}

func (P *Parser) ParseCall(F Ast.Expression) Ast.Expression {
	E := &Ast.CallExpr{Token: P.Cur, Func: F}
	E.Args = P.ParseExprList(Lexer.Rparen)
	return E
}

func (P *Parser) ParseIndex(Left Ast.Expression) Ast.Expression {
	E := &Ast.IndexExpr{Token: P.Cur, Left: Left}
	P.Next()
	E.Index = P.ParseExpr(Lowest)
	if !P.Exp(Lexer.Rbracket) {
		return nil
	}
	return E
}

func (P *Parser) ParseDot(Left Ast.Expression) Ast.Expression {
	E := &Ast.DotExpr{Token: P.Cur, Left: Left}
	if !P.Exp(Lexer.Ident) {
		return nil
	}
	E.Right = &Ast.Ident{Token: P.Cur, Value: P.Cur.Lit}
	return E
}

func (P *Parser) ParseNew() Ast.Expression {
	E := &Ast.NewExpr{Token: P.Cur}
	P.Next()
	E.Class = P.ParseExpr(Call)
	if !P.Exp(Lexer.Lparen) {
		return nil
	}
	E.Args = P.ParseExprList(Lexer.Rparen)
	return E
}

func (P *Parser) ParseImport() Ast.Expression {
	E := &Ast.ImportExpr{Token: P.Cur}
	P.Next()
	E.Path = P.ParseExpr(Prefix)
	return E
}

func (P *Parser) ParseExprList(End Lexer.TokenType) []Ast.Expression {
	List := []Ast.Expression{}
	if P.Peek.Type == End {
		P.Next()
		return List
	}
	P.Next()
	List = append(List, P.ParseExpr(Lowest))
	for P.Peek.Type == Lexer.Comma {
		P.Next()
		P.Next()
		List = append(List, P.ParseExpr(Lowest))
	}
	if !P.Exp(End) {
		return nil
	}
	return List
}

func (P *Parser) Exp(T Lexer.TokenType) bool {
	if P.Peek.Type == T {
		P.Next()
		return true
	}
	P.Errs = append(P.Errs, fmt.Sprintf("Expected next token %s, got %s at line %d", T, P.Peek.Type, P.Peek.Line))
	return false
}

func (P *Parser) PeekPrec() int {
	if Prc, Ok := Precs[P.Peek.Type]; Ok {
		return Prc
	}
	return Lowest
}

func (P *Parser) CurPrec() int {
	if Prc, Ok := Precs[P.Cur.Type]; Ok {
		return Prc
	}
	return Lowest
}

func (P *Parser) RegPre(T Lexer.TokenType, F PrefixFn) {
	P.PreFns[T] = F
}

func (P *Parser) RegIn(T Lexer.TokenType, F InfixFn) {
	P.InFns[T] = F
}