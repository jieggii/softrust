package domain

import "strings"

func NormalizeQuery(query string) string {
	query = strings.ToLower(query)
	query = strings.Trim(query, " \t")
	return query
}
