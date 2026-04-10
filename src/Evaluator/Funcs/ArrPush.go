package Funcs

import "webonly/src/Object"

var ArrPush = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		if len(Args) != 2 || Args[0].Type() != Object.ArrObj {
			return &Object.Err{Msg: "ArrPush wants an Array and an Element"}
		}
		Arr := Args[0].(*Object.Array)
		Arr.Elems = append(Arr.Elems, Args[1])
		return &Object.Num{Value: float64(len(Arr.Elems))}
	},
}