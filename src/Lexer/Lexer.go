package Lexer

import "strings"

type Lexer struct {
	Input  string
	Pos    int
	Read   int
	Char   byte
	Line   int
	Column int
	InCode bool
}

func Create(Input string) *Lexer {
	Lex := &Lexer{Input: Input, Line: 1, Column: 0, InCode: false}
	Lex.Advance()
	return Lex
}

func (Lex *Lexer) Advance() {
	if Lex.Read >= len(Lex.Input) {
		Lex.Char = 0
	} else {
		Lex.Char = Lex.Input[Lex.Read]
	}
	Lex.Pos = Lex.Read
	Lex.Read++
	Lex.Column++
}

func (Lex *Lexer) Peek() byte {
	if Lex.Read >= len(Lex.Input) {
		return 0
	}
	return Lex.Input[Lex.Read]
}

func (Lex *Lexer) Next() Token {
	if !Lex.InCode {
		Start := Lex.Pos
		StartLine := Lex.Line
		StartCol := Lex.Column
		for Lex.Char != 0 {
			if Lex.Char == '<' && strings.HasPrefix(Lex.Input[Lex.Pos:], "<?webonly") {
				if Lex.Pos > Start {
					Text := Lex.Input[Start:Lex.Pos]
					return Token{Type: Html, Lit: Text, Line: StartLine, Column: StartCol}
				}
				for I := 0; I < 9; I++ {
					Lex.Advance()
				}
				Lex.InCode = true
				return Lex.Next()
			}
			if Lex.Char == '\n' {
				Lex.Line++
				Lex.Column = 0
			}
			Lex.Advance()
		}
		if Lex.Pos > Start {
			return Token{Type: Html, Lit: Lex.Input[Start:Lex.Pos], Line: StartLine, Column: StartCol}
		}
		return Token{Type: Eof, Lit: ""}
	}

	Lex.SkipSpace()

	if Lex.Char == '?' && Lex.Peek() == '>' {
		Lex.Advance()
		Lex.Advance()
		Lex.InCode = false
		return Lex.Next()
	}

	Tok := Token{Line: Lex.Line, Column: Lex.Column}
	switch Lex.Char {
	case '=':
		if Lex.Peek() == '=' {
			Prev := Lex.Char
			Lex.Advance()
			Tok = Token{Type: Eq, Lit: string(Prev) + string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
		} else {
			Tok = Token{Type: Assign, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
		}
	case '+':
		Tok = Token{Type: Plus, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case '-':
		Tok = Token{Type: Minus, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case '!':
		if Lex.Peek() == '=' {
			Prev := Lex.Char
			Lex.Advance()
			Tok = Token{Type: Neq, Lit: string(Prev) + string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
		} else {
			Tok = Token{Type: Bang, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
		}
	case '/':
		Tok = Token{Type: Slash, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case '*':
		Tok = Token{Type: Ast, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case '<':
		Tok = Token{Type: Lt, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case '>':
		Tok = Token{Type: Gt, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case ';':
		Tok = Token{Type: Semi, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case ',':
		Tok = Token{Type: Comma, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case '(':
		Tok = Token{Type: Lparen, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case ')':
		Tok = Token{Type: Rparen, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case '[':
		Tok = Token{Type: Lbracket, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case ']':
		Tok = Token{Type: Rbracket, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case ':':
		Tok = Token{Type: Colon, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case '.':
		Tok = Token{Type: Dot, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	case '"':
		Tok.Type = Str
		Tok.Lit = Lex.ReadStr()
	case 0:
		Tok.Type = Eof
		Tok.Lit = ""
	default:
		if Lex.IsAlpha(Lex.Char) {
			Tok.Lit = Lex.ReadIdent()
			Tok.Type = Lookup(Tok.Lit)
			return Tok
		} else if Lex.IsNum(Lex.Char) {
			Tok.Type = Num
			Tok.Lit = Lex.ReadNum()
			return Tok
		}
		Tok = Token{Type: Illegal, Lit: string(Lex.Char), Line: Tok.Line, Column: Tok.Column}
	}
	Lex.Advance()
	return Tok
}

func (Lex *Lexer) ReadStr() string {
	Start := Lex.Pos + 1
	for {
		Lex.Advance()
		if Lex.Char == '"' || Lex.Char == 0 {
			break
		}
	}
	return Lex.Input[Start:Lex.Pos]
}

func (Lex *Lexer) ReadIdent() string {
	Start := Lex.Pos
	for Lex.IsAlpha(Lex.Char) || Lex.IsNum(Lex.Char) {
		Lex.Advance()
	}
	return Lex.Input[Start:Lex.Pos]
}

func (Lex *Lexer) ReadNum() string {
	Start := Lex.Pos
	for Lex.IsNum(Lex.Char) || Lex.Char == '.' {
		Lex.Advance()
	}
	return Lex.Input[Start:Lex.Pos]
}

func (Lex *Lexer) IsAlpha(Ch byte) bool {
	return (Ch >= 'a' && Ch <= 'z') || (Ch >= 'A' && Ch <= 'Z') || Ch == '_' || Ch == '$'
}

func (Lex *Lexer) IsNum(Ch byte) bool {
	return Ch >= '0' && Ch <= '9'
}

func (Lex *Lexer) SkipSpace() {
	for Lex.Char == ' ' || Lex.Char == '\t' || Lex.Char == '\n' || Lex.Char == '\r' {
		if Lex.Char == '\n' {
			Lex.Line++
			Lex.Column = 0
		}
		Lex.Advance()
	}
}