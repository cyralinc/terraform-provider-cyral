package datamap

import (
	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DataMapRequest struct {
	DataMap `json:"dataMap,omitempty"`
}

// This is called 'DataMap' and not 'Datamap', because although we consider
// 'datamap' to be a single word in the resource name 'cyral_repository_datamap'
// for ease of writing, 'data map' is actually two words in English.
type DataMap struct {
	Labels map[string]*DataMapMapping `json:"labels,omitempty"`
}

func (dm *DataMap) WriteToSchema(d *schema.ResourceData) error {
	var mappings []interface{}
	for label, mapping := range dm.Labels {
		mappingContents := make(map[string]interface{})

		var attributes []string
		if mapping != nil {
			attributes = mapping.Attributes
		}

		mappingContents["label"] = label
		mappingContents["attributes"] = attributes

		mappings = append(mappings, mappingContents)
	}

	return d.Set("mapping", mappings)
}

func (dm *DataMap) equal(other DataMap) bool {
	for label, thisMapping := range dm.Labels {
		if otherMapping, ok := other.Labels[label]; ok {
			if !utils.ElementsMatch(thisMapping.Attributes, otherMapping.Attributes) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

type DataMapMapping struct {
	Attributes []string `json:"attributes,omitempty"`
}
