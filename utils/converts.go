package utils

import (
	"database/sql"
	"golang.org/x/exp/constraints"
	"strconv"
)

// SqlItptS converts NullInt64 to *string
func SqlItptS(i sql.NullInt64) *string {
	var result *string
	result = nil
	if i.Valid {
		tmp := strconv.Itoa(int(i.Int64))
		result = &tmp
	}
	return result
}

// PtStI converts nullables, *string to *int
func PtStI[O constraints.Integer](i *string) (*O, error) {
	var result *O
	if i != nil {
		i, err := strconv.Atoi(*i)
		if err != nil {
			return nil, err
		}
		ui := O(i)
		result = &ui
	}
	return result, nil
}
