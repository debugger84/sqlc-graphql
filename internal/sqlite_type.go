package golang

import (
	"log"
	"strings"

	"github.com/debugger84/sqlc-graphql/internal/debug"
	"github.com/debugger84/sqlc-graphql/internal/opts"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

func sqliteType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	dt := strings.ToLower(sdk.DataType(col.Type))
	notNull := col.NotNull || col.IsArray

	switch dt {

	case "int", "integer", "tinyint", "smallint", "mediumint", "bigint", "unsignedbigint", "int2", "int8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "blob":
		return "[]byte"

	case "real", "double", "doubleprecision", "float":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	case "boolean", "bool":
		if notNull {
			return "bool"
		}
		return "sql.NullBool"

	case "date", "datetime", "timestamp":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "any":
		return "interface{}"

	}

	switch {

	case strings.HasPrefix(dt, "character"),
		strings.HasPrefix(dt, "varchar"),
		strings.HasPrefix(dt, "varyingcharacter"),
		strings.HasPrefix(dt, "nchar"),
		strings.HasPrefix(dt, "nativecharacter"),
		strings.HasPrefix(dt, "nvarchar"),
		dt == "text",
		dt == "clob":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case strings.HasPrefix(dt, "decimal"), dt == "numeric":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	default:
		if debug.Active {
			log.Printf("unknown SQLite type: %s\n", dt)
		}

		return "interface{}"

	}
}
