package Funcs

import (
	"strings"
	"webonly/src/Object"
)

var StrLower = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		if len(Args) != 1 || Args[0].Type() != Object.StrObj {
			return &Object.Err{Msg: "StrLower wants one Str argument"}
		}
		return &Object.Str{Value: strings.ToLower(Args[0].(*Object.Str).Value)}
	},
}