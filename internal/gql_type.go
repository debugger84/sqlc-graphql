package golang

import (
	"fmt"
	"github.com/debugger84/sqlc-graphql/internal/opts"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
	"strings"
)

type GqlField struct {
	Name    string
	Type    string
	Comment string
}

type GqlStruct struct {
	Name    string
	Comment string
	Fields  []GqlField
}

func gqlType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	// Check if the column's type has been overridden
	for _, override := range options.Overrides {
		oride := override.ShimOverride

		if override.GqlType == "" {
			continue
		}
		cname := col.Name
		if col.OriginalName != "" {
			cname = col.OriginalName
		}
		sameTable := override.Matches(col.Table, req.Catalog.DefaultSchema)
		dbTypeParts := strings.Split(oride.DbType, ".")
		schema := ""
		typeName := dbTypeParts[len(dbTypeParts)-1]
		if len(dbTypeParts) > 1 {
			schema = dbTypeParts[0]
		}
		if (oride.Column != "" && sdk.MatchString(oride.ColumnName, cname) && sameTable) ||
			(col.Type != nil && col.Type.Name == typeName && col.Type.Schema == schema) {
			tn := override.GqlType
			if col.NotNull {
				tn += "!"
			}
			if col.IsSqlcSlice {
				tn = fmt.Sprintf("[%s]!", tn)
			}

			return tn
		}
	}
	typ := gqlInnerType(req, options, col)
	if col.NotNull {
		return typ + "!"
	}

	if col.IsSqlcSlice {
		return fmt.Sprintf("[%s]!", typ)
	}
	return typ
}

func gqlInnerType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	gotype := goInnerType(req, options, col)
	switch gotype {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "Int"
	case "float32", "float64":
		return "Float"
	case "bool":
		return "Boolean"
	case "string":
		return "String"
	case "time.Time", "sql.NullTime":
		return "Time"
	case "uuid.UUID", "uuid.NullUUID":
		return "UUID"

	case "netip.Addr", "*netip.Addr", "netip.Prefix", "*netip.Prefix", "net.HardwareAddr", "pgtype.Range[pgtype.Date]", "pgtype.Multirange[pgtype.Range[pgtype.Date]]":
		return "String"

	case "sql.NullInt8", "sql.NullInt16", "sql.NullInt32", "sql.NullInt64", "sql.NullUint", "sql.NullUint8", "sql.NullUint16", "sql.NullUint32", "sql.NullUint64":
		return "Int"
	case "sql.NullFloat32", "sql.NullFloat64":
		return "Float"
	case "sql.NullBool":
		return "Boolean"
	case "sql.NullString":
		return "String"
	case "[]byte", "pgtype.JSON", "pgtype.JSONB", "json.RawMessage", "pqtype.NullRawMessage":
		return "JSON"
	case "interface{}":
		return "Unknown"
	}

	tmpGqlType := gotype
	parts := strings.Split(gotype, ".")
	if len(parts) > 1 {
		tmpGqlType = parts[len(parts)-1]
	}
	tmpGqlType = strings.TrimPrefix(tmpGqlType, "Null")

	return tmpGqlType
}
