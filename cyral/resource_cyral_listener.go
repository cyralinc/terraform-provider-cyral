package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

// SidecarListener struct for sidecar listener.
type SidecarListener struct {
	SidecarId        string            `json:"-"`
	ListenerId       string            `json:"id"`
	RepoTypes        []string          `json:"repoTypes"`
	UnixListener     *UnixListener     `json:"unixListener,omitempty"`
	TcpListener      *TcpListener      `json:"tcpListener,omitempty"`
	Multiplexed      bool              `json:"multiplexed,omitempty"`
	MysqlSettings    *MysqlSettings    `json:"mysqlSettings,omitempty"`
	S3Settings       *S3Settings       `json:"s3Settings,omitempty"`
	DynamoDbSettings *DynamoDbSettings `json:"dynamoDBSettings,omitempty"`
}
type UnixListener struct {
	File string `json:"file"`
}
type TcpListener struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port"`
}
type MysqlSettings struct {
	DbVersion    string `json:"dbVersion,omitempty"`
	CharacterSet string `json:"characterSet,omitempty"`
}
type S3Settings struct {
	ProxyMode bool `json:"proxyMode,omitempty"`
}
type DynamoDbSettings struct {
	ProxyMode bool `json:"proxyMode,omitempty"`
}

var ReadSidecarListenersConfig = ResourceOperationConfig{
	Name:       "SidecarListenersResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/sidecars/%s/listeners/%s",
			c.ControlPlane,
			d.Get("sidecar_id").(string),
			d.Id())
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData { return &ReadSidecarListenersAPIResponse{} },
}

type ReadSidecarListenersAPIResponse struct {
	Listener *SidecarListener `json:"listener"`
}
type CreateListenerAPIResponse struct {
	ListenerId string `json:"listenerId"`
}

func (c CreateListenerAPIResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(c.ListenerId)
	_ = d.Set("listener_id", c.ListenerId)
	return nil
}

func (data ReadSidecarListenersAPIResponse) WriteToSchema(d *schema.ResourceData) error {
	if data.Listener != nil {
		_ = d.Set("sidecar_id", data.Listener.SidecarId)
		_ = d.Set("listener_id", data.Listener.ListenerId)
		_ = d.Set("repo_types", data.Listener.RepoTypes)
		_ = d.Set("unix_listener_file", data.Listener.UnixListener.File)
		_ = d.Set("tcp_listener_port", data.Listener.TcpListener.Port)
		_ = d.Set("tcp_listener_host", data.Listener.TcpListener.Host)
		_ = d.Set("multiplexed", data.Listener.Multiplexed)
		_ = d.Set("mysql_settings_db_version", data.Listener.MysqlSettings.DbVersion)
		_ = d.Set("mysql_settings_character_set", data.Listener.MysqlSettings.CharacterSet)
		_ = d.Set("s3_settings_proxy_mode", data.Listener.S3Settings.ProxyMode)
		_ = d.Set("dynamodb_settings_proxy_mode", data.Listener.DynamoDbSettings.ProxyMode)
	}
	return nil
}

// SidecarListenerResource represents the payload of a create or update a listener request
type SidecarListenerResource struct {
	ListenerConfig SidecarListener `json:"listenerConfig"`
}

