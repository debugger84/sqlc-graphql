package golang

import (
	"context"
	"fmt"
	"github.com/debugger84/sqlc-graphql/internal/opts"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

//type tmplCtx struct {
//	Q           string
//	Package     string
//	SQLDriver   opts.SQLDriver
//	Enums       []Enum
//	Structs     []Struct
//	GoQueries   []Query
//	SqlcVersion string
//
//	// TODO: Race conditions
//	SourceName string
//
//	EmitJSONTags              bool
//	JsonTagsIDUppercase       bool
//	EmitDBTags                bool
//	EmitPreparedQueries       bool
//	EmitInterface             bool
//	EmitEmptySlices           bool
//	EmitMethodsWithDBArgument bool
//	EmitEnumValidMethod       bool
//	EmitAllEnumValues         bool
//	UsesCopyFrom              bool
//	UsesBatch                 bool
//	OmitSqlcVersion           bool
//	BuildTags                 string
//}
//
//func (t *tmplCtx) OutputQuery(sourceName string) bool {
//	return t.SourceName == sourceName
//}
//
//func (t *tmplCtx) codegenDbarg() string {
//	if t.EmitMethodsWithDBArgument {
//		return "db DBTX, "
//	}
//	return ""
//}
//
//// Called as a global method since subtemplate queryCodeStdExec does not have
//// access to the toplevel tmplCtx
//func (t *tmplCtx) codegenEmitPreparedQueries() bool {
//	return t.EmitPreparedQueries
//}
//
//func (t *tmplCtx) codegenQueryMethod(q Query) string {
//	db := "q.db"
//	if t.EmitMethodsWithDBArgument {
//		db = "db"
//	}
//
//	switch q.Cmd {
//	case ":one":
//		if t.EmitPreparedQueries {
//			return "q.queryRow"
//		}
//		return db + ".QueryRowContext"
//
//	case ":many":
//		if t.EmitPreparedQueries {
//			return "q.query"
//		}
//		return db + ".QueryContext"
//
//	default:
//		if t.EmitPreparedQueries {
//			return "q.exec"
//		}
//		return db + ".ExecContext"
//	}
//}
//
//func (t *tmplCtx) codegenQueryRetval(q Query) (string, error) {
//	switch q.Cmd {
//	case ":one":
//		return "row :=", nil
//	case ":many":
//		return "rows, err :=", nil
//	case ":exec":
//		return "_, err :=", nil
//	case ":execrows", ":execlastid":
//		return "result, err :=", nil
//	case ":execresult":
//		return "return", nil
//	default:
//		return "", fmt.Errorf("unhandled q.Cmd case %q", q.Cmd)
//	}
//}

func Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	options, err := opts.Parse(req)
	if err != nil {
		return nil, err
	}

	if err := opts.ValidateOpts(options); err != nil {
		return nil, err
	}

	if options.DefaultSchema != "" {
		req.Catalog.DefaultSchema = options.DefaultSchema
	}
	enums := buildEnums(req, options)
	structs := buildStructs(req, options)
	queries, err := buildQueries(req, options, structs)
	if err != nil {
		return nil, err
	}

	if options.OmitUnusedStructs {
		enums, structs = filterUnusedStructs(enums, structs, queries)
	}

	if err := validate(enums, structs, queries); err != nil {
		return nil, err
	}

	resp, err := generateGql(req, options, enums, structs, queries)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func validate(enums []Enum, structs []Struct, queries []Query) error {
	enumNames := make(map[string]struct{})
	for _, enum := range enums {
		enumNames[enum.Name] = struct{}{}
		enumNames["Null"+enum.Name] = struct{}{}
	}
	structNames := make(map[string]struct{})
	for _, struckt := range structs {
		if _, ok := enumNames[struckt.Name]; ok {
			return fmt.Errorf("struct name conflicts with enum name: %s", struckt.Name)
		}
		structNames[struckt.Name] = struct{}{}
	}

	return nil
}

func filterUnusedStructs(enums []Enum, structs []Struct, queries []Query) ([]Enum, []Struct) {
	keepTypes := make(map[string]struct{})

	for _, query := range queries {
		if !query.Arg.isEmpty() {
			keepTypes[query.Arg.Type()] = struct{}{}
			if query.Arg.IsStruct() {
				for _, field := range query.Arg.Struct.Fields {
					keepTypes[field.Type] = struct{}{}
				}
			}
		}
		if query.hasRetType() {
			keepTypes[query.Ret.Type()] = struct{}{}
			if query.Ret.IsStruct() {
				for _, field := range query.Ret.Struct.Fields {
					keepTypes[field.Type] = struct{}{}
					for _, embedField := range field.EmbedFields {
						keepTypes[embedField.Type] = struct{}{}
					}
				}
			}
		}
	}

	keepEnums := make([]Enum, 0, len(enums))
	for _, enum := range enums {
		_, keep := keepTypes[enum.Name]
		_, keepNull := keepTypes["Null"+enum.Name]
		if keep || keepNull {
			keepEnums = append(keepEnums, enum)
		}
	}

	keepStructs := make([]Struct, 0, len(structs))
	for _, st := range structs {
		if _, ok := keepTypes[st.Name]; ok {
			keepStructs = append(keepStructs, st)
		}
	}

	return keepEnums, keepStructs
}
