package client

import "fmt"

func MakeStrForSD(sensitiveData SensitiveData) string {
	var sd string

	for key, value := range sensitiveData {
		s1 := fmt.Sprintf("map[%s]", key)
		sd = sd + s1
		for _, r := range value {
			s2 := fmt.Sprintf("repo: %s, attributes: [ ", r.Name)
			sd += s2
			for _, attr := range r.Attributes {
				s3 := fmt.Sprintf("%s ", attr)
				sd += s3
			}
			sd += "]"
		}
	}

	return sd
}

// type Datamap struct {
// 	// ID     string         `json:"id"`
// 	Labels []DatamapLabel `json:"labels,omitempty"`
// }

// type DatamapLabel struct {
// 	Name string      `json:"name,omitempty"`
// 	Info []LabelInfo `json:"info,omitempty"`
// }

// type RepoAttrs struct {
// 	Repo       string   `json:"repo,omitempty"`
// 	Attributes []string `json:"attributes,omitempty"`
// }

type SensitiveData map[string][]*RepoAttrs

type DataMap struct {
	SensitiveData SensitiveData `json:"sensitiveData" yaml:"sensitiveData"`
}

type RepoAttrs struct {
	Name       string   `json:"repo" yaml:"repo"`
	Attributes []string `json:"attributes" yaml:"attributes"`
}
