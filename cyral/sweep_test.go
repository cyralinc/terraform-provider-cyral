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

// TODO (probably next MAJOR): the sweepers should be put in a dedicated
// package. In this file, we should just have an _import_ of all the resource
// packages, like below, just to trigger the _init_ effect. Each package would
// be responsible to implement their sweepers, and include the
// resource.AddTestSweeper calls in their _init_ function.
//
// import (
//
//	_ "github.com/cyralinc/terraform-provider-cyral/internal/resources/repository"
//	_ "github.com/cyralinc/terraform-provider-cyral/internal/resources/sidecar"
//
// )
//
// Obviously, this requires a complete reorganization of our packaging
// structure. Since we currently export many functions, changing the package
// structure is technically a breaking change, so it might be best to leave this
// for the next MAJOR release. -aholmquist 2022-08-10
func init() {
	sidecarRepositoryName := sidecarResourceName + "_" + repositoryResourceName
	resource.AddTestSweepers(sidecarRepositoryName, &resource.Sweeper{
		Name: sidecarRepositoryName,
		F:    sweepSidecarAndRepository,
	})
	resource.AddTestSweepers(roleResourceName, &resource.Sweeper{
		Name: roleResourceName,
		F:    sweepRole,
	})
	resource.AddTestSweepers(integrationIdPResourceName, &resource.Sweeper{
		Name: integrationIdPResourceName,
		F:    sweepIntegrationIdP,
	})
	resource.AddTestSweepers(policyResourceName, &resource.Sweeper{
		Name: policyResourceName,
		F:    sweepPolicy,
	})
	// TODO: add sweepers for rest of resources -aholmquist 2022-08-10
}

// These must be combined to ensure that sidecars are deleted before
// repositories, otherwise deletion of repositories which have bound ports might
// fail.
func sweepSidecarAndRepository(_ string) error {
	err := sweepSidecar("")
	if err != nil {
		return err
	}
	err = sweepRepository("")
	if err != nil {
		return err
	}
	return nil
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

func sweepSidecar(_ string) error {
	c, err := newClientFromEnv()
	if err != nil {
		return err
	}
	sidecars, err := listSidecars(c)
	if err != nil {
		return err
	}
	for _, sidecar := range sidecars {
		if !hasAccTestPrefix(sidecar.Sidecar.Name) {
			continue
		}
		url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane,
			sidecar.ID)
		_, err := c.DoRequest(url, http.MethodDelete, nil)
		if err != nil {
			return fmt.Errorf("delete request returned error: %w", err)
		}
	}
	return nil
}

func sweepRole(_ string) error {
	c, err := newClientFromEnv()
	if err != nil {
		return err
	}
	resp, err := listRoles(c)
	if err != nil {
		return err
	}
	roles := resp.Groups
	for _, role := range roles {
		if !hasAccTestPrefix(role.Name) {
			continue
		}
		url := fmt.Sprintf("https://%s/v1/users/groups/%s", c.ControlPlane,
			role.ID)
		_, err := c.DoRequest(url, http.MethodDelete, nil)
		if err != nil {
			return fmt.Errorf("delete request returned error: %w", err)
		}
	}
	return nil
}

func sweepPolicy(_ string) error {
	c, err := newClientFromEnv()
	if err != nil {
		return err
	}
	policies, err := listPolicies(c)
	if err != nil {
		return err
	}
	for _, policy := range policies {
		if !hasAccTestPrefix(policy.Meta.Name) {
			continue
		}
		url := fmt.Sprintf("https://%s/v1/policies/%s",
			c.ControlPlane, policy.Meta.ID)
		_, err := c.DoRequest(url, http.MethodDelete, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func sweepIntegrationIdP(_ string) error {
	c, err := newClientFromEnv()
	if err != nil {
		return err
	}
	resp, err := listIdPIntegrations(c)
	if err != nil {
		return fmt.Errorf("failed to get IdP integrations: %w", err)
	}

	integrations := resp.Connections.Connections
	for _, integration := range integrations {
		if !hasAccTestPrefix(integration.DisplayName) {
			continue
		}
		url := fmt.Sprintf("https://%s/v1/integrations/saml/%s",
			c.ControlPlane, integration.Alias)
		_, err := c.DoRequest(url, http.MethodDelete, nil)
		if err != nil {
			return fmt.Errorf("delete request returned error: %w", err)
		}
	}

	return nil
}
