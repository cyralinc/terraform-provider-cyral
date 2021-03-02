package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateRepoResponse struct {
	ID string `json:"ID"`
}

type GetRepoByIDResponse struct {
	Repo RepoData `json:"repo"`
}

type RepoData struct {
	ID       string `json:"id"`
	RepoType string `json:"type"`
	Name     string `json:"name"`
	Host     string `json:"repoHost"`
	Port     int    `json:"repoPort"`
}

func resourceCyralRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCyralRepositoryCreate,
		ReadContext:   resourceCyralRepositoryRead,
		UpdateContext: resourceCyralRepositoryUpdate,
		DeleteContext: resourceCyralRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCyralRepositoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceCyralRepositoryCreate")
	c := m.(*client.Client)

	dataFromResource, err := getRepoDataFromResource(c, d)
	if err != nil {
		return createError("Unable to create repository", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Resource info: %#v", dataFromResource)

	url := fmt.Sprintf("https://%s/v1/repos", c.ControlPlane)
	log.Printf("[DEBUG] POST URL: %s", url)
	payloadBytes, err := json.Marshal(dataFromResource)
	if err != nil {
		return createError("Unable to create repository",
			fmt.Sprintf("failed to encode 'create repo' payload: %v", err))
	}

	log.Printf("[DEBUG] POST payload: %s", string(payloadBytes))
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return createError("Unable to create repository",
			fmt.Sprintf("unable to create 'create repo' request; err: %v", err))
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))

	log.Printf("[DEBUG] Executing POST")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return createError("Unable to create repository",
			fmt.Sprintf("unable to execute 'create repo' request. Check the control plane address; err: %v", err))
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusConflict {
		return createError("Unable to create repository",
			fmt.Sprintf("repository name already exists in control plane"))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return createError("Unable to create repository",
			fmt.Sprintf("unable to read data from request body; err: %v", err))
	}
	log.Printf("[DEBUG] Response body: %s", string(body))

	if res.StatusCode != http.StatusCreated {
		return createError("Unable to create repository",
			fmt.Sprintf("unexpected response from 'create repo' request; status code: %d; body: %q",
				res.StatusCode, body))
	}

	unmarshalledBody := CreateRepoResponse{}
	if err := json.Unmarshal(body, &unmarshalledBody); err != nil {
		return createError("Unable to create repository",
			fmt.Sprintf("unable to unmarshall json; err: %v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", unmarshalledBody)

	d.SetId(unmarshalledBody.ID)

	return resourceCyralRepositoryRead(ctx, d, m)
}

func resourceCyralRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceCyralRepositoryRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, d.Id())

	log.Printf("[DEBUG] GET URL: %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return createError("Unable to read repository",
			fmt.Sprintf("unable to create new request; err: %v", err))
	}

	log.Printf("[DEBUG] Executing GET")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))
	log.Printf("[DEBUG] GET request: %#v", req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return createError("Unable to read repository",
			fmt.Sprintf("unable to execute request at resourceCyralRepositoryRead."+
				" Check the control plane address; err: %v", err))
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return createError("Unable to read repository",
			fmt.Sprintf("unable to read data from request body at resourceCyralRepositoryRead; err: %v", err))
	}
	log.Printf("[DEBUG] Response body: %s", string(body))

	// Not an error, nor any data was found
	if res.StatusCode == http.StatusNotFound {
		return createError("Unable to read repository",
			fmt.Sprintf("repo not found; id: %s", d.Id()))
	}

	if res.StatusCode != http.StatusOK {
		return createError("Unable to read repository",
			fmt.Sprintf("unexpected response from resourceCyralRepositoryRead; status code: %d; body: %q", res.StatusCode, res.Body))
	}

	unmarshalledBody := GetRepoByIDResponse{}
	if err := json.Unmarshal(body, &unmarshalledBody); err != nil {
		return createError("Unable to read repository",
			fmt.Sprintf("unable to get repo json by id, err: %v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", unmarshalledBody)

	d.Set("type", unmarshalledBody.Repo.RepoType)
	d.Set("host", unmarshalledBody.Repo.Host)
	d.Set("port", unmarshalledBody.Repo.Port)
	d.Set("name", unmarshalledBody.Repo.Name)

	log.Printf("[DEBUG] End resourceCyralRepositoryRead")

	return diag.Diagnostics{}
}

func resourceCyralRepositoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceCyralRepositoryUpdate")
	c := m.(*client.Client)

	dataFromResource, err := getRepoDataFromResource(c, d)
	if err != nil {
		return createError("Unable to update repository", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, d.Id())
	payloadBytes, err := json.Marshal(dataFromResource)
	if err != nil {
		return createError("Unable to update repository",
			fmt.Sprintf("failed to encode 'update repo' payload: %v", err))
	}

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return createError("Unable to update repository",
			fmt.Sprintf("unable to create 'update repo' request; err: %v", err))
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return createError("Unable to update repository",
			fmt.Sprintf("unable to execute 'update repo' request. Check the control plane address; err: %v", err))
	}

	if res.StatusCode == http.StatusNotFound {
		return createError("Unable to update repository",
			fmt.Sprintf("repository not found; id: %s", d.Id()))
	} else if res.StatusCode == http.StatusConflict {
		return createError("Unable to update repository",
			fmt.Sprintf("repository name already exists in control plane"))
	} else if res.StatusCode != http.StatusOK {
		return createError("Unable to update repository",
			fmt.Sprintf("unexpected response from 'update repo'; status code: %d; body: %q",
				res.StatusCode, res.Body))
	}

	return resourceCyralRepositoryRead(ctx, d, m)
}

func resourceCyralRepositoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceCyralRepositoryDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, d.Id())
	log.Printf("[DEBUG] DELETE URL: %s", url)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return createError("Unable to delete repository",
			fmt.Sprintf("unable to create 'delete repo' request; err: %v", err))
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return createError("Unable to delete repository",
			fmt.Sprintf("unable execute 'delete repo' request. Check the control plane address; err: %v", err))
	}

	if res.StatusCode != http.StatusOK {
		return createError("Unable to delete repository",
			fmt.Sprintf("unexpected response from 'delete repo' request; status code: %d; body: %q", res.StatusCode, res.Body))
	}

	log.Printf("[DEBUG] End resourceCyralRepositoryDelete")

	return diag.Diagnostics{}
}

func getRepoDataFromResource(c *client.Client, d *schema.ResourceData) (RepoData, error) {
	repoType := d.Get("type").(string)
	err := containsRepoType(repoType)
	if err != nil {
		return RepoData{}, err
	}

	return RepoData{
		ID:       d.Id(),
		RepoType: repoType,
		Host:     d.Get("host").(string),
		Name:     d.Get("name").(string),
		Port:     d.Get("port").(int),
	}, nil
}

func containsRepoType(repoType string) error {
	// This code was copied here to remove dependency of CRUD,
	// but we should move the CRUD code to CRUD-API (or somewhere
	// else) in the future.
	repoTypes := map[string]bool{
		"bigquery":   true,
		"cassandra":  true,
		"dremio":     true,
		"galera":     true,
		"mariadb":    true,
		"mongodb":    true,
		"mysql":      true,
		"oracle":     true,
		"postgresql": true,
		"snowflake":  true,
		"s3":         true,
		"sqlserver":  true,
	}
	if repoTypes[repoType] == false {
		return fmt.Errorf("repo type must be one of %v", repoTypes)
	}
	return nil
}
