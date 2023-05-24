package cyral

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Common keys.
const (
	IDKey           = "id"
	HostKey         = "host"
	PortKey         = "port"
	TypeKey         = "type"
	RepositoryIDKey = "repository_id"
	BindingIDKey    = "binding_id"
	SidecarIDKey    = "sidecar_id"
	ListenerIDKey   = "listener_id"
)

func urlQuery(kv map[string]string) string {
	queryStr := "?"
	for k, v := range kv {
		queryStr += fmt.Sprintf("&%s=%s", k, v)
	}
	return queryStr
}

func elementsMatch(this, other []string) bool {
	if len(this) != len(other) {
		return false
	}
	copyThis := append([]string{}, this...)
	copyOther := append([]string{}, other...)
	sort.Strings(copyThis)
	sort.Strings(copyOther)
	for i, elemThis := range copyThis {
		if elemOther := copyOther[i]; elemThis != elemOther {
			return false
		}
	}
	return true
}

func listToStr(attributes []string) string {
	if len(attributes) == 0 {
		return "[]"
	}
	s := "["
	s += fmt.Sprintf(`"%s"`, attributes[0])
	if len(attributes) > 1 {
		for _, attribute := range attributes[1:] {
			s += fmt.Sprintf(`, "%s"`, attribute)
		}
	}
	s += "]"
	return s
}

func listToStrNoQuotes(l []string) string {
	if len(l) == 0 {
		return "[]"
	}
	s := "["
	s += l[0]
	if len(l) > 1 {
		for _, attribute := range l[1:] {
			s += fmt.Sprintf(`, %s`, attribute)
		}
	}
	s += "]"
	return s
}

func supportedTypesMarkdown(types []string) string {
	var s string
	for _, typ := range types {
		s += fmt.Sprintf("\n  - `%s`", typ)
	}
	return s
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

func validationStringLenAtLeast(min int) schema.SchemaValidateFunc {
	return validation.StringLenBetween(min, math.MaxInt)
}

func typeSetNonEmpty(d *schema.ResourceData, attname string) bool {
	return len(d.Get(attname).(*schema.Set).List()) > 0
}

func getStrList(m map[string]interface{}, attName string) []string {
	var attStrs []string
	for _, valIface := range m[attName].([]interface{}) {
		attStrs = append(attStrs, valIface.(string))
	}
	return attStrs
}

func schemaAllComputed(s map[string]*schema.Schema) map[string]*schema.Schema {
	for k, _ := range s {
		s[k].Optional = false
		s[k].Required = false
		s[k].Computed = true
		s[k].Default = nil
		s[k].MaxItems = 0
		s[k].ExactlyOneOf = nil
		s[k].ValidateFunc = nil
		if s[k].Elem != nil {
			schemaAllComputed(s[k].Elem.(*schema.Resource).Schema)
		}
	}

	return s
}

func boolAsString(v bool) string {
	if v {
		return "true"
	} else {
		return "false"
	}
}
