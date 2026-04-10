package Funcs

import (
	"strings"
	"webonly/src/Object"
)

var StrUpper = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		if len(Args) != 1 || Args[0].Type() != Object.StrObj {
			return &Object.Err{Msg: "StrUpper wants one Str argument"}
		}
		return &Object.Str{Value: strings.ToUpper(Args[0].(*Object.Str).Value)}
	},
}