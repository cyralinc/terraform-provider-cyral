package listener

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

var resourceContextHandler = core.DefaultContextHandler{
	ResourceName:                  resourceName,
	ResourceType:                  resourcetype.Resource,
	SchemaReaderFactory:           func() core.SchemaReader { return &SidecarListenerResource{} },
	SchemaWriterFactoryGetMethod:  func(_ *schema.ResourceData) core.SchemaWriter { return &ReadSidecarListenerAPIResponse{} },
	SchemaWriterFactoryPostMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &CreateListenerAPIResponse{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/sidecars/%s/listeners",
			c.ControlPlane,
			d.Get(utils.SidecarIDKey).(string))
	},
	GetPutDeleteURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/sidecars/%s/listeners/%s",
			c.ControlPlane,
			d.Get(utils.SidecarIDKey).(string),
			d.Get(utils.ListenerIDKey).(string),
		)
	},
}

// resourceSidecarListener returns the schema and methods for provisioning a sidecar listener
// Sidecar listeners API is {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID
// GET {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID (Get one listener)
// POST {{baseURL}}/sidecars/:sidecarID/listeners/ (Create a listener)
// PUT {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID (Update a listener)
// DELETE {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID (Delete a listener)
func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [sidecar listeners](https://cyral.com/docs/sidecars/sidecar-listeners)." +
			"\n~> **Warning** Multiple listeners can be associated to a single sidecar as long as " +
			"`host` and `port` are unique. If `host` is omitted, then `port` must be unique.",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		UpdateContext: resourceContextHandler.UpdateContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),

		Schema: getSidecarListenerSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				ids, err := utils.UnMarshalComposedID(d.Id(), "/", 2)
				if err != nil {
					return nil, err
				}
				_ = d.Set(utils.SidecarIDKey, ids[0])
				_ = d.Set(utils.ListenerIDKey, ids[1])
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func getSidecarListenerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		utils.ListenerIDKey: {
			Description: "ID of the listener that will be bound to the sidecar.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Computed:    true,
		},
		utils.SidecarIDKey: {
			Description: "ID of the sidecar that the listener will be bound to.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		RepoTypesKey: {
			Description: "List of repository types that the listener supports. Currently limited to one repo type from supported repo types:" + utils.SupportedValuesAsMarkdown(repository.RepositoryTypes()),
			Type:        schema.TypeList,
			Required:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		NetworkAddressKey: {
			Description: "The network address that the sidecar listens on.",
			Type:        schema.TypeSet,
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					utils.HostKey: {
						Description: "Host where the sidecar will listen for the given repository, in the case where the sidecar is deployed on a host " +
							"with multiple network interfaces. If omitted, the sidecar will assume the default \"0.0.0.0\" and listen on all network interfaces.",
						Type:     schema.TypeString,
						Optional: true,
					},
					utils.PortKey: {
						Description: "Port where the sidecar will listen for the given repository.",
						Type:        schema.TypeInt,
						Required:    true,
					},
				},
			},
		},
		MySQLSettingsKey: {
			Description: "MySQL settings represents the listener settings for a [`mysql`, `galera`, `mariadb`] data repository.",
			Type:        schema.TypeSet,
			Optional:    true,
			// Notice the MaxItems: 1 here. This ensures that the user can only specify one this block.
			MaxItems:      1,
			ConflictsWith: []string{S3SettingsKey, DynamoDbSettingsKey, SQLServerSettingsKey},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					DbVersionKey: {
						Description: "MySQL advertised DB version. Required (and only relevant) for listeners of " +
							"types `mysql` and `mariadb`. This value represents the MySQL/MariaDB server version that " +
							"the Cyral sidecar will use to present itself to client applications. Different applications, " +
							"especially JDBC-based ones, may behave differently according to the version of the " +
							"database they are connecting to. It is crucial that version value specified in this " +
							"field to be either the same value as the underlying database version, or to be a " +
							"compatible one. For a compatibility reference, refer to our " +
							"[public docs](https://cyral.com/docs/sidecars/manage/bind-repo). " +
							"Example values: `\"5.7.3\"`, `\"8.0.4\"` or `\"10.2.1\"`.",
						Type:     schema.TypeString,
						Optional: true,
					},
					CharacterSetKey: {
						Description: "MySQL character set. Optional (and only relevant) for listeners of " +
							"types `mysql` and `mariadb`. The sidecar automatically derives this value out of the server " +
							"version specified in the dbVersion field. This field should only be populated if the database " +
							"was configured, at deployment time, to use a global character set different from the database " +
							"default. The char set is extracted from the collation informed. The list of possible collations " +
							"can be extracted from the column `collation` by running the command `SHOW COLLATION` in " +
							"the target database.",
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		S3SettingsKey: {
			Description: "S3 settings.",
			Type:        schema.TypeSet,
			Optional:    true,
			// Notice the MaxItems: 1 here. This ensures that the user can only specify one this block.
			MaxItems:      1,
			ConflictsWith: []string{MySQLSettingsKey, DynamoDbSettingsKey, SQLServerSettingsKey},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					ProxyModeKey: {
						Description: "S3 proxy mode. Only relevant for S3 listeners. Allowed values: [true, false]. " +
							"Defaults to `false`. " +
							"When `true`, instructs the sidecar to operate as an HTTP Proxy server. Client " +
							"applications need to be explicitly configured to send the traffic through an HTTP " +
							"proxy server, represented by the Cyral sidecar endpoint + the S3 listening port. " +
							"It is indicated when connecting from CLI applications, such as `aws cli`, or through " +
							"the AWS SDK. This listener mode is functional for client applications using either " +
							"AWS native credentials, e.g. Access Key ID/Secret Access Key, or Cyral-Provided access " +
							"tokens (Single Sign-On connections). " +
							"When `false`, instructs the sidecar to mimic the actual behavior of AWS S3, meaning " +
							"client applications will not be aware of a middleware HTTP proxy in the path to S3. " +
							"This listener mode is only compatible with applications using Cyral-Provided access tokens " +
							"and is must used when configuring the Cyral S3 Browser. This mode is currently not " +
							"recommended for any other use besides the Cyral S3 Browser.",
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		DynamoDbSettingsKey: {
			Description: "DynamoDB settings.",
			Type:        schema.TypeSet,
			Optional:    true,
			// Notice the MaxItems: 1 here. This ensures that the user can only specify one this block.
			MaxItems:      1,
			ConflictsWith: []string{S3SettingsKey, MySQLSettingsKey, SQLServerSettingsKey},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					ProxyModeKey: {
						Description: "DynamoDB proxy mode. Only relevant for listeners of type `dynamodb` or " +
							"`dynamodbstreams` and must always be set to `true` for these listener types. " +
							"Defaults to false. " +
							"When `true`, instructs the sidecar to operate as an HTTP Proxy server. Client " +
							"applications need to be explicitly configured to send the traffic through an HTTP " +
							"proxy server, represented by the Cyral sidecar endpoint + the DynamoDB listening port. " +
							"It is indicated when connecting from CLI applications, such as `aws cli`, or through " +
							"the AWS SDK." +
							"Setting this value to `false` for the `dynamodb` and `dynamodbstreams` listeners types " +
							"is currently not allowed and is reserved for future use.",
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		SQLServerSettingsKey: {
			Description: "SQL Server settings.",
			Type:        schema.TypeSet,
			Optional:    true,
			// Notice the MaxItems: 1 here. This ensures that the user can only specify one this block.
			MaxItems:      1,
			ConflictsWith: []string{S3SettingsKey, MySQLSettingsKey, DynamoDbSettingsKey},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					SQLServerVersionKey: {
						Description: "Advertised SQL Server version. Required (and only relevant) for " +
							"Listeners of type 'sqlserver' " +
							"The format of the version should be <major>.<minor>.<build_number> " +
							"API will validate that the version is a valid version number. " +
							"Major version is an integer in range 0-255. " +
							"Minor version is an integer in range 0-255. " +
							"Build number is an integer in range 0-65535. " +
							"Example: 16.0.1000 " +
							"To get the version of the SQL Server runtime, run the following query: " +
							"SELECT SERVERPROPERTY('productversion') " +
							"Note: If the query returns a four part version number, only the first three parts " +
							"should be used. Example: 16.0.1000.6 -> 16.0.1000",
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}
