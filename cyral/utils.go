package cyral

import (
	"fmt"
	"strings"
)

func urlQuery(kv map[string]string) string {
	queryStr := "?"
	for k, v := range kv {
		queryStr += fmt.Sprintf("&%s=%s", k, v)
	}
	return queryStr
}

func marshalComposedID(id1, id2 string) string {
	return fmt.Sprintf("%s-%s", id1, id2)
}

func unmarshalComposedID(id, sep string) (string, string, error) {
	ids := strings.Split(id, sep)
	if len(ids) < 2 {
		return "", "", fmt.Errorf("unexpected ID syntax. Correct ID " +
			"syntax is {id1}-{id1}")
	}
	return ids[0], ids[1], nil
}
