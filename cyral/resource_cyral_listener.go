package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

// Listener struct. Proto ref: https://github.com/cyralinc/wrapper-client-go/blob/main/messages/listeners.proto#L112
// the XxxSettings and XxxListener fields are oneof in protocol buffer, I assume they should be individual attributes here.
// For these, there should be one present in payload, if more one will be picked - so all marked as omitempty.
// golang code will need to enforce this manually. Perhaps there is a assertOneOf function in the codebase?
// I also assume it is good practise to keep sidecarId out of JSON since that is only used as a path parameter
type SidecarListener struct {
	SidecarId        string           `json:"-"`
	ListenerId       string           `json:"id"`
	RepoTypes        []string         `json:"repoTypes"`
	UnixListener     UnixListener     `json:"unixListener,omitempty"`
	TcpListener      TcpListener      `json:"tcpListener,omitempty"`
	Multiplexed      bool             `json:"multiplexed,omitempty"`
	MysqlSettings    MysqlSettings    `json:"mysqlSettings,omitempty"`
	S3Settings       S3Settings       `json:"s3Settings,omitempty"`
	DynamoDbSettings DynamoDbSettings `json:"dynamoDBSettings,omitempty"`
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
			d.Get("listener_id").(string))
	},
	//here
	NewResponseData: func(_ *schema.ResourceData) ResponseData { return &SidecarListenersAPIResponse{} },
}

type SidecarListenersAPIResponse struct {
	Listener *SidecarListener `json:"listener"`
}

func (data SidecarListenersAPIResponse) WriteToSchema(d *schema.ResourceData) error {
	if data.Listener == nil {
		d.SetId(d.Get("listener_id").(string))
	} else {
		d.Set("sidecar_id", data.Listener.SidecarId)
		d.Set("listener_id", data.Listener.ListenerId)
		d.Set("repo_types", data.Listener.RepoTypes)
		d.Set("unix_listener_file", data.Listener.UnixListener.File)
		d.Set("tcp_listener_port", data.Listener.TcpListener.Port)
		d.Set("tcp_listener_host", data.Listener.TcpListener.Host)
		d.Set("multiplexed", data.Listener.Multiplexed)
		d.Set("mysql_settings_db_version", data.Listener.MysqlSettings.DbVersion)
		d.Set("mysql_settings_character_set", data.Listener.MysqlSettings.CharacterSet)
		d.Set("s3_settings_proxy_mode", data.Listener.S3Settings.ProxyMode)
		d.Set("dynamodb_settings_proxy_mode", data.Listener.DynamoDbSettings.ProxyMode)
	}
	return nil
}

// SidecarListenerResource represents the payload of a create listener request
type SidecarListenerResource struct {
	//TODO, create and update are different:
	// create MUST NOT have id set
	// update MUST have id set
	ListenerConfig SidecarListener `json:"listenerConfig"`
}

