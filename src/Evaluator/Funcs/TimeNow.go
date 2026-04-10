package Funcs

import (
	"time"
	"webonly/src/Object"
)

var TimeNow = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		return &Object.Str{Value: time.Now().Format(time.RFC3339)}
	},
}