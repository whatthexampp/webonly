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
	Modulo   TokenType = "%"
	Bang     TokenType = "!"
	Lt       TokenType = "<"
	Gt       TokenType = ">"
	Eq       TokenType = "=="
	Neq      TokenType = "!="
	And      TokenType = "&&"
	Or       TokenType = "||"
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
	Null     TokenType = "Null"
	If       TokenType = "If"
	Elseif   TokenType = "Elseif"
	Else     TokenType = "Else"
	While    TokenType = "While"
	Ret      TokenType = "Ret"
	End      TokenType = "End"
	Endif    TokenType = "Endif"
	Class    TokenType = "Class"
	Extends  TokenType = "Extends"
	New      TokenType = "New"
	Const    TokenType = "Const"
	Enum     TokenType = "Enum"
	Import   TokenType = "Import"
	Public   TokenType = "Public"
)

var Keywords = map[string]TokenType{
	"function": Func,
	"true":     True,
	"false":    False,
	"null":     Null,
	"if":       If,
	"elseif":   Elseif,
	"else":     Else,
	"while":    While,
	"return":   Ret,
	"end":      End,
	"endif":    Endif,
	"class":    Class,
	"extends":  Extends,
	"new":      New,
	"const":    Const,
	"enum":     Enum,
	"import":   Import,
	"public":   Public,
}

func Lookup(Identifier string) TokenType {
	if Tok, Ok := Keywords[strings.ToLower(Identifier)]; Ok {
		return Tok
	}
	return Ident
}