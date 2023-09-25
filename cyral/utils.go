package cyral

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Common keys.
const (
	IDKey           = "id"
	NameKey         = "name"
	DescriptionKey  = "description"
	HostKey         = "host"
	PortKey         = "port"
	TypeKey         = "type"
	RepositoryIDKey = "repository_id"
	BindingIDKey    = "binding_id"
	SidecarIDKey    = "sidecar_id"
	ListenerIDKey   = "listener_id"
)

func convertToInterfaceList[T any](list []T) []any {
	if list == nil {
		return nil
	}
	interfaceList := make([]any, len(list))
	for index, item := range list {
		interfaceList[index] = item
	}
	return interfaceList
}

func convertFromInterfaceList[T any](interfaceList []any) []T {
	if interfaceList == nil {
		return nil
	}
	list := make([]T, len(interfaceList))
	for index, item := range interfaceList {
		list[index] = item.(T)
	}
	return list
}

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

func validationDurationString(value interface{}, key string) (warnings []string, errors []error) {
	duration, ok := value.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", key))
		return warnings, errors
	}
	if !strings.HasSuffix(duration, "s") {
		errors = append(errors, fmt.Errorf(
			"expected %s to end with a 's' suffix. For example: `300s`, `60s`, `10.50s` etc. Got `%v`",
			key, duration,
		))
	}
	return warnings, errors
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

func convertSchemaFieldsToComputed(s map[string]*schema.Schema) map[string]*schema.Schema {
	for k, _ := range s {
		s[k] = &schema.Schema{
			Description: s[k].Description,
			Type:        s[k].Type,
			Computed:    true,
			Elem:        s[k].Elem,
		}
		if elem, ok := s[k].Elem.(*schema.Resource); ok && elem != nil {
			convertSchemaFieldsToComputed(elem.Schema)
		}
	}

	return s
}

// setKeysAsNewComputedIfPlanHasChanges is intended to be used in resource CustomizeDiff functions to set
// computed fields that are expected to change as "new computed" (known after apply) so that terraform can
// detect changes in those fields and update them in the resource state correctly in the same plan operation.
// Otherwise, if this function is not called, terraform will not detect a change in those computed fields during
// the initial update operation and the changes will only be detected in the subsequent terraform plan.
// For reference:
// - https://github.com/hashicorp/terraform/issues/15857
func setKeysAsNewComputedIfPlanHasChanges(resourceDiff *schema.ResourceDiff, keys []string) {
	changedKeys := resourceDiff.GetChangedKeysPrefix("")
	log.Printf("[DEBUG] changedKeys: %+v", changedKeys)
	hasChanges := len(changedKeys) > 0
	log.Printf("[DEBUG] hasChanges: %t", hasChanges)
	if hasChanges {
		for _, key := range keys {
			resourceDiff.SetNewComputed(key)
		}
	}
}
