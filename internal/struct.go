package golang

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/debugger84/sqlc-graphql/internal/opts"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

type Struct struct {
	Table     *plugin.Identifier
	Name      string
	Fields    []Field
	Comment   string
	ModelPath string
}

func StructName(name string, options *opts.Options) string {
	if rename := options.Rename[name]; rename != "" {
		return rename
	}
	out := ""
	name = strings.Map(
		func(r rune) rune {
			if unicode.IsLetter(r) {
				return r
			}
			if unicode.IsDigit(r) {
				return r
			}
			return rune('_')
		}, name,
	)

	for _, p := range strings.Split(name, "_") {
		out += strings.Title(p)
	}

	// If a name has a digit as its first char, prepand an underscore to make it a valid Go name.
	r, _ := utf8.DecodeRuneInString(out)
	if unicode.IsDigit(r) {
		return "_" + out
	} else {
		return out
	}
}
