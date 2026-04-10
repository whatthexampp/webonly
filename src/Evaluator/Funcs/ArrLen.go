package Funcs

import "webonly/src/Object"

var ArrLen = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		if len(Args) != 1 || Args[0].Type() != Object.ArrObj {
			return &Object.Err{Msg: "ArrLen wants one Array argument"}
		}
		return &Object.Num{Value: float64(len(Args[0].(*Object.Array).Elems))}
	},
}