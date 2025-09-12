package helper

import (
	"database/sql"
	"fmt"
	"strconv"
)

func ParseIDsNumeric(strIDs []string) ([]int64, error) {
	out := make([]int64, 0, len(strIDs))
	for _, s := range strIDs {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid numeric id: %s", s)
		}
		out = append(out, id)
	}
	return out, nil
}

func NullString(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}
