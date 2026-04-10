package Funcs

import (
	"time"
	"webonly/src/Object"
)

var TimeFmt = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		if len(Args) != 1 || Args[0].Type() != Object.StrObj {
			return &Object.Err{Msg: "TimeFmt wants one Str argument"}
		}
		return &Object.Str{Value: time.Now().Format(Args[0].(*Object.Str).Value)}
	},
}