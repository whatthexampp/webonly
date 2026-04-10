package Funcs

import (
	"strings"
	"webonly/src/Object"
)

var StrSplit = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		if len(Args) != 2 || Args[0].Type() != Object.StrObj || Args[1].Type() != Object.StrObj {
			return &Object.Err{Msg: "StrSplit wants two Str arguments"}
		}
		Parts := strings.Split(Args[0].(*Object.Str).Value, Args[1].(*Object.Str).Value)
		Arr := &Object.Array{Elems: []Object.Object{}}
		for _, P := range Parts {
			Arr.Elems = append(Arr.Elems, &Object.Str{Value: P})
		}
		return Arr
	},
}