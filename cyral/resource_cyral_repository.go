package cyral

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	r "github.com/cyralinc/crud/repo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var supportedRepoTypes = map[string]bool{
	"galera":     true,
	"mariadb":    true,
	"mysql":      true,
	"bigquery":   true,
	"dremio":     true,
	"postgresql": true,
	"snowflake":  true,
	"sqlserver":  true,
	"cassandra":  true,
	"mongodb":    true,
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
	repoType, repoHost, repoPort, repoName, requireTLS, err := getRepoInfo(config, d)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/repos", config.controlPlane)
	payloadStr := fmt.Sprintf("{\"type\":\"%s\",\"name\":\"%s\",\"hostName\":\"%s\",\"port\":%d,\"repo_tls\":%t}",
		repoType, repoName, repoHost, repoPort, requireTLS)
	payload := strings.NewReader(payloadStr)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("unable to create new repo request, err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable execute new request, err: %v", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to read data from request body, err: %v", err)
	}

	type jsonResponse struct {
		id string `json:"ID" form:"-"`
	}
	idMap := map[string]string{}
	err = json.Unmarshal(body, &idMap)
	id := idMap["ID"]
	if err != nil {
		return fmt.Errorf("unable to get repo ID from json, err: %v", err)
	}

	config.repoID = id
	d.SetId(id)

	return resourceCyralRepositoryRead(d, m)
}

func resourceCyralRepositoryRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	log.Printf("############ >> resourceCyralRepositoryFindById")
	log.Printf(fmt.Sprintf("############ >> d.Id(): %s", d.Id()))
	err := resourceCyralRepositoryFindById(config, d.Id())

	if err == nil {
		//d.SetId("")
	}
	log.Printf("############ >> 6")

	d.Set("type", "mysql")
	d.Set("host", "some_host")
	d.Set("port", 3306)
	d.Set("name", "some_random_repo_name")
	d.Set("require_tls", true)

	log.Printf("############ >> 7")

	return nil
}

func resourceCyralRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceCyralRepositoryRead(d, m)
}

func resourceCyralRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCyralRepositoryFindById(config *Config, id string) error {
	log.Printf("############ >> 1")
	url := fmt.Sprintf("https://%s/repos/%s", config.controlPlane, config.repoID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	log.Printf("############ >> 2")
	if err != nil {
		return fmt.Errorf("unable to create new request, err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer: %s", config.token))
	log.Printf("############ >> 3")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable execute new repo request, err: %v", err)
	}

	log.Printf("############ >> 4")
	defer res.Body.Close()
	if res.Body == http.NoBody {
		return fmt.Errorf("unable to read data from request body, err: %v", err)
	}
	log.Printf("############ >> 5")

	return nil
}

func getRepoInfo(config *Config, d *schema.ResourceData) (string, string, int, string, bool, error) {
	repoType := d.Get("type").(string)
	err := containsRepoType(repoType)
	if err != nil {
		return "", "", 0, "", false, err
	}

	repoHost := d.Get("host").(string)
	repoPort := d.Get("port").(int)
	repoName := d.Get("name").(string)
	requireTLS := d.Get("require_tls").(bool)

	return repoType, repoHost, repoPort, repoName, requireTLS, nil
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
