package tools

import (
	"fmt"
	"strconv"
)

const ErrConvert = -10

type QueryParams struct {
	Limit      int
	Offset     int
	Query      string
	OrderField string
	OrderBy    int
}

func NewQueryParams(limit, offset, query, orderField, orderBy string) (*QueryParams, error) {
	limitInt := ConvertStrToInt(limit)
	if limitInt < 0 {
		return nil, fmt.Errorf("invalid limit param")
	}
	offsetInt := ConvertStrToInt(offset)
	if offsetInt < 0 {
		return nil, fmt.Errorf("invalid offset param")
	}
	orderByInt := ConvertStrToInt(orderBy)

	switch orderByInt {
	case 1, 0, -1:
	default:
		orderByInt = ErrConvert
	}

	switch orderField {
	case "Name", "name", "":
		orderField = "name"
	case "Id", "id":
		orderField = "id"
	case "Age", "age":
		orderField = "age"
	default:
		return nil, fmt.Errorf("invalid order_field param")
	}

	return &QueryParams{
		Limit:      limitInt,
		Offset:     offsetInt,
		Query:      query,
		OrderField: orderField,
		OrderBy:    orderByInt,
	}, nil
}

func ConvertStrToInt(numStr string) int {
	numInt, err := strconv.Atoi(numStr)
	if err != nil {
		numInt = ErrConvert
	}
	return numInt
}
