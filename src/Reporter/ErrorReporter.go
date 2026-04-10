package Reporter

import (
	"fmt"
	"strings"
)

type Reporter struct {
	Code string
	Path string
}
//TOp ten besst error handlers in ALL of programming.

func Create(Code string, Path string) *Reporter {
	return &Reporter{Code: Code, Path: Path}
}

func (R *Reporter) FmtErr(Msg string, Line int) string {
	Lines := strings.Split(R.Code, "\n")
	Snip := ""
	if Line > 0 && Line <= len(Lines) {
		Snip = Lines[Line-1]
	}
 ///TODO: TURN THIS INTO AN HTML TEMPLATE
	Fmt := "<div style='background:#f8d7da; color:#721c24; padding:15px; border-radius:5px; font-family:monospace;'>"
	Fmt += "<b>Webonly Error</b><br>"
	Fmt += fmt.Sprintf("<b>File:</b> %s<br>", R.Path)
	Fmt += fmt.Sprintf("<b>Line:</b> %d<br><br>", Line)
	Fmt += fmt.Sprintf("<b>Details:</b> %s<br><br>", Msg)
	if Snip != "" {
		Fmt += "<b>Snippet:</b><br>"
		Fmt += fmt.Sprintf("<pre style='background:#f1b0b7; padding:10px;'>%s</pre>", strings.TrimSpace(Snip))
	}
	Fmt += "</div>"
	return Fmt
}

func (R *Reporter) FmtParseErrs(Errs []string) string {
	Fmt := "<div style='background:#f8d7da; color:#721c24; padding:15px; border-radius:5px; font-family:monospace;'>"
	Fmt += fmt.Sprintf("<b>Parse Errors in:</b> %s<br><ul>", R.Path)
	for _, E := range Errs {
		Fmt += fmt.Sprintf("<li>%s</li>", E)
	}
	Fmt += "</ul></div>"
	return Fmt
}