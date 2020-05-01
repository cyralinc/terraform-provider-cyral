package cyral

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	r "github.com/cyralinc/crud/repo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type resourceCyralRepositoryData struct {
	RepoType   string `json:"type"`
	Name       string `json:"name"`
	Host       string `json:"hostName"`
	Port       int    `json:"port"`
	RequireTLS bool   `json:"repo_tls"`
}

type getRepoResponse struct {
	ID   string                 `json:"id"`
	Repo map[string]interface{} `json:"repo"`
	TLS  bool                   `json:"TLS"`
}

func resourceCyralRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceCyralRepositoryCreate,
		Read:   resourceCyralRepositoryRead,
		Update: resourceCyralRepositoryUpdate,
		Delete: resourceCyralRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"require_tls": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceCyralRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	repoData, err := getRepoData(config, d)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/repos", config.controlPlane)
	payloadBytes, err := json.Marshal(repoData)
	if err != nil {
		return fmt.Errorf("failed to encode newRepo payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return fmt.Errorf("unable to create newRepo request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", config.tokenType, config.token))
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode == http.StatusConflict {
		return fmt.Errorf("repository name already exists in control plane")
	}

	if err != nil {
		return fmt.Errorf("unable execute newRepo request; err: %v", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to read data from request body; err: %v", err)
	}

	jsonMap := map[string]string{}
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		return fmt.Errorf("unable to get repo ID from json; err: %v", err)
	}

	id := jsonMap["ID"]
	config.repoID = id
	d.SetId(id)

	return resourceCyralRepositoryRead(d, m)
}

func resourceCyralRepositoryRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	repoData, err := resourceCyralRepositoryFindByID(config, d.Id())

	if err != nil {
		return err
	}

	if repoData == nil {
		fmt.Errorf("repo not found; id: %s", d.Id())
	}

	d.Set("type", repoData.RepoType)
	d.Set("host", repoData.Host)
	d.Set("port", repoData.Port)
	d.Set("name", repoData.Name)
	d.Set("require_tls", repoData.RequireTLS)

	return nil
}

func resourceCyralRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	// TODO Warn users if they modify `require_tls` parameter in .tf, as it
	// is not possible to change it once the repo is created.
	config := m.(*Config)
	repoData, err := getRepoData(config, d)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/repos/%s", config.controlPlane, d.Id())
	payloadBytes, err := json.Marshal(repoData)
	if err != nil {
		return fmt.Errorf("failed to encode updateRepo payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return fmt.Errorf("unable to create updateRepo request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", config.tokenType, config.token))
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("repository not found; id: %s", d.Id())
	}
	if res.StatusCode == http.StatusConflict {
		return fmt.Errorf("repository name already exists in control plane")
	}

	if err != nil {
		return fmt.Errorf("unable execute updateRepo request; err: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response from updateRepo; status code: %d", res.StatusCode)
	}

	return resourceCyralRepositoryRead(d, m)
}

func resourceCyralRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	url := fmt.Sprintf("https://%s/repos/%s", config.controlPlane, d.Id())

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create deleteRepo request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", config.tokenType, config.token))
	_, err = http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("unable execute deleteRepo request; err: %v", err)
	}

	return nil
}

func resourceCyralRepositoryFindByID(config *Config, id string) (*resourceCyralRepositoryData, error) {
	url := fmt.Sprintf("https://%s/repos/%s", config.controlPlane, id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create new request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", config.tokenType, config.token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable execute new repo request; err: %v", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read data from request body; err: %v", err)
	}

	// Not an error, nor any data was found
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	repoRespJSON := getRepoResponse{
		Repo: map[string]interface{}{},
	}
	err = json.Unmarshal(body, &repoRespJSON)
	if err != nil {
		return nil, fmt.Errorf("unable to get repo json by id, err: %v", err)
	}
	if repoRespJSON.ID != id {
		return nil, fmt.Errorf("unable to get repo json by id, err: %v", err)
	}

	return &resourceCyralRepositoryData{
		RepoType:   repoRespJSON.Repo["type"].(string),
		Host:       repoRespJSON.Repo["repoHost"].(string),
		Name:       repoRespJSON.Repo["name"].(string),
		Port:       int(repoRespJSON.Repo["repoPort"].(float64)),
		RequireTLS: repoRespJSON.TLS,
	}, nil
}

func validateGetRepoRequest(repoRespJSON getRepoResponse) error {
	if _, ok := repoRespJSON.Repo["name"]; !ok {
		return fmt.Errorf("missing 'name' field in GET repo response %v", repoRespJSON)
	}
	if _, ok := repoRespJSON.Repo["type"]; !ok {
		return fmt.Errorf("missing 'type' field in GET repo response %v", repoRespJSON)
	}
	if _, ok := repoRespJSON.Repo["repoHost"]; !ok {
		return fmt.Errorf("missing 'repoHost' field in GET repo response %v", repoRespJSON)
	}
	if _, ok := repoRespJSON.Repo["repoPort"]; !ok {
		return fmt.Errorf("missing 'repoPort' field in GET repo response %v", repoRespJSON)
	}
	return nil
}

func getRepoData(config *Config, d *schema.ResourceData) (resourceCyralRepositoryData, error) {
	repoType := d.Get("type").(string)
	err := containsRepoType(repoType)
	if err != nil {
		return resourceCyralRepositoryData{}, err
	}

	repoData := resourceCyralRepositoryData{
		RepoType:   repoType,
		Host:       d.Get("host").(string),
		Name:       d.Get("name").(string),
		Port:       d.Get("port").(int),
		RequireTLS: d.Get("require_tls").(bool),
	}

	return repoData, nil
}

func containsRepoType(repoType string) error {
	repoTypes := r.GetRepoTypes()
	mapRepoTypes := make(map[string]bool)
	for _, r := range repoTypes {
		mapRepoTypes[r] = true
	}
	if mapRepoTypes[repoType] == false {
		return fmt.Errorf("repo type must be one of %v", repoTypes)
	}
	return nil
}
