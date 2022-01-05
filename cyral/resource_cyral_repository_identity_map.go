package cyral

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rickb777/date/period"
)

type AccessDuration struct {
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

func (data AccessDuration) formatTime() string {
	return period.New(0, 0,
		data.Days,
		data.Hours,
		data.Minutes,
		data.Seconds,
	).String()
}

func (data *AccessDuration) getTimeFromString(payload string) error {
	accessDurationPeriod, err := period.Parse(payload)
	if err != nil {
		return err
	}

	accessDurationNormalized := accessDurationPeriod.Normalise(false)

	data.Days = accessDurationNormalized.Days()
	data.Hours = accessDurationNormalized.Hours()
	data.Minutes = accessDurationNormalized.Minutes()
	data.Seconds = accessDurationNormalized.Seconds()

	return nil
}

type RepositoryIdentityMapAPIBody struct {
	AccessDuration string `json:"accessDuration,omitempty"`
}

type RepositoryIdentityMapResource struct {
	RepositoryId          string          `json:"-"`
	IdentityType          string          `json:"-"`
	IdentityName          string          `json:"-"`
	RepositoryAccountUUID string          `json:"-"`
	AccessDuration        *AccessDuration `json:"accessDuration,omitempty"`
}

func (data RepositoryIdentityMapResource) WriteToSchema(d *schema.ResourceData) {
	d.Set("repository_id", data.RepositoryId)

	if err := data.isIdentityTypeValid(); err != nil {
		panic(err)
	}
	d.Set("identity_type", data.IdentityType)
	d.Set("repository_local_account_id", data.RepositoryAccountUUID)
	d.Set("identity_name", data.IdentityName)
	if data.AccessDuration != nil {
		d.Set("access_duration", []interface{}{
			map[string]interface{}{
				"days":    data.AccessDuration.Days,
				"hours":   data.AccessDuration.Hours,
				"minutes": data.AccessDuration.Minutes,
				"seconds": data.AccessDuration.Seconds,
			},
		})
	}
}

func (data *RepositoryIdentityMapResource) ReadFromSchema(d *schema.ResourceData) {
	data.RepositoryId = d.Get("repository_id").(string)
	data.IdentityType = d.Get("identity_type").(string)
	if err := data.isIdentityTypeValid(); err != nil {
		panic(err)
	}
	data.RepositoryAccountUUID = d.Get("repository_local_account_id").(string)
	data.IdentityName = d.Get("identity_name").(string)

	if _, hasAcessDuration := d.GetOk("access_duration"); hasAcessDuration {
		data.AccessDuration = &AccessDuration{}
		acess := d.Get("access_duration").(*schema.Set)

		for _, id := range acess.List() {
			idMap := id.(map[string]interface{})

			data.AccessDuration.Days = idMap["days"].(int)
			data.AccessDuration.Hours = idMap["hours"].(int)
			data.AccessDuration.Minutes = idMap["minutes"].(int)
			data.AccessDuration.Seconds = idMap["seconds"].(int)
		}
	}
}

func (resource *RepositoryIdentityMapResource) UnmarshalJSON(data []byte) error {
	var response RepositoryIdentityMapAPIBody
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	if response.AccessDuration != "" {
		resource.AccessDuration = &AccessDuration{}
		if err := resource.AccessDuration.getTimeFromString(response.AccessDuration); err != nil {
			return err
		}
	}

	return nil
}

func (resource *RepositoryIdentityMapResource) MarshalJSON() ([]byte, error) {
	payload := RepositoryIdentityMapAPIBody{}
	if resource.AccessDuration != nil {
		payload.AccessDuration = resource.AccessDuration.formatTime()
	}

	return json.Marshal(payload)
}

func (data RepositoryIdentityMapResource) isIdentityTypeValid() error {
	if !(data.IdentityType == "user" || data.IdentityType == "group") {
		return errors.New("invalid identity type")
	}
	return nil
}

type RepositoryIdentityMapAPIResponse struct {
	AccessDuration *AccessDuration `json:"accessDuration,omitempty"`
}

func (data RepositoryIdentityMapAPIResponse) WriteToSchema(d *schema.ResourceData) {
	d.SetId(fmt.Sprintf("%s-%s", d.Get("repository_id").(string),
		d.Get("repository_local_account_id").(string)))
	if data.AccessDuration != nil {
		d.Set("access_duration", []interface{}{
			map[string]interface{}{
				"days":    data.AccessDuration.Days,
				"hours":   data.AccessDuration.Hours,
				"minutes": data.AccessDuration.Minutes,
				"seconds": data.AccessDuration.Seconds,
			},
		})
	}
}

func (data *RepositoryIdentityMapAPIResponse) ReadFromSchema(d *schema.ResourceData) {
	if _, hasAcessDuration := d.GetOk("access_duration"); hasAcessDuration {
		data.AccessDuration = &AccessDuration{}
		acess := d.Get("access_duration").(*schema.Set)

		for _, id := range acess.List() {
			idMap := id.(map[string]interface{})

			data.AccessDuration.Days = idMap["days"].(int)
			data.AccessDuration.Hours = idMap["hours"].(int)
			data.AccessDuration.Minutes = idMap["minutes"].(int)
			data.AccessDuration.Seconds = idMap["seconds"].(int)
		}
	}
}

func (resource *RepositoryIdentityMapAPIResponse) UnmarshalJSON(data []byte) error {
	var response RepositoryIdentityMapAPIBody
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	if response.AccessDuration != "" && response.AccessDuration != "P0D" {
		resource.AccessDuration = &AccessDuration{}
		if err := resource.AccessDuration.getTimeFromString(response.AccessDuration); err != nil {
			return err
		}
	} else {
		resource.AccessDuration = nil
	}

	return nil
}

var ReadRepositoryIdentityMapConfig = ResourceOperationConfig{
	Name:       "RepositoryIdentityMapResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/repos/%s/identityMaps/%s/%s/%s",
			c.ControlPlane,
			d.Get("repository_id").(string),
			d.Get("identity_type").(string),
			d.Get("identity_name").(string),
			d.Get("repository_local_account_id").(string))
	},
	ResponseData: &RepositoryIdentityMapAPIResponse{},
}

