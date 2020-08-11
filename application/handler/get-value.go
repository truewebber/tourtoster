package handler

import (
	"html"
	"net/url"
	"strconv"
	"strings"
)

func stringValue(key string, v url.Values) string {
	return strings.TrimSpace(v.Get(key))
}

func stringEscapedValue(key string, v url.Values) string {
	strVal := stringValue(key, v)

	return html.EscapeString(strVal)
}

func intValue(key string, v url.Values) (int, error) {
	strVal := stringValue(key, v)

	return strconv.Atoi(strVal)
}
