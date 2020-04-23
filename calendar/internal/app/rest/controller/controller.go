package controller

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"strconv"
)

//type IService interface {}

// Controller is the parent struct for all controllers
type Controller struct {
}

var matchedParams = []string{
}


func (c Controller) parseUint(ctx *routing.Context, paramName string) (uint, error) {
	paramVal, err := strconv.ParseUint(ctx.Param(paramName), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(paramVal), nil
}

// ExtractQueryFromContext extracts a query from context
func (c Controller) ExtractQueryFromContext(ctx *routing.Context) map[string]interface{} {
	query := make(map[string]interface{}, 1)

	for _, paramName := range matchedParams {

		if paramVal := ctx.Param(paramName); paramVal != "" {
			query[paramName] = paramVal
		}
	}

	return query
}