func resourceRepositoryIdentityMap(deprecate bool) *schema.Resource {
	deprecationMessage := ""

	if deprecate {
		deprecationMessage = "Use `cyral_repository_identity_map` instead."
	}

	return &schema.Resource{
		DeprecationMessage: deprecationMessage,
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "RepositoryIdentityMapResourceCreate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/repos/%s/identityMaps/%s/%s/%s",
						c.ControlPlane,
						d.Get("repository_id").(string),
						d.Get("identity_type").(string),
						d.Get("identity_name").(string),
						d.Get("repository_local_account_id").(string))
				},
				ResourceData: &RepositoryIdentityMapResource{},
				ResponseData: &RepositoryIdentityMapAPIResponse{},
			}, ReadRepositoryIdentityMapConfig,
		),
		ReadContext: ReadResource(ReadRepositoryIdentityMapConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "RepositoryIdentityMapResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/repos/%s/identityMaps/%s/%s/%s",
						c.ControlPlane,
						d.Get("repository_id").(string),
						d.Get("identity_type").(string),
						d.Get("identity_name").(string),
						d.Get("repository_local_account_id").(string))
				},
				ResourceData: &RepositoryIdentityMapResource{},
			}, ReadRepositoryIdentityMapConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "RepositoryIdentityMapResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/repos/%s/identityMaps/%s/%s/%s",
						c.ControlPlane,
						d.Get("repository_id").(string),
						d.Get("identity_type").(string),
						d.Get("identity_name").(string),
						d.Get("repository_local_account_id").(string))
				},
			},
		),
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository_local_account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identity_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identity_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_duration": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"hours": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"minutes": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
