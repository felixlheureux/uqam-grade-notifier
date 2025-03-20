package db

import (
	"fmt"
	"reflect"

	"github.com/uptrace/bun"
)

// Operator symbols mapping
var operatorSymbols = map[string]string{
	"gte":      ">=",
	"gt":       ">",
	"lte":      "<=",
	"lt":       "<",
	"any":      "IS NOT NULL",
	"contains": "LIKE",
}

// isZeroValue checks if the value is a zero value of its type
func isZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		return val.IsNil()
	}

	return reflect.DeepEqual(v, reflect.Zero(val.Type()).Interface())
}

// AppendFiltersToQuery dynamically builds the SQL query with filters
func AppendFiltersToQuery(query *bun.SelectQuery, tableName string, filters map[string]interface{}) *bun.SelectQuery {
	for key, value := range filters {
		columnName := fmt.Sprintf("%s.%s", tableName, key)

		if !isZeroValue(value) {
			switch v := value.(type) {
			case string, int, int64, float64, bool:
				query = query.Where("? = ?", bun.Ident(columnName), v)
			case []interface{}:
				if len(v) > 0 {
					query = query.Where("? IN (?)", bun.Ident(columnName), bun.In(v))
				}
			case map[string]interface{}:
				for operator, val := range v {
					if symbol, exists := operatorSymbols[operator]; exists && !isZeroValue(val) {
						switch operator {
						case "any":
							query = query.Where("? IS NOT NULL", bun.Ident(columnName))
						case "contains":
							query = query.Where("LOWER(?) LIKE LOWER(?)", bun.Ident(columnName), fmt.Sprintf("%%%v%%", val))
						default:
							query = query.Where("? "+symbol+" ?", bun.Ident(columnName), val)
						}
					}
				}
			case *string, *int, *int64, *float64, *bool:
				query = query.Where("? = ?", bun.Ident(columnName), reflect.ValueOf(v).Elem().Interface())
			}
		}
	}

	return query
}
