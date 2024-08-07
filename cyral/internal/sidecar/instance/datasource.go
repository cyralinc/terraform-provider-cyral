package instance

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve sidecar instances.",
		// The DefaultContextHandler is NOT used here as this data source intentionally
		// does not handle 404 errors, returning them to the user.
		ReadContext: core.ReadResource(core.ResourceOperationConfig{
			ResourceName: "SidecarInstanceDataSourceRead",
			Type:         operationtype.Read,
			HttpMethod:   http.MethodGet,
			URLFactory: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf(
					"https://%s/v2/sidecars/%s/instances",
					c.ControlPlane, d.Get(utils.SidecarIDKey),
				)
			},
			SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter {
				return &SidecarInstances{}
			},
		}),
		Schema: map[string]*schema.Schema{
			utils.SidecarIDKey: {
				Description: "Sidecar identifier.",
				Type:        schema.TypeString,
				Required:    true,
			},
			utils.IDKey: {
				Description: "Data source identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			SidecarInstanceListKey: {
				Description: "List of existing sidecar instances.",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						utils.IDKey: {
							Description: "Instance identifier. Varies according to the computing platform that " +
								"the sidecar is deployed to.",
							Type:     schema.TypeString,
							Computed: true,
						},
						MetadataKey: {
							Description: "Instance metadata.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									VersionKey: {
										Description: "Sidecar version that the instance is using.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									DynamicVersionKey: {
										Description: "If true, indicates that the instance has dynamic versioning, " +
											"that means that the version is not fixed at template level and it can be " +
											"automatically upgraded.",
										Type:     schema.TypeBool,
										Computed: true,
									},
									CapabilitiesKey: {
										Description: "Set of capabilities that can be enabled or disabled. **Note**: This " +
											"field is per-instance, not per-sidecar, because not all sidecar instances might be " +
											"in sync at some point in time.",
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												RecyclableKey: {
													Description: "Indicates if sidecar instance will be recycled (e.g., by an ASG) " +
														"if it reports itself as unhealthy.",
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									StartTimestampKey: {
										Description: "The time when the instance started.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									LastRegistrationKey: {
										Description: "The last time the instance reported to the Control Plane.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									RecyclingKey: {
										Description: "Indicates whether the Control Plane has asked the instance to mark " +
											"itself unhealthy so that it is recycled by the infrastructure.",
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						MonitoringKey: {
							Description: "Instance monitoring information, such as its overall health.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									utils.StatusKey: {
										Description: "Aggregated status of all the sidecar services.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									ServicesKey: {
										Description: "Sidecar instance services monitoring information.",
										Type:        schema.TypeMap,
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeSet,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													utils.StatusKey: {
														Description: "Aggregated status of sidecar service.",
														Type:        schema.TypeString,
														Computed:    true,
													},
													MetricsPortKey: {
														Description: "Metrics port for service monitoring.",
														Type:        schema.TypeInt,
														Computed:    true,
													},
													ComponentsKey: {
														Description: "Map of name to monitoring component. A component is a " +
															"monitored check on the service that has its own status.",
														Type:     schema.TypeMap,
														Computed: true,
														Elem: &schema.Schema{
															Type: schema.TypeSet,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	utils.StatusKey: {
																		Description: "Component status.",
																		Type:        schema.TypeString,
																		Computed:    true,
																	},
																	utils.DescriptionKey: {
																		Description: "Describes what the type of check the component represents.",
																		Type:        schema.TypeString,
																		Computed:    true,
																	},
																	ErrorKey: {
																		Description: "Error that describes what caused the current status.",
																		Type:        schema.TypeString,
																		Computed:    true,
																	},
																},
															},
														},
													},
													utils.HostKey: {
														Description: "Service host on the deployment.",
														Type:        schema.TypeString,
														Computed:    true,
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