// ReadFromSchema populates the SidecarListenerResource from the schema
func (s *SidecarListenerResource) ReadFromSchema(d *schema.ResourceData) error {
	var tcpListener TcpListener
	ifRepoTypes := d.Get("repo_types").([]interface{})
	repoTypes := make([]string, len(ifRepoTypes))
	for i, v := range ifRepoTypes {
		repoTypes[i] = v.(string)
	}
	s.ListenerConfig = SidecarListener{
		SidecarId:   d.Get("sidecar_id").(string),
		ListenerId:  d.Get("listener_id").(string),
		RepoTypes:   repoTypes,
		Multiplexed: d.Get("multiplexed").(bool),
	}

	if v, ok := d.GetOk("tcp_listener_port"); ok {

		tcpListener.Port = v.(int)
		if v, ok := d.GetOk("tcp_listener_host"); ok {
			tcpListener.Host = v.(string)
		}
		s.ListenerConfig.TcpListener = &tcpListener
	}
	var unixListener UnixListener
	if v, ok := d.GetOk("unix_listener_file"); ok {
		unixListener.File = v.(string)
		s.ListenerConfig.UnixListener = &unixListener
	}
	//if mysqlsettings set then set dbversion and charset
	var mysqlSettings MysqlSettings
	if v, ok := d.GetOk("mysql_settings_db_version"); ok {
		mysqlSettings.DbVersion = v.(string)
		if v, ok = d.GetOk("mysql_settings_character_set"); ok {
			mysqlSettings.CharacterSet = v.(string)
		}
		s.ListenerConfig.MysqlSettings = &mysqlSettings
	}
	// if s3settings set then set proxy mode
	var s3Settings S3Settings
	if v, ok := d.GetOk("s3_settings_proxy_mode"); ok {
		s3Settings.ProxyMode = v.(bool)
		s.ListenerConfig.S3Settings = &s3Settings
	}
	// if dynamodbsettings set then set proxy mode
	var dynamoDbSettings DynamoDbSettings
	if v, ok := d.GetOk("dynamodb_settings_proxy_mode"); ok {
		dynamoDbSettings.ProxyMode = v.(bool)
		s.ListenerConfig.DynamoDbSettings = &dynamoDbSettings
	}
	//print the struct in json format
	//get json string from struct
	jsonString, _ := json.Marshal(s)
	log.Printf("[DEBUG] ReadFromSchema, struct: %s", string(jsonString))
	return nil
}

// resourceSidecarListener returns the schema and methods for provisioning a sidecar listener
// Sidecar listeners API is {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID
// GET {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID (Get one listener)
// POST {{baseURL}}/sidecars/:sidecarID/listeners/ (Create a listener)
// PUT {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID (Update a listener)
// DELETE {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID (Delete a listener)
func resourceSidecarListener() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [sidecar listeners](https://cyral.com/docs/sidecars/sidecar-listeners).",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "SidecarListenersResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/sidecars/%s/listeners",
						c.ControlPlane,
						d.Get("sidecar_id").(string))

				},
				NewResourceData: func() ResourceData { return &SidecarListenerResource{} },
				NewResponseData: func(_ *schema.ResourceData) ResponseData { return &CreateListenerAPIResponse{} },
			}, ReadSidecarListenersConfig,
		),
		ReadContext: ReadResource(ReadSidecarListenersConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "SidecarListenersResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/sidecars/%s/listeners/%s",
						c.ControlPlane,
						d.Get("sidecar_id").(string),
						d.Id())

				},
				NewResourceData: func() ResourceData { return &SidecarListenerResource{} },
			}, ReadSidecarListenersConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "SidecarListenersResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/sidecars/%s/listeners/%s",
						c.ControlPlane,
						d.Get("sidecar_id").(string),
						d.Get("listener_id").(string))
				},
			},
		),
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Description: "ID of the listener that will be bound to the sidecar.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
			},
			"sidecar_id": {
				Description: "ID of the sidecar that the listener will be bound to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"repo_types": {
				Description: "List of repository types that the listener supports. Currently limited to one repo type, eg [\"mysql\"]",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"unix_listener_file": {
				Description: "File in which the sidecar will listen for the given repository. Required for unix listeners and mutual exclusive with tcp_listener_port.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tcp_listener_port": {
				Description: "Port in which the sidecar will listen for the given repository. Required for tcp listeners and mutual exclusive with unix_listener_file.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"tcp_listener_host": {
				Description: "Host in which the sidecar will listen for the given repository. Omit to listen on all interfaces.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"multiplexed": {
				Description: "Multiplexed listener, defaults to not multiplexing (false). Not supported for all repository types.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"mysql_settings_db_version": {
				Description: "MySQL DB version. Required (and only relevant) for multiplexed listeners of type mysql",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"mysql_settings_character_set": {
				Description: "MySQL character set. Optional and only relevant for multiplexed listeners of type mysql.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"s3_settings_proxy_mode": {
				Description: "S3 proxy mode, only relevant for S3 listeners. Defaults to false.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"dynamodb_settings_proxy_mode": {
				Description: "DynamoDB proxy mode, only relevant for DynamoDB listeners. Defaults to false.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				ids, err := unmarshalComposedID(d.Id(), "-", 2)
				if err != nil {
					return nil, err
				}
				_ = d.Set("sidecar_id", ids[0])
				d.SetId(ids[1])
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
