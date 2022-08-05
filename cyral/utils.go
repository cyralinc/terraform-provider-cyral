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

func marshalComposedID(ids []string, sep string) string {
	return strings.Join(ids, sep)
}

func unmarshalComposedID(id, sep string, numFields int) ([]string, error) {
	ids := strings.Split(id, sep)
	if len(ids) < numFields {
		return nil, fmt.Errorf("unexpected ID syntax. Correct ID " +
			fmt.Sprintf("syntax uses separator %q and contains %d "+
				"fields", sep, numFields))
	}
	return ids, nil
}
