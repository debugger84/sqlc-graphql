package golang_test

import (
	"context"
	"encoding/json"
	golang "github.com/debugger84/sqlc-graphql/internal"
	"github.com/debugger84/sqlc-graphql/internal/opts"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerate(t *testing.T) {
	ctx := context.Background()
	t.Run(
		"Generate query", func(t *testing.T) {
			factory := NewGenReqFactory()
			req := factory.GenerateRequest()

			resp, err := golang.Generate(ctx, req)

			t.Log("Given the 'select * from authors' SQL query is passed to the generator")
			t.Log("When the generator is called")
			t.Log("	Then the generator should return a response without an error")
			require.NoError(t, err)
			t.Log("	And the response should contain the generated code")
			require.NotNil(t, resp)
			require.Len(t, resp.Files, 2)
			snaps.WithConfig(snaps.Ext("."+resp.Files[0].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[0].Contents))
			snaps.WithConfig(snaps.Ext("."+resp.Files[1].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[1].Contents))
		},
	)

	t.Run(
		"Generate paginated query", func(t *testing.T) {
			factory := NewGenReqFactory()
			req := factory.GenerateRequest()
			req.Queries[0].Comments = append(
				req.Queries[0].Comments,
				"gql: Query.authorsPaginated",
				"paginated:cursor:name,id",
			)

			resp, err := golang.Generate(ctx, req)

			t.Log("Given the 'select * from authors' SQL query is passed to the generator")
			t.Log("Given the comments with gql and paginated are added to the query")
			t.Log("When the generator is called")
			t.Log("	Then the generator should return a response without an error")
			require.NoError(t, err)
			t.Log("	And the response should contain the generated code")
			require.NotNil(t, resp)
			require.Len(t, resp.Files, 2)
			snaps.WithConfig(snaps.Ext("."+resp.Files[0].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[0].Contents))
			snaps.WithConfig(snaps.Ext("."+resp.Files[1].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[1].Contents))
		},
	)
	t.Run(
		"Generate offset pagination query", func(t *testing.T) {
			factory := NewGenReqFactory()
			req := factory.GenerateRequest()
			columns := req.Queries[0].Columns
			req.Queries[0] = &plugin.Query{
				Text:    "select id, name, status from authors",
				Name:    "PaginatedAuthors",
				Cmd:     ":many",
				Columns: columns,
				Params:  []*plugin.Parameter{},
				Comments: []string{
					"gql: Query.paginatedAuthors",
					"gql: Query.authorsPaginated",
					"paginated:offset",
				},
				Filename: "authors.sql",
			}

			resp, err := golang.Generate(ctx, req)

			t.Log("Given the 'select * from authors' SQL query is passed to the generator")
			t.Log("Given the comments with gql and paginated are added to the query")
			t.Log("When the generator is called")
			t.Log("	Then the generator should return a response without an error")
			require.NoError(t, err)
			t.Log("	And the response should contain the generated code")
			require.NotNil(t, resp)
			require.Len(t, resp.Files, 2)
			snaps.WithConfig(snaps.Ext("."+resp.Files[0].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[0].Contents))
			snaps.WithConfig(snaps.Ext("."+resp.Files[1].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[1].Contents))
		},
	)

	t.Run(
		"Exclude field", func(t *testing.T) {
			factory := NewGenReqFactory()
			factory.options.Exclude = []string{"Author.name"}
			req := factory.GenerateRequest()

			resp, err := golang.Generate(ctx, req)

			t.Log("Given the option to exclude the 'name' field is passed to the generator")
			t.Log("When the generator is called")
			t.Log("	Then the generator should return a response without an error")
			require.NoError(t, err)
			t.Log("	And the response should not contain the name field in the generated Author type")
			require.NotNil(t, resp)
			require.Len(t, resp.Files, 2)
			for _, file := range resp.Files {
				if file.Name == "schema.graphql" {
					snaps.WithConfig(snaps.Ext("."+file.Name)).
						MatchStandaloneSnapshot(t, string(file.Contents))
					require.NotContains(t, string(file.Contents), "name:")
				}
			}
		},
	)

	t.Run(
		"Add directive to query", func(t *testing.T) {
			factory := NewGenReqFactory()
			factory.options.Directives = []opts.Directive{
				{
					Model:     "Query",
					Field:     "author",
					Directive: "authGuard",
				},
				{
					Model:     "Author",
					Field:     "name",
					Directive: "authGuard",
				},
			}
			req := factory.GenerateRequest()

			resp, err := golang.Generate(ctx, req)

			t.Log("Given the options with adding directives to the query and the Author type are passed to the generator")
			t.Log("When the generator is called")
			t.Log("	Then the generator should return a response without an error")
			require.NoError(t, err)
			t.Log("	And the response should contain authGuard directive in the generated code")
			require.NotNil(t, resp)
			require.Len(t, resp.Files, 2)
			snaps.WithConfig(snaps.Ext("."+resp.Files[0].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[0].Contents))
			snaps.WithConfig(snaps.Ext("."+resp.Files[1].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[1].Contents))

			for _, file := range resp.Files {
				if file.Name == "schema.graphql" {
					require.Contains(
						t,
						string(file.Contents),
						"name: String @authGuard",
						"The name field of Author type should contain authGuard directive",
					)
				} else {
					require.Contains(
						t,
						string(file.Contents),
						"author(id: UUID!): Author! @authGuard",
						"The author query should contain authGuard directive",
					)
				}
			}
		},
	)
	t.Run(
		"Overwrite SQL type by custom Gql type", func(t *testing.T) {
			factory := NewGenReqFactory()
			factory.options.Overrides = []opts.Override{
				{
					DBType:  "text",
					GqlType: "MyText",
				},
				{
					DBType:  "public.status",
					GqlType: "MyStatus",
				},
			}
			req := factory.GenerateRequest()

			resp, err := golang.Generate(ctx, req)

			t.Log("Given the options with adding directives to the query and the Author type are passed to the generator")
			t.Log("When the generator is called")
			t.Log("	Then the generator should return a response without an error")
			require.NoError(t, err)
			t.Log("	And the response should contain authGuard directive in the generated code")
			require.NotNil(t, resp)
			require.Len(t, resp.Files, 2)
			snaps.WithConfig(snaps.Ext("."+resp.Files[0].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[0].Contents))
			snaps.WithConfig(snaps.Ext("."+resp.Files[1].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[1].Contents))

			for _, file := range resp.Files {
				if file.Name == "schema.graphql" {
					require.Contains(
						t,
						string(file.Contents),
						"name: MyText",
						"The name field should be of type MyText",
					)
					require.Contains(
						t,
						string(file.Contents),
						"status: MyStatus",
						"The status field should be of type MyStatus",
					)
				} else {
					require.Contains(
						t,
						string(file.Contents),
						"author(id: UUID!): Author!",
						"The author query should return Author",
					)
				}
			}
		},
	)
	t.Run(
		"Overwrite Gql type of a column", func(t *testing.T) {
			factory := NewGenReqFactory()
			factory.options.Overrides = []opts.Override{
				{
					Column:  "authors.name",
					GqlType: "MyText",
				},
			}
			req := factory.GenerateRequest()

			resp, err := golang.Generate(ctx, req)

			t.Log("Given the options with adding directives to the query and the Author type are passed to the generator")
			t.Log("When the generator is called")
			t.Log("	Then the generator should return a response without an error")
			require.NoError(t, err)
			t.Log("	And the response should contain authGuard directive in the generated code")
			require.NotNil(t, resp)
			require.Len(t, resp.Files, 2)
			snaps.WithConfig(snaps.Ext("."+resp.Files[0].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[0].Contents))
			snaps.WithConfig(snaps.Ext("."+resp.Files[1].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[1].Contents))

			for _, file := range resp.Files {
				if file.Name == "schema.graphql" {
					require.Contains(
						t,
						string(file.Contents),
						"name: MyText",
						"The name field should be of type MyText",
					)
				}
			}
		},
	)

	t.Run(
		"Generate mutation", func(t *testing.T) {

			factory := NewGenReqFactory()
			factory.query.Text = "insert into authors (id, name, status) values ($1, $2, $3) returning id, name, status"
			factory.query.Name = "InsertAuthor"
			factory.query.Cmd = ":exec"
			factory.query.Comments = []string{
				"gql: Mutation.createAuthor",
			}

			req := factory.GenerateRequest()
			resp, err := golang.Generate(ctx, req)

			t.Log("Given the insert SQL query is passed to the generator")
			t.Log("When the generator is called")
			t.Log("	Then the generator should return a response without an error")
			require.NoError(t, err)
			t.Log("	And the response should contain the generated code")
			require.NotNil(t, resp)
			require.Len(t, resp.Files, 2)
			snaps.WithConfig(snaps.Ext("."+resp.Files[0].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[0].Contents))
			snaps.WithConfig(snaps.Ext("."+resp.Files[1].Name)).
				MatchStandaloneSnapshot(t, string(resp.Files[1].Contents))
		},
	)
}

type genReqFactory struct {
	engine     string
	schemaName string
	tableIdent *plugin.Identifier
	columns    []*plugin.Column
	catalog    *plugin.Catalog
	query      *plugin.Query
	options    opts.Options
}

func NewGenReqFactory() genReqFactory {
	engine := "postgresql"
	schemaName := "public"
	tableIdent := &plugin.Identifier{
		Catalog: "",
		Schema:  schemaName,
		Name:    "authors",
	}
	columns := getDefaultColumns(tableIdent)
	return genReqFactory{
		engine:     engine,
		schemaName: schemaName,
		tableIdent: tableIdent,
		columns:    columns,
		catalog:    getDefaultCatalog(tableIdent, schemaName, columns),
		query:      getDefaultQuery(columns),
		options:    getDefaultOptions(schemaName),
	}
}

func getDefaultOptions(schemaName string) opts.Options {
	return opts.Options{
		Package:                     "authors/storage",
		Out:                         "./",
		Overrides:                   []opts.Override{},
		Rename:                      nil,
		OutputModelsFileName:        "",
		OutputQuerierFileName:       "",
		OutputCopyfromFileName:      "",
		OutputFilesSuffix:           "",
		EmitExactTableNames:         false,
		InflectionExcludeTableNames: nil,
		QueryParameterLimit:         nil,
		OmitSqlcVersion:             false,
		OmitUnusedStructs:           false,
		DefaultSchema:               schemaName,
		GenCommonParts:              false,
		Exclude:                     nil,
		Directives:                  nil,
	}
}

func getDefaultQuery(columns []*plugin.Column) *plugin.Query {
	return &plugin.Query{
		Text:    "select id, name, status from authors where id = $1",
		Name:    "GetAuthor",
		Cmd:     ":one",
		Columns: columns,
		Params: []*plugin.Parameter{
			{
				Column: columns[0],
			},
		},
		Comments: []string{
			"gql: Query.author",
		},
		Filename: "authors.sql",
	}
}

func (f genReqFactory) SetEngine(engine string) genReqFactory {
	f.engine = engine
	return f
}

func (f genReqFactory) SetSchemaName(schemaName string) genReqFactory {
	oldSchemaName := f.schemaName
	f.schemaName = schemaName
	f.tableIdent.Schema = schemaName
	for _, col := range f.columns {
		if col.Type.Schema == oldSchemaName {
			col.Type.Schema = schemaName
		}
	}
	f.options.DefaultSchema = schemaName
	return f
}

func getDefaultColumns(tableIdent *plugin.Identifier) []*plugin.Column {
	return []*plugin.Column{
		{
			Name:    "id",
			NotNull: true,
			Table:   tableIdent,
			Type: &plugin.Identifier{
				Name: "uuid",
			},
		},
		{
			Name:    "name",
			NotNull: false,
			Table:   tableIdent,
			Type: &plugin.Identifier{
				Name: "text",
			},
		},
		{
			Name:    "status",
			NotNull: true,
			Table:   tableIdent,
			Type: &plugin.Identifier{
				Name:   "status",
				Schema: tableIdent.Schema,
			},
		},
	}
}

func getDefaultCatalog(tableIdent *plugin.Identifier, schemaName string, columns []*plugin.Column) *plugin.Catalog {
	return &plugin.Catalog{
		DefaultSchema: schemaName,
		Schemas: []*plugin.Schema{
			{
				Comment: "",
				Name:    schemaName,
				Tables: []*plugin.Table{
					{
						Rel:     tableIdent,
						Columns: columns,
						Comment: "Authors",
					},
				},
				Enums: []*plugin.Enum{
					{
						Name: "status",
						Vals: []string{"active", "inactive"},
					},
				},
				CompositeTypes: nil,
			},
		},
	}
}

func (f genReqFactory) GenerateRequest() *plugin.GenerateRequest {
	req := &plugin.GenerateRequest{}

	pluginOptions := f.options

	jsonOpts, err := json.Marshal(&pluginOptions)
	if err != nil {
		panic(err)
	}

	settings := &plugin.Settings{
		Engine: f.engine,
	}

	req.Catalog = f.catalog
	req.Queries = []*plugin.Query{f.query}
	req.SqlcVersion = "v1.27.0"
	req.PluginOptions = jsonOpts
	req.Settings = settings

	return req
}
