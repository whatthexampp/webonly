package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"webonly/src/Config"
	"webonly/src/Evaluator"
	"webonly/src/Evaluator/Funcs"
	"webonly/src/Lexer"
	"webonly/src/Object"
	"webonly/src/Parser"
	"webonly/src/Reporter"
)

func main() {
	Addr := "127.0.0.1:43211"
	L, Err := net.Listen("tcp", Addr)
	if Err != nil {
		fmt.Printf("Startup Error: %s\n", Err.Error())
		return
	}

	for {
		Conn, AccErr := L.Accept()
		if AccErr == nil {
			go Handle(Conn)
		}
	}
}

func Handle(C net.Conn) {
	defer C.Close()
	Buf := make([]byte, 4096)
	N, ReadErr := C.Read(Buf)
	if ReadErr != nil {
		return
	}

	ReqLine := strings.TrimSpace(string(Buf[:N]))
	Parts := strings.SplitN(ReqLine, "?", 2)
	ReqPath := Parts[0]
	QueryString := ""
	if len(Parts) > 1 {
		QueryString = Parts[1]
	}

	CfgPath := "/etc/webonly/Webonly.json"
	if runtime.GOOS == "windows" {
		CfgPath = "C:\\webonly\\Webonly.json"
	}
	Cfg := Config.Load(CfgPath)

	if Cfg.HideExtensions && !strings.HasSuffix(ReqPath, ".wo") {
		ReqPath += ".wo"
	}

	Path := ReqPath
	Src, FileErr := os.ReadFile(Path)
	if FileErr != nil {
		C.Write([]byte(fmt.Sprintf("Error: File Not Found -> %s", Path)))
		return
	}

	Code := string(Src)
	Lex := Lexer.Create(Code)
	Par := Parser.Create(Lex)
	Prog := Par.ParseProg()
	Rep := Reporter.Create(Code, Path)

	if len(Par.Errs) > 0 {
		C.Write([]byte(Rep.FmtParseErrs(Par.Errs)))
		return
	}

	Env := Object.NewEnv()

	PrintFn := &Object.Builtin{
		Fn: func(Args ...Object.Object) Object.Object {
			for _, A := range Args {
				C.Write([]byte(A.Inspect()))
			}
			return &Object.Null{}
		},
	}

	QueryFn := &Object.Builtin{
		Fn: func(Args ...Object.Object) Object.Object {
			if len(Args) != 1 || Args[0].Type() != Object.StrObj {
				return &Object.Null{}
			}
			Key := Args[0].(*Object.Str).Value
			Vals, _ := url.ParseQuery(QueryString)
			if Val := Vals.Get(Key); Val != "" {
				return &Object.Str{Value: Val}
			}
			return &Object.Null{}
		},
	}

	Env.SetConst("echo", PrintFn)
	Env.SetConst("__webonly_html_out", PrintFn)
	Env.SetConst("HttpQuery", QueryFn)

	for K, V := range Funcs.GetAll() {
		Env.SetConst(K, V)
	}

	MainEnv := Object.NewEncEnv(Env)
	AbsPath, _ := filepath.Abs(Path)
	MainEnv.SetConst("__FILE__", &Object.Str{Value: AbsPath})

	Res := Evaluator.Eval(Prog, MainEnv)

	if Evaluator.IsErr(Res) {
		Msg := Res.(*Object.Err).Msg
		Line := 0
		fmt.Sscanf(Msg, "Line %d:", &Line)
		C.Write([]byte(Rep.FmtErr(Msg, Line)))
	}
}