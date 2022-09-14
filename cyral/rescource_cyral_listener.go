package cyral

import (
	"context"
	"fmt"
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
	SidecarId        string
	ListenerId       string           `json:"id"`
	UnixListener     UnixListener     `json:"unixListener,omitempty"`
	TCPListener      TCPListener      `json:"tcpListener,omitempty"`
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

// resourceSidecarListener returns the schema and methods for provisioning a sidecar listener
// Sidecar listeners API is {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID
// GET {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID (Get one listener)
// GET {{baseURL}}/sidecars/:sidecarID/listeners/ (Get all listeners)
// POST {{baseURL}}/sidecars/:sidecarID/listeners/ (Create a listener)
// PUT {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID (Update a listener)
// DELETE {{baseURL}}/sidecars/:sidecarID/listeners/:listenerID (Delete a listener)
func resourceSidecarListener() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages [sidecar listeners](https://cyral.com/docs/sidecars/sidecar-listeners).",
		CreateContext: resourceSidecarListenerCreate,
		ReadContext:   resourceSidecarListenerRead,
		UpdateContext: resourceSidecarListenerUpdate,
		DeleteContext: resourceSidecarListenerDelete,
		//TODO fix the correct flags for Required (should map with omitempty in struct)
		//TODO ForceNew, read what that is..
		Schema: map[string]*schema.Schema{
			"sidecar_id": {
				Description: "ID of the sidecar that the listener will be bound to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"listener_id": {
				Description: "ID of the listener that will be bound to the sidecar.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
	unixListenerFile := d.Get("unix_listener_file").(string)
	tcpListenerPort := d.Get("tcp_listener_port").(uint32)
	tcpListenerHost := d.Get("tcp_listener_host").(string)
	multiplexed := d.Get("multiplexed").(bool)
	mysqlSettingsDbVersion := d.Get("mysql_settings_db_version").(string)
	mysqlSettingsCharacterSet := d.Get("mysql_settings_character_set").(string)
	s3SettingsProxyMode := d.Get("s3_settings_proxy_mode").(bool)
	dynamoDbSettingsProxyMode := d.Get("dynamodb_settings_proxy_mode").(bool)

	listener := SidecarListener{
		SidecarId:  sidecarID,
		ListenerId: listenerID,
		UnixListener: UnixListener{
			File: unixListenerFile,
		},
		TCPListener: TCPListener{
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
