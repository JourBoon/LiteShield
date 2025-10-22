package cache

import (
	"net/http"
	"encoding/hex"
	"crypto/sha1"
	"sort"
	"strings"
)

func GenerateCacheKey(r *http.Request) string {
	method := r.Method

	u := r.URL
	query := u.Query()

	var keys []string
	for k := range query {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var queryParts []string
	for _, k := range keys {
		values := query[k]
		sort.Strings(values)
		for _, v := range values {
			queryParts = append(queryParts, k+"="+v)
		}
	}

	queryString := strings.Join(queryParts, "&")
	fullPath := u.Path
	if queryString != "" {
		fullPath += "?" + queryString
	}

	key := method + ":" + fullPath

	h := sha1.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}