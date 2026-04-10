package Funcs

import (
	"strings"
	"webonly/src/Object"
)

var StrReplace = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		if len(Args) != 3 || Args[0].Type() != Object.StrObj || Args[1].Type() != Object.StrObj || Args[2].Type() != Object.StrObj {
			return &Object.Err{Msg: "StrReplace wants three Str arguments"}
		}
		Val := strings.ReplaceAll(Args[0].(*Object.Str).Value, Args[1].(*Object.Str).Value, Args[2].(*Object.Str).Value)
		return &Object.Str{Value: Val}
	},
}