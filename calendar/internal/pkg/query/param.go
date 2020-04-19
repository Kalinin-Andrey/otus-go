package query

import (
	"context"
)

// QueryParams const
const QueryParams = "QueryParams"

// Params is the struct for params
type Params struct {
	Where	map[string]interface{}
	SortOrder	string
}

// ExtractParams func
func ExtractParams(ctx context.Context) *Params {
	var params *Params
	if p, ok := ctx.Value(QueryParams).(Params); ok {
		params = &p
	}
	return params
}

