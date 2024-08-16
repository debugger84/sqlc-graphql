package golang

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/debugger84/sqlc-graphql/internal/opts"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
	"slices"
	"strings"
	"text/template"
)

type gqlTmplCtx struct {
	ModelPackage string
	Enums        []Enum
	Structs      []Struct
	GoQueries    []Query
	SqlcVersion  string

	// TODO: Race conditions
	SourceName string

	OmitSqlcVersion bool
}

func (t *gqlTmplCtx) OutputQuery(sourceName string) bool {
	return t.SourceName == sourceName
}
func (t *gqlTmplCtx) ParamsName(InputName string) string {
	return strings.TrimRight(InputName, "Input") + "Params"
}

func generateGql(
	req *plugin.GenerateRequest,
	options *opts.Options,
	enums []Enum,
	structs []Struct,
	queries []Query,
) (*plugin.GenerateResponse, error) {
	excludedFields, err := getGqlExcluded(options)
	if err != nil {
		return nil, err
	}
	structs = filterStructs(structs, excludedFields)

	tctx := gqlTmplCtx{
		ModelPackage:    options.Package,
		Enums:           enums,
		Structs:         structs,
		SqlcVersion:     req.SqlcVersion,
		OmitSqlcVersion: options.OmitSqlcVersion,
		GoQueries:       queries,
	}

	funcMap := template.FuncMap{
		"lowerTitle": sdk.LowerTitle,
		"hasPrefix":  strings.HasPrefix,
	}

	tmpl := template.Must(
		template.New("table").
			Funcs(funcMap).
			ParseFS(
				templates,
				"templates/*.tmpl",
			),
	)

	output := map[string]string{}

	execute := func(name, templateName string) error {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		err := tmpl.ExecuteTemplate(w, templateName, &tctx)
		w.Flush()
		if err != nil {
			return err
		}

		if !strings.HasSuffix(name, ".graphql") {
			name = strings.TrimSuffix(name, ".sql")
			name += ".graphql"
		}
		output[name] = string(b.Bytes())
		return nil
	}

	gqlFileName := "schema.graphql"

	if err := execute(gqlFileName, "modelsGqlFile"); err != nil {
		return nil, err
	}

	if options.GenCommonParts {
		if err := execute("common.graphql", "commonGqlFile"); err != nil {
			return nil, err
		}
	}
	files := map[string]struct{}{}
	for _, gq := range queries {
		files[gq.SourceName] = struct{}{}
	}

	for source := range files {
		if err := execute(source, "gqlQueryFile"); err != nil {
			return nil, err
		}
	}
	resp := plugin.GenerateResponse{}

	for filename, code := range output {
		//if options.Out != "" {
		//	filename = options.Out + "/" + filename
		//}
		resp.Files = append(
			resp.Files, &plugin.File{
				Name:     filename,
				Contents: []byte(code),
			},
		)
	}

	return &resp, nil
}

func getGqlExcluded(options *opts.Options) (map[string][]string, error) {
	res := make(map[string][]string)
	if options == nil || options.Exclude == nil {
		return nil, nil
	}
	for _, exclude := range options.Exclude {
		parts := strings.Split(exclude, ".")

		if len(parts) == 0 {
			continue
		}

		if len(parts) != 2 {
			return nil, errors.New("invalid exclude format. It should be in the format of 'GqlTypeName.fieldName'")
		}
		typeName := parts[0]
		fieldName := parts[1]
		if _, ok := res[typeName]; !ok {
			res[typeName] = make([]string, 0)
		}
		res[typeName] = append(res[typeName], fieldName)
	}
	return res, nil
}

func filterStructs(structs []Struct, excludeFields map[string][]string) []Struct {
	var result []Struct
	for _, s := range structs {
		var fields []Field
		for _, f := range s.Fields {
			if _, ok := excludeFields[s.Name]; ok {
				if slices.Contains(excludeFields[s.Name], f.Name) {
					continue
				}
			}
			fields = append(fields, f)
		}
		s.Fields = fields
		result = append(result, s)
	}
	return result
}
