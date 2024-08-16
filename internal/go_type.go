package golang

import (
	"github.com/debugger84/sqlc-graphql/internal/opts"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

func goInnerType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	columnType := sdk.DataType(col.Type)
	notNull := col.NotNull || col.IsArray

	// package overrides have a higher precedence
	for _, override := range options.Overrides {
		oride := override.ShimOverride
		if oride.GoType.TypeName == "" {
			continue
		}
		if oride.DbType != "" && oride.DbType == columnType && oride.Nullable != notNull && oride.Unsigned == col.Unsigned {
			return oride.GoType.TypeName
		}
	}

	// TODO: Extend the engine interface to handle types
	switch req.Settings.Engine {
	case "mysql":
		return mysqlType(req, options, col)
	case "postgresql":
		return postgresType(req, options, col)
	case "sqlite":
		return sqliteType(req, options, col)
	default:
		return "interface{}"
	}
}