func (s *SidecarListenerResource) ReadFromSchema(d *schema.ResourceData) error {
	// for create we do not allow ListenerId to be set, and for update we mandate ListenerId to be set.
	// We need to make sure this works as expected.
	//iterate over all attributes and log them
	log.Printf("[DEBUG] ReadFromSchema, schema: %#v", d)
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
		s.ListenerConfig.TcpListener = tcpListener
	}
	var unixListener UnixListener
	if v, ok := d.GetOk("unix_listener_file"); ok {
		unixListener.File = v.(string)
		s.ListenerConfig.UnixListener = unixListener
	}
	//if mysqlsettings set then set dbversion and charset
	var mysqlSettings MysqlSettings
	if v, ok := d.GetOk("mysql_settings_db_version"); ok {
		mysqlSettings.DbVersion = v.(string)
		if v, ok = d.GetOk("mysql_settings_character_set"); ok {
			mysqlSettings.CharacterSet = v.(string)
		}
		s.ListenerConfig.MysqlSettings = mysqlSettings
	}
	// if s3settings set then set proxy mode
	var s3Settings S3Settings
	if v, ok := d.GetOk("s3_settings_proxy_mode"); ok {
		s3Settings.ProxyMode = v.(bool)
		s.ListenerConfig.S3Settings = s3Settings
	}
	// if dynamodbsettings set then set proxy mode
	var dynamoDbSettings DynamoDbSettings
	if v, ok := d.GetOk("dynamodb_settings_proxy_mode"); ok {
		dynamoDbSettings.ProxyMode = v.(bool)
		s.ListenerConfig.DynamoDbSettings = dynamoDbSettings
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
// GET {{baseURL}}/sidecars/:sidecarID/listeners/ (Get all listeners)
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
				NewResponseData: func(_ *schema.ResourceData) ResponseData { return &SidecarListenersAPIResponse{} },
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
						d.Get("listener_id").(string))

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
				Computed:    true,
			},
			"sidecar_id": {
				Description: "ID of the sidecar that the listener will be bound to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"repo_types": {
				Description: "List of repository types that the listener supports.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"unix_listener_file": {
				Description: "File in which the sidecar will listen for the given repository.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tcp_listener_port": {
				Description: "Port in which the sidecar will listen for the given repository.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"tcp_listener_host": {
				Description: "Host in which the sidecar will listen for the given repository.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"multiplexed": {
				Description: "Multiplexed listener.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"mysql_settings_db_version": {
				Description: "MySQL version.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"mysql_settings_character_set": {
				Description: "MySQL character set.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"s3_settings_proxy_mode": {
				Description: "S3 proxy mode.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"dynamodb_settings_proxy_mode": {
				Description: "DynamoDB proxy mode.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
}
func resourceSidecarListenerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	//TODO pull this out to a helper function as it's most likely re-used in other resource action functions
	sidecarID := d.Get("sidecar_id").(string)
	listenerID := d.Get("listener_id").(string)
	repoTypes := d.Get("repo_types").([]string)
	unixListenerFile := d.Get("unix_listener_file").(string)
	tcpListenerPort := d.Get("tcp_listener_port").(int)
	tcpListenerHost := d.Get("tcp_listener_host").(string)
	multiplexed := d.Get("multiplexed").(bool)
	mysqlSettingsDbVersion := d.Get("mysql_settings_db_version").(string)
	mysqlSettingsCharacterSet := d.Get("mysql_settings_character_set").(string)
	s3SettingsProxyMode := d.Get("s3_settings_proxy_mode").(bool)
	dynamoDbSettingsProxyMode := d.Get("dynamodb_settings_proxy_mode").(bool)

	listener := SidecarListener{
		SidecarId:  sidecarID,
		ListenerId: listenerID,
		RepoTypes:  repoTypes,
		UnixListener: UnixListener{
			File: unixListenerFile,
		},
		TcpListener: TcpListener{
			Port: tcpListenerPort,
			Host: tcpListenerHost,
		},
		Multiplexed: multiplexed,
		MysqlSettings: MysqlSettings{
			DbVersion:    mysqlSettingsDbVersion,
			CharacterSet: mysqlSettingsCharacterSet,
		},
		S3Settings: S3Settings{
			ProxyMode: s3SettingsProxyMode,
		},
		DynamoDbSettings: DynamoDbSettings{
			ProxyMode: dynamoDbSettingsProxyMode,
		},
	}
	url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane,
		listener.SidecarId)
	if _, err := c.DoRequest(url, http.MethodPut, listener); err != nil {
		return createError("Unable create listener resource", fmt.Sprintf("%v", err))
	}
	//TODO, check response and get id from there
	d.SetId(listenerID)
	return diags

}
func resourceSidecarListenerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//TODO
	return nil
}
func resourceSidecarListenerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//TODO
	return nil
}
func resourceSidecarListenerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//TODO
	return nil
}
