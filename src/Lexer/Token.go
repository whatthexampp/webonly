package Lexer

import "strings"

type TokenType string

type Token struct {
	Type   TokenType
	Lit    string
	Line   int
	Column int
}

const (
	Illegal  TokenType = "Illegal"
	Eof      TokenType = "Eof"
	Html     TokenType = "Html"
	Ident    TokenType = "Ident"
	Num      TokenType = "Num"
	Str      TokenType = "Str"
	Assign   TokenType = "="
	Plus     TokenType = "+"
	Minus    TokenType = "-"
	Ast      TokenType = "*"
	Slash    TokenType = "/"
	Bang     TokenType = "!"
	Lt       TokenType = "<"
	Gt       TokenType = ">"
	Eq       TokenType = "=="
	Neq      TokenType = "!="
	Comma    TokenType = ","
	Semi     TokenType = ";"
	Lparen   TokenType = "("
	Rparen   TokenType = ")"
	Lbracket TokenType = "["
	Rbracket TokenType = "]"
	Colon    TokenType = ":"
	Dot      TokenType = "."
	Func     TokenType = "Func"
	True     TokenType = "True"
	False    TokenType = "False"
	If       TokenType = "If"
	Elseif   TokenType = "Elseif"
	Else     TokenType = "Else"
	While    TokenType = "While"
	Ret      TokenType = "Ret"
	End      TokenType = "End"
	Class    TokenType = "Class"
	New      TokenType = "New"
)

var Keywords = map[string]TokenType{
	"function": Func,
	"true":     True,
	"false":    False,
	"if":       If,
	"elseif":   Elseif,
	"else":     Else,
	"while":    While,
	"return":   Ret,
	"end":      End,
	"class":    Class,
	"new":      New,
}

func Lookup(Identifier string) TokenType {
	if Tok, Ok := Keywords[strings.ToLower(Identifier)]; Ok {
		return Tok
	}
	return Ident
}