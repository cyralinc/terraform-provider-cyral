package cyral

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers("cyral_"+repositoryResourceName, &resource.Sweeper{
		Name: repositoryResourceName,
		F:    sweepRepository,
	})
	// TODO: add sweepers for rest of resources. This should take some
	// time to implement... -aholmquist 2022-08-10
}

func sweepRepository(_ string) error {
	c, err := newClientFromEnv()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://%s/v1/repos?name=^%s", c.ControlPlane,
		tprovACCPrefix)
	reposBytes, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return fmt.Errorf("get request returned error: %w", err)
	}
	repos := GetReposResponse{}
	if err := json.Unmarshal(reposBytes, &repos); err != nil {
		return fmt.Errorf("error unmarshaling resp: %w", err)
	}
	for _, repo := range repos.Repos {
		url = fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, repo.ID)
		_, err := c.DoRequest(url, http.MethodDelete, nil)
		if err != nil {
			return fmt.Errorf("delete request returned error: %w", err)
		}
	}
	return nil
}
