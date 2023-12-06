package sweep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/deprecated"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/policy"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/role"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
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
	sidecarRepositoryName := utils.SidecarResourceName + "_" + utils.RepositoryResourceName
	resource.AddTestSweepers(sidecarRepositoryName, &resource.Sweeper{
		Name: sidecarRepositoryName,
		F:    sweepSidecarAndRepository,
	})
	resource.AddTestSweepers(utils.RoleResourceName, &resource.Sweeper{
		Name: utils.RoleResourceName,
		F:    sweepRole,
	})
	resource.AddTestSweepers(utils.IntegrationIdPResourceName, &resource.Sweeper{
		Name: utils.IntegrationIdPResourceName,
		F:    sweepIntegrationIdP,
	})
	resource.AddTestSweepers(utils.PolicyResourceName, &resource.Sweeper{
		Name: utils.PolicyResourceName,
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
	c, err := client.FromEnv()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://%s/v1/repos?name=^%s", c.ControlPlane,
		utils.TFProvACCPrefix)
	reposBytes, err := c.DoRequest(context.Background(), url, http.MethodGet, nil)
	if err != nil {
		return fmt.Errorf("get request returned error: %w", err)
	}
	repos := repository.GetReposResponse{}
	if err := json.Unmarshal(reposBytes, &repos); err != nil {
		return fmt.Errorf("error unmarshaling resp: %w", err)
	}
	for _, repo := range repos.Repos {
		url = fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, repo.ID)
		_, err := c.DoRequest(context.Background(), url, http.MethodDelete, nil)
		if err != nil {
			return fmt.Errorf("delete request returned error: %w", err)
		}
	}
	return nil
}

func sweepSidecar(_ string) error {
	c, err := client.FromEnv()
	if err != nil {
		return err
	}
	sidecars, err := sidecar.ListSidecars(c)
	if err != nil {
		return err
	}
	for _, sidecar := range sidecars {
		if !utils.HasAccTestPrefix(sidecar.Sidecar.Name) {
			continue
		}
		url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane,
			sidecar.ID)
		_, err := c.DoRequest(context.Background(), url, http.MethodDelete, nil)
		if err != nil {
			return fmt.Errorf("delete request returned error: %w", err)
		}
	}
	return nil
}

func sweepRole(_ string) error {
	c, err := client.FromEnv()
	if err != nil {
		return err
	}
	resp, err := role.ListRoles(c)
	if err != nil {
		return err
	}
	roles := resp.Groups
	for _, role := range roles {
		if !utils.HasAccTestPrefix(role.Name) {
			continue
		}
		url := fmt.Sprintf("https://%s/v1/users/groups/%s", c.ControlPlane,
			role.ID)
		_, err := c.DoRequest(context.Background(), url, http.MethodDelete, nil)
		if err != nil {
			return fmt.Errorf("delete request returned error: %w", err)
		}
	}
	return nil
}

func sweepPolicy(_ string) error {
	c, err := client.FromEnv()
	if err != nil {
		return err
	}
	policies, err := policy.ListPolicies(c)
	if err != nil {
		return err
	}
	for _, policy := range policies {
		if !utils.HasAccTestPrefix(policy.Meta.Name) {
			continue
		}
		url := fmt.Sprintf("https://%s/v1/policies/%s",
			c.ControlPlane, policy.Meta.ID)
		_, err := c.DoRequest(context.Background(), url, http.MethodDelete, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func sweepIntegrationIdP(_ string) error {
	c, err := client.FromEnv()
	if err != nil {
		return err
	}
	resp, err := deprecated.ListIdPIntegrations(c)
	if err != nil {
		return fmt.Errorf("failed to get IdP integrations: %w", err)
	}

	integrations := resp.Connections.Connections
	for _, integration := range integrations {
		if !utils.HasAccTestPrefix(integration.DisplayName) {
			continue
		}
		url := fmt.Sprintf("https://%s/v1/integrations/saml/%s",
			c.ControlPlane, integration.Alias)
		_, err := c.DoRequest(context.Background(), url, http.MethodDelete, nil)
		if err != nil {
			return fmt.Errorf("delete request returned error: %w", err)
		}
	}

	return nil
}
