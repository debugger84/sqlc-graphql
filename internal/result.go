package golang

import (
	"fmt"
	"sort"
	"strings"

	"github.com/debugger84/sqlc-graphql/internal/inflection"
	"github.com/debugger84/sqlc-graphql/internal/opts"
	"github.com/sqlc-dev/plugin-sdk-go/metadata"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

func buildEnums(req *plugin.GenerateRequest, options *opts.Options) []Enum {
	var enums []Enum
	for _, schema := range req.Catalog.Schemas {
		if schema.Name == "pg_catalog" || schema.Name == "information_schema" {
			continue
		}
		for _, enum := range schema.Enums {
			var enumName string
			if schema.Name == req.Catalog.DefaultSchema {
				enumName = enum.Name
			} else {
				enumName = schema.Name + "_" + enum.Name
			}

			e := Enum{
				Name:    StructName(enumName, options),
				Comment: enum.Comment,
			}

			seen := make(map[string]struct{}, len(enum.Vals))
			for i, v := range enum.Vals {
				value := EnumReplace(v)
				if _, found := seen[value]; found || value == "" {
					value = fmt.Sprintf("value_%d", i)
				}
				e.Constants = append(
					e.Constants, Constant{
						Name:  StructName(enumName+"_"+value, options),
						Value: v,
						Type:  e.Name,
					},
				)
				seen[value] = struct{}{}
			}
			enums = append(enums, e)
		}
	}
	if len(enums) > 0 {
		sort.Slice(enums, func(i, j int) bool { return enums[i].Name < enums[j].Name })
	}
	return enums
}

func buildStructs(req *plugin.GenerateRequest, options *opts.Options) []Struct {
	var structs []Struct
	for _, schema := range req.Catalog.Schemas {
		if schema.Name == "pg_catalog" || schema.Name == "information_schema" {
			continue
		}
		for _, table := range schema.Tables {
			var tableName string
			if schema.Name == req.Catalog.DefaultSchema {
				tableName = table.Rel.Name
			} else {
				tableName = schema.Name + "_" + table.Rel.Name
			}
			structName := tableName
			if !options.EmitExactTableNames {
				structName = inflection.Singular(
					inflection.SingularParams{
						Name:       structName,
						Exclusions: options.InflectionExcludeTableNames,
					},
				)
			}
			modelName := StructName(structName, options)
			s := Struct{
				Table:   &plugin.Identifier{Schema: schema.Name, Name: table.Rel.Name},
				Name:    modelName,
				Comment: table.Comment,
			}
			s.ModelPath = options.Package + "." + s.Name
			for _, column := range table.Columns {
				fieldName := StructName(column.Name, options)
				s.Fields = append(
					s.Fields, Field{
						Name:      fieldName,
						Type:      gqlType(req, options, column),
						Comment:   column.Comment,
						Directive: parseDirective(options.Directives, modelName, fieldName),
					},
				)
			}
			structs = append(structs, s)
		}
	}
	if len(structs) > 0 {
		sort.Slice(structs, func(i, j int) bool { return structs[i].Name < structs[j].Name })
	}
	return structs
}

type goColumn struct {
	id int
	*plugin.Column
	embed *goEmbed
}

type goEmbed struct {
	modelType string
	modelName string
	fields    []Field
}

// look through all the structs and attempt to find a matching one to embed
// We need the name of the struct and its field names.
func newGoEmbed(embed *plugin.Identifier, structs []Struct, defaultSchema string) *goEmbed {
	if embed == nil {
		return nil
	}

	for _, s := range structs {
		embedSchema := defaultSchema
		if embed.Schema != "" {
			embedSchema = embed.Schema
		}

		// compare the other attributes
		if embed.Catalog != s.Table.Catalog || embed.Name != s.Table.Name || embedSchema != s.Table.Schema {
			continue
		}

		fields := make([]Field, len(s.Fields))
		for i, f := range s.Fields {
			fields[i] = f
		}

		return &goEmbed{
			modelType: s.Name,
			modelName: s.Name,
			fields:    fields,
		}
	}

	return nil
}

func columnName(c *plugin.Column, pos int) string {
	if c.Name != "" {
		return c.Name
	}
	return fmt.Sprintf("column_%d", pos+1)
}

func paramName(p *plugin.Parameter) string {
	if p.Column.Name != "" {
		return argName(p.Column.Name)
	}
	return fmt.Sprintf("dollar_%d", p.Number)
}

func argName(name string) string {
	out := ""
	for i, p := range strings.Split(name, "_") {
		if i == 0 {
			out += strings.ToLower(p)
		} else if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}
	return out
}

func buildQueries(req *plugin.GenerateRequest, options *opts.Options, structs []Struct) ([]Query, error) {
	qs := make([]Query, 0, len(req.Queries))
	for _, query := range req.Queries {
		if query.Name == "" {
			continue
		}
		if query.Cmd == "" {
			continue
		}

		comments := query.Comments
		var extendedType string
		resolverName := query.Name
		returnType := ""
		directive := ""
		for i, comment := range comments {
			parts := strings.Split(comment, ":")
			if len(parts) > 1 && strings.Trim(parts[0], " ") == "gql" {
				gql := strings.Trim(parts[1], " ")
				resolverInfo := strings.Split(gql, ".")
				extendedType = resolverInfo[0]
				if len(resolverInfo) > 1 {
					resolverName = resolverInfo[1]
				}
				if len(parts) > 2 {
					returnType = strings.TrimSpace(parts[2])
					directiveParts := strings.Split(returnType, "@")
					if len(directiveParts) > 1 {
						returnType = strings.TrimSpace(directiveParts[0])
						directive = "@" + strings.Join(directiveParts[1:], "@")
					}
				}
				comments = append(comments[:i], comments[i+1:]...)
				break
			}
		}
		if extendedType == "" {
			continue
		}

		paginated := false
		cursorPagination := false
		for i, comment := range comments {
			comment = strings.TrimSpace(comment)
			if strings.HasPrefix(comment, "paginated") {
				paginated = true
				comments = append(comments[:i], comments[i+1:]...)
				if strings.Contains(comment, "cursor") {
					cursorPagination = true
				}
				break
			}
		}

		parsedDirective := parseDirective(options.Directives, extendedType, resolverName)
		if parsedDirective != "" {
			if directive != "" {
				directive += " " + parsedDirective
			} else {
				directive = parsedDirective
			}
		}
		gq := Query{
			Cmd:              query.Cmd,
			Comments:         comments,
			MethodName:       query.Name,
			SourceName:       query.Filename,
			ExtendedType:     extendedType,
			ResolverName:     resolverName,
			Directive:        directive,
			Paginated:        paginated,
			CursorPagination: cursorPagination,
		}

		if returnType == "" {
			returnType = gq.MethodName + "Row"
		}

		qpl := int(*options.QueryParameterLimit)

		if paginated {
			number := int32(len(query.Params) + 1)
			if cursorPagination {
				query.Params = append(
					query.Params, &plugin.Parameter{
						Number: number,
						Column: &plugin.Column{
							Name:         "first",
							NotNull:      true,
							IsNamedParam: true,
							Type: &plugin.Identifier{
								Name: "int",
							},
						},
					}, &plugin.Parameter{
						Number: number + 1,
						Column: &plugin.Column{
							Name:         "after",
							NotNull:      true,
							IsNamedParam: true,
							Type: &plugin.Identifier{
								Name: "string",
							},
						},
					},
				)
			} else {
				query.Params = append(
					query.Params, &plugin.Parameter{
						Number: number,
						Column: &plugin.Column{
							Name:         "limit",
							NotNull:      true,
							IsNamedParam: true,
							Type: &plugin.Identifier{
								Name: "int",
							},
						},
					}, &plugin.Parameter{
						Number: number + 1,
						Column: &plugin.Column{
							Name:         "offset",
							NotNull:      true,
							IsNamedParam: true,
							Type: &plugin.Identifier{
								Name: "int",
							},
						},
					},
				)
			}
		}

		if len(query.Params) == 1 && qpl != 0 {
			p := query.Params[0]
			gq.Arg = QueryValue{
				Name:   escape(paramName(p)),
				DBName: p.Column.GetName(),
				Typ:    gqlType(req, options, p.Column),
				Column: p.Column,
			}
		} else if len(query.Params) >= 1 {
			var cols []goColumn
			for _, p := range query.Params {
				cols = append(
					cols, goColumn{
						id:     int(p.Number),
						Column: p.Column,
					},
				)
			}
			s, err := columnsToStruct(req, options, resolverName+"Input", cols, false)
			if err != nil {
				return nil, err
			}
			s.Fields = addDefaultDirectivesToPaginationInputFields(s.Fields)
			gq.Arg = QueryValue{
				Emit:      true,
				Name:      "request",
				Struct:    s,
				ModelPath: options.Package + "." + gq.MethodName + "Params",
			}

			if len(query.Params) <= qpl {
				gq.Arg.Emit = false
			}
		}

		if len(query.Columns) == 1 && query.Columns[0].EmbedTable == nil {
			c := query.Columns[0]
			name := columnName(c, 0)
			name = strings.Replace(name, "$", "_", -1)
			gq.Ret = QueryValue{
				Name:      escape(name),
				DBName:    name,
				Typ:       gqlType(req, options, c),
				ModelPath: options.Package + "." + gq.MethodName,
			}
		} else if putOutColumns(query) {
			var gs *Struct
			var emit bool

			for _, s := range structs {
				if len(s.Fields) != len(query.Columns) {
					continue
				}
				same := true
				for i, f := range s.Fields {
					c := query.Columns[i]
					sameName := f.Name == StructName(columnName(c, i), options)
					sameType := f.Type == gqlType(req, options, c)
					sameTable := sdk.SameTableName(c.Table, s.Table, req.Catalog.DefaultSchema)
					if !sameName || !sameType || !sameTable {
						same = false
					}
				}
				if same {
					buf := s
					gs = &buf
					break
				}
			}

			modelPath := options.Package + "." + gq.MethodName
			if gs == nil {
				var columns []goColumn
				for i, c := range query.Columns {
					columns = append(
						columns, goColumn{
							id:     i,
							Column: c,
							embed:  newGoEmbed(c.EmbedTable, structs, req.Catalog.DefaultSchema),
						},
					)
				}
				var err error
				gs, err = columnsToStruct(req, options, returnType, columns, true)
				if err != nil {
					return nil, err
				}
				emit = true
				modelPath = options.Package + "." + gq.MethodName + "Row"
			}
			gs.ModelPath = modelPath

			gq.Ret = QueryValue{
				Emit:      emit,
				Name:      "i",
				Struct:    gs,
				Typ:       gs.Name + "!",
				ModelPath: modelPath,
			}
		}

		qs = append(qs, gq)
	}
	sort.Slice(qs, func(i, j int) bool { return qs[i].MethodName < qs[j].MethodName })
	return qs, nil
}

var cmdReturnsData = map[string]struct{}{
	metadata.CmdBatchMany: {},
	metadata.CmdBatchOne:  {},
	metadata.CmdMany:      {},
	metadata.CmdOne:       {},
}

func putOutColumns(query *plugin.Query) bool {
	_, found := cmdReturnsData[query.Cmd]
	return found
}

// It's possible that this method will generate duplicate JSON tag values
//
//	Columns: count, count,   count_2
//	 Fields: Count, Count_2, Count2
//
// JSON tags: count, count_2, count_2
//
// This is unlikely to happen, so don't fix it yet
func columnsToStruct(
	req *plugin.GenerateRequest,
	options *opts.Options,
	name string,
	columns []goColumn,
	useID bool,
) (*Struct, error) {
	gs := Struct{
		Name: name,
	}
	seen := map[string][]int{}
	suffixes := map[int]int{}
	for i, c := range columns {
		colName := columnName(c.Column, i)
		tagName := colName

		// override col/tag with expected model name
		if c.embed != nil {
			colName = c.embed.modelName
			tagName = SetCaseStyle(colName, "snake")
		}

		fieldName := StructName(colName, options)
		baseFieldName := fieldName
		// Track suffixes by the ID of the column, so that columns referring to the same numbered parameter can be
		// reused.
		suffix := 0
		if o, ok := suffixes[c.id]; ok && useID {
			suffix = o
		} else if v := len(seen[fieldName]); v > 0 && !c.IsNamedParam {
			suffix = v + 1
		}
		suffixes[c.id] = suffix
		if suffix > 0 {
			tagName = fmt.Sprintf("%s_%d", tagName, suffix)
			fieldName = fmt.Sprintf("%s_%d", fieldName, suffix)
		}

		f := Field{
			Name:   fieldName,
			DBName: colName,
			Column: c.Column,
		}
		f.Directive = parseDirective(options.Directives, name, fieldName)
		if c.embed == nil {
			f.Type = gqlType(req, options, c.Column)
		} else {
			f.Type = c.embed.modelType
			f.EmbedFields = c.embed.fields
		}

		gs.Fields = append(gs.Fields, f)
		if _, found := seen[baseFieldName]; !found {
			seen[baseFieldName] = []int{i}
		} else {
			seen[baseFieldName] = append(seen[baseFieldName], i)
		}
	}

	// If a field does not have a known type, but another
	// field with the same name has a known type, assign
	// the known type to the field without a known type
	for i, field := range gs.Fields {
		if len(seen[field.Name]) > 1 && field.Type == "interface{}" {
			for _, j := range seen[field.Name] {
				if i == j {
					continue
				}
				otherField := gs.Fields[j]
				if otherField.Type != field.Type {
					field.Type = otherField.Type
				}
				gs.Fields[i] = field
			}
		}
	}

	err := checkIncompatibleFieldTypes(gs.Fields)
	if err != nil {
		return nil, err
	}

	return &gs, nil
}

func checkIncompatibleFieldTypes(fields []Field) error {
	fieldTypes := map[string]string{}
	for _, field := range fields {
		if fieldType, found := fieldTypes[field.Name]; !found {
			fieldTypes[field.Name] = field.Type
		} else if field.Type != fieldType {
			return fmt.Errorf("named param %s has incompatible types: %s, %s", field.Name, field.Type, fieldType)
		}
	}
	return nil
}

func addRetValuesToStructs(structs []Struct, queries []Query) []Struct {
	for _, q := range queries {
		if q.Ret.Struct != nil {
			if q.Ret.Emit {
				structs = append(structs, *q.Ret.Struct)
			}
			if q.Paginated {
				if q.CursorPagination {
					structs = addConnectionStruct(*q.Ret.Struct, structs)
				} else {
					structs = addPageStruct(*q.Ret.Struct, structs)
				}
			}
		}
	}
	return structs
}

func addPageStruct(original Struct, structs []Struct) []Struct {
	pageName := original.Name + "Page"
	for _, s := range structs {
		if s.Name == pageName {
			return structs
		}
	}
	pageStruct := Struct{
		Name:      pageName,
		ModelPath: original.ModelPath + "Page",
		Fields: []Field{
			{
				Name:    "Items",
				DBName:  "",
				Type:    "[" + original.Name + "!]!",
				Comment: "",
				Column: &plugin.Column{
					Name:    "items",
					NotNull: true,
					IsArray: true,
				},
				EmbedFields: nil,
			},
			{
				Name:   "Total",
				DBName: "",
				Type:   "Int!",
			},
			{
				Name:   "HasNext",
				DBName: "",
				Type:   "Boolean!",
			},
		},
	}
	structs = append(structs, pageStruct)
	return structs
}

func addConnectionStruct(original Struct, structs []Struct) []Struct {
	connectionName := original.Name + "Connection"
	edgeName := original.Name + "Edge"
	for _, s := range structs {
		if s.Name == connectionName {
			return structs
		}
	}

	edgeStruct := Struct{
		Name:      edgeName,
		ModelPath: original.ModelPath + "Edge",
		Fields: []Field{
			{
				Name: "Node",
				Type: original.Name + "!",
				Column: &plugin.Column{
					Name:    "node",
					NotNull: true,
					IsArray: false,
					Type:    &plugin.Identifier{Name: original.Name},
				},
			},
			{
				Name: "Cursor",
				Type: "String!",
				Column: &plugin.Column{
					Name:    "cursor",
					NotNull: true,
					IsArray: false,
					Type:    &plugin.Identifier{Name: "string"},
				},
			},
		},
	}

	connectionStruct := Struct{
		Name:      connectionName,
		ModelPath: original.ModelPath + "Connection",
		Fields: []Field{
			{
				Name:    "Edges",
				DBName:  "",
				Type:    "[" + edgeName + "!]!",
				Comment: "",
				Column: &plugin.Column{
					Name:    "edges",
					NotNull: true,
					IsArray: true,
					Type:    &plugin.Identifier{Name: edgeName},
				},
				EmbedFields: nil,
			},
			{
				Name:   "PageInfo",
				DBName: "",
				Type:   "PageInfo!",
				Column: &plugin.Column{
					Name:    "pageInfo",
					NotNull: true,
					IsArray: false,
					Type:    &plugin.Identifier{Name: "PageInfo"},
				},
			},
		},
	}

	structs = append(structs, connectionStruct, edgeStruct)
	return structs
}

func parseDirective(directives []opts.Directive, modelName, fieldName string) string {
	res := make([]string, 0)
	for _, directive := range directives {
		if strings.EqualFold(modelName, directive.Model) {
			if directive.Field == "" || strings.EqualFold(fieldName, directive.Field) {
				dn := strings.TrimSpace(directive.Directive)
				if !strings.HasPrefix(dn, "@") {
					dn = "@" + dn
				}
				res = append(res, dn)
			}
		}
	}

	return strings.Join(res, " ")
}

func addDefaultDirectivesToPaginationInputFields(fields []Field) []Field {
	res := make([]Field, 0, len(fields))
	for _, f := range fields {
		if f.Name == "First" && !strings.Contains(f.Directive, "@goField") {
			f.Directive += " @goField(name: \"limit\")"
		}
		if f.Name == "After" && !strings.Contains(f.Directive, "@goField") {
			f.Directive += " @goField(name: \"cursor\")"
		}
		f.Directive = strings.TrimSpace(f.Directive)
		res = append(res, f)
	}
	return res
}
