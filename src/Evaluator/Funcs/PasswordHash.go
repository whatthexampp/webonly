package Funcs

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"webonly/src/Object"
)

var PasswordHash = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		if len(Args) < 1 || len(Args) > 2 {
			return &Object.Err{Msg: "PasswordHash requires 1 or 2 arguments"}
		}
		if Args[0].Type() != Object.StrObj {
			return &Object.Err{Msg: "PasswordHash requires string password"}
		}

		Password := Args[0].(*Object.Str).Value
		Algo := "argon2"

		if len(Args) == 2 {
			if Args[1].Type() != Object.StrObj {
				return &Object.Err{Msg: "PasswordHash requires string algorithm"}
			}
			Algo = Args[1].(*Object.Str).Value
		}

		switch Algo {
		case "bcrypt":
			Hash, Err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
			if Err != nil {
				return &Object.Err{Msg: "Bcrypt generation failed"}
			}
			return &Object.Str{Value: string(Hash)}
		case "argon2", "argon2id":
			Salt := make([]byte, 16)
			if _, Err := rand.Read(Salt); Err != nil {
				return &Object.Err{Msg: "Random source failed"}
			}
			Time := uint32(1)
			Memory := uint32(64 * 1024)
			Threads := uint8(2)
			KeyLen := uint32(32)

			Hash := argon2.IDKey([]byte(Password), Salt, Time, Memory, Threads, KeyLen)
			B64Salt := base64.RawStdEncoding.EncodeToString(Salt)
			B64Hash := base64.RawStdEncoding.EncodeToString(Hash)

			Res := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, Memory, Time, Threads, B64Salt, B64Hash)
			return &Object.Str{Value: Res}
		default:
			return &Object.Err{Msg: "Unsupported algorithm"}
		}
	},
}