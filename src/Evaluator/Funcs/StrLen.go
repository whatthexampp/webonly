package Funcs

import "webonly/src/Object"

var StrLen = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		if len(Args) != 1 || Args[0].Type() != Object.StrObj {
			return &Object.Err{Msg: "StrLen wants one Str argument"}
		}
		return &Object.Num{Value: float64(len(Args[0].(*Object.Str).Value))}
	},
}