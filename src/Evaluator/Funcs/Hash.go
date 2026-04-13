package Funcs

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"webonly/src/Object"
)

var Hash = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		if len(Args) != 2 {
			return &Object.Err{Msg: "Hash requires 2 arguments"}
		}
		if Args[0].Type() != Object.StrObj || Args[1].Type() != Object.StrObj {
			return &Object.Err{Msg: "Hash arguments must be strings"}
		}

		Data := Args[0].(*Object.Str).Value
		Algo := Args[1].(*Object.Str).Value

		var HashBytes []byte

		switch Algo {
		case "md5":
			Sum := md5.Sum([]byte(Data))
			HashBytes = Sum[:]
		case "sha1":
			Sum := sha1.Sum([]byte(Data))
			HashBytes = Sum[:]
		case "sha256":
			Sum := sha256.Sum256([]byte(Data))
			HashBytes = Sum[:]
		case "sha512":
			Sum := sha512.Sum512([]byte(Data))
			HashBytes = Sum[:]
		default:
			return &Object.Err{Msg: "Unsupported hash algorithm"}
		}

		return &Object.Str{Value: hex.EncodeToString(HashBytes)}
	},
}