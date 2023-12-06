package teams

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type MsTeamsIntegration struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (data MsTeamsIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.Set("name", data.Name)
	d.Set("url", data.URL)
	return nil
}

func (data *MsTeamsIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.Name = d.Get("name").(string)
	data.URL = d.Get("url").(string)
	return nil
}
