package cyral

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/senseyeio/duration"
)

type AccessDuration struct {
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

func (data AccessDuration) formatTime() string {
	return duration.Duration{
		D:  data.Days,
		TH: data.Hours,
		TM: data.Minutes,
		TS: data.Seconds,
	}.String()
}

func (data *AccessDuration) getTimeFromString(payload string) error {
	dur, err := duration.ParseISO8601(payload)
	if err != nil {
		return err
	}
	data.Days = dur.D
	data.Hours = dur.TH
	data.Minutes = dur.TM
	data.Seconds = dur.TS

	return nil
}

type IdentityMapAPIBody struct {
	AccessDuration string `json:"accessDuration,omitempty"`
}

type IdentityMapResource struct {
	RepositoryId          string          `json:"-"`
	IdentityType          string          `json:"-"`
	IdentityName          string          `json:"-"`
	RepositoryAccountUUID string          `json:"-"`
	AccessDuration        *AccessDuration `json:"accessDuration,omitempty"`
}

func (data IdentityMapResource) WriteToSchema(d *schema.ResourceData) {
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

func (data *IdentityMapResource) ReadFromSchema(d *schema.ResourceData) {
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

func (resource *IdentityMapResource) UnmarshalJSON(data []byte) error {
	var response IdentityMapAPIBody
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

func (resource *IdentityMapResource) MarshalJSON() ([]byte, error) {
	payload := IdentityMapAPIBody{}
	if resource.AccessDuration != nil {
		payload.AccessDuration = resource.AccessDuration.formatTime()
	}

	return json.Marshal(payload)
}

func (data IdentityMapResource) isIdentityTypeValid() error {
	if !(data.IdentityType == "user" || data.IdentityType == "group") {
		return errors.New("invalid identity type")
	}
	return nil
}

type IdentityMapAPIResponse struct {
	AccessDuration *AccessDuration `json:"accessDuration,omitempty"`
}

func (data IdentityMapAPIResponse) WriteToSchema(d *schema.ResourceData) {
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

func (data *IdentityMapAPIResponse) ReadFromSchema(d *schema.ResourceData) {
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

func (resource *IdentityMapAPIResponse) UnmarshalJSON(data []byte) error {
	var response IdentityMapAPIBody
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

var ReadIdentityMapConfig = ResourceOperationConfig{
	Name:       "IdentityMapResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/repos/%s/identityMaps/%s/%s/%s",
			c.ControlPlane,
			d.Get("repository_id").(string),
			d.Get("identity_type").(string),
			d.Get("identity_name").(string),
			d.Get("repository_local_account_id").(string))
	},
	ResponseData: &IdentityMapAPIResponse{},
}

func resourceIdentityMap() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "IdentityMapResourceCreate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/repos/%s/identityMaps/%s/%s/%s",
						c.ControlPlane,
						d.Get("repository_id").(string),
						d.Get("identity_type").(string),
						d.Get("identity_name").(string),
						d.Get("repository_local_account_id").(string))
				},
				ResourceData: &IdentityMapResource{},
				ResponseData: &IdentityMapAPIResponse{},
			}, ReadIdentityMapConfig,
		),
		ReadContext: ReadResource(ReadIdentityMapConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "IdentityMapResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/repos/%s/identityMaps/%s/%s/%s",
						c.ControlPlane,
						d.Get("repository_id").(string),
						d.Get("identity_type").(string),
						d.Get("identity_name").(string),
						d.Get("repository_local_account_id").(string))
				},
				ResourceData: &IdentityMapResource{},
			}, ReadIdentityMapConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "IdentityMapResourceDelete",
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
			},
			"repository_local_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identity_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identity_name": {
				Type:     schema.TypeString,
				Required: true,
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
