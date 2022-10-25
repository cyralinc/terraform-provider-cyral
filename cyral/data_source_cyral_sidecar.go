package cyral

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type dataSourceSidecarListResponse []IdentifiedSidecarInfo

func (resp dataSourceSidecarListResponse) WriteToSchema(d *schema.ResourceData) error {
	sidecarList := make([]interface{}, 0, len(resp))

	for _, sidecarInfo := range resp {
		sidecar := sidecarInfo.Sidecar
		sidecarInfoMap := map[string]interface{}{
			"id": sidecarInfo.ID,
			"sidecar": []interface{}{map[string]interface{}{
				"name":                       sidecar.Name,
				"labels":                     sidecar.Labels,
				"properties":                 sidecar.Properties,
				"services":                   sidecar.Services,
				"user_endpoint":              sidecar.UserEndpoint,
				"certificate_bundle_secrets": sidecar.CertificateBundleSecrets,
			}},
		}
		sidecarList = append(sidecarList, sidecarInfoMap)
	}

	return d.Set("sidecar_list", sidecarList)
}

func dataSourceSidecar() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter [sidecars](https://cyral.com/docs/sidecars/sidecar-install) that exist in the Cyral Control Plane.",
		ReadContext: ReadResource(dataSourceRoleReadConfig()),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Filter the results by a regular expression (regex) that matches names of existing roles.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"sidecar_list": {
				Description: "List of existing sidecars satisfying given filter criteria.",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "ID of the sidecar in the Cyral environment.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"sidecar": {
							Description: "Sidecar information.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "Sidecar name.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"cloud": {
										Description: "Cloud provider the sidecar is deployed in.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"endpoint": {
										Description: "Sidecar endpoint (DNS name or IP address).",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"user_endpoint": {
										Description: "User-defined endpoint (also referred as `alias`) that can be used to " +
											"override the sidecar DNS endpoint shown in the UI.",
										Type:     schema.TypeString,
										Computed: true,
									},
									"services": {
										Description: "Special configuration for each of the sidecar services.",
										Type:        schema.TypeMap,
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"labels": {
										Description: "Labels attached to the sidecar.",
										Type:        schema.TypeList,
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"properties": {
										Description: "Special configuration for the sidecar.",
										Type:        schema.TypeMap,
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"certificate_bundle_secrets": {
										Description: "Configuration that holds data about the location of a particular TLS certificate " +
											"bundle in a secrets manager.",
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeSet,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"engine": {
														Description: "Name of the engine used with the given secrets.",
														Type:        schema.TypeString,
														Computed:    true,
													},
													"secret_id": &schema.Schema{
														Description: "Identifier or location of the secret.",
														Type:        schema.TypeString,
														Computed:    true,
													},
													"type": {
														Description: "Type identifies the secret manager used to store the secret. " +
															"For example: `aws`, `k8s`.",
														Type:     schema.TypeString,
														Computed: true,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
