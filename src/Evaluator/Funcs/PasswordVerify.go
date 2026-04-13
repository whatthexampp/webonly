package Funcs

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"webonly/src/Object"
)

var PasswordVerify = &Object.Builtin{
	Fn: func(Args ...Object.Object) Object.Object {
		if len(Args) != 2 {
			return &Object.Err{Msg: "PasswordVerify requires 2 arguments"}
		}
		if Args[0].Type() != Object.StrObj || Args[1].Type() != Object.StrObj {
			return &Object.Err{Msg: "PasswordVerify arguments must be strings"}
		}

		Password := Args[0].(*Object.Str).Value
		HashStr := Args[1].(*Object.Str).Value

		if strings.HasPrefix(HashStr, "$2a$") || strings.HasPrefix(HashStr, "$2b$") || strings.HasPrefix(HashStr, "$2y$") {
			Err := bcrypt.CompareHashAndPassword([]byte(HashStr), []byte(Password))
			return &Object.Bool{Value: Err == nil}
		}

		if strings.HasPrefix(HashStr, "$argon2id$") {
			Parts := strings.Split(HashStr, "$")
			if len(Parts) != 6 {
				return &Object.Bool{Value: false}
			}

			var Version int
			_, Err := fmt.Sscanf(Parts[2], "v=%d", &Version)
			if Err != nil || Version != argon2.Version {
				return &Object.Bool{Value: false}
			}

			var Memory uint32
			var Time uint32
			var Threads uint8
			_, Err = fmt.Sscanf(Parts[3], "m=%d,t=%d,p=%d", &Memory, &Time, &Threads)
			if Err != nil {
				return &Object.Bool{Value: false}
			}

			Salt, Err := base64.RawStdEncoding.DecodeString(Parts[4])
			if Err != nil {
				return &Object.Bool{Value: false}
			}

			DecodedHash, Err := base64.RawStdEncoding.DecodeString(Parts[5])
			if Err != nil {
				return &Object.Bool{Value: false}
			}

			ComputedHash := argon2.IDKey([]byte(Password), Salt, Time, Memory, Threads, uint32(len(DecodedHash)))

			if subtle.ConstantTimeCompare(DecodedHash, ComputedHash) == 1 {
				return &Object.Bool{Value: true}
			}
			return &Object.Bool{Value: false}
		}

		return &Object.Bool{Value: false}
	},
}