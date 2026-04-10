package Funcs

import "webonly/src/Object"

func GetAll() map[string]*Object.Builtin {
	return map[string]*Object.Builtin{
		"StrLen":     StrLen,
		"StrUpper":   StrUpper,
		"StrLower":   StrLower,
		"StrSplit":   StrSplit,
		"StrReplace": StrReplace,
		"TimeNow":    TimeNow,
		"TimeFmt":    TimeFmt,
		"ArrPush":    ArrPush,
		"ArrLen":     ArrLen,
	}
}