package cyral

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

// create a constant block for schema keys

const (
	SidecarIdKey        = "sidecar_id"
	ListenerIdKey       = "listener_id"
	RepoTypesKey        = "repo_types"
	UnixListenerKey     = "unix_listener"
	TcpListenerKey      = "tcp_listener"
	PortKey             = "port"
	HostKey             = "host"
	FileKey             = "file"
	MultiplexedKey      = "multiplexed"
	MysqlSettingsKey    = "mysql_settings"
	DbVersionKey        = "db_version"
	CharacterSetKey     = "character_set"
	S3SettingsKey       = "s3_settings"
	ProxyModeKey        = "proxy_mode"
	DynamoDbSettingsKey = "dynamodb_settings"
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
			d.Get(SidecarIdKey).(string),
			d.Get(ListenerIdKey).(string))
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData { return &ReadSidecarListenersAPIResponse{} },
}

type ReadSidecarListenersAPIResponse struct {
	ListenerConfig *SidecarListener `json:"listenerConfig"`
}
type CreateListenerAPIResponse struct {
	ListenerId string `json:"listenerId"`
}

func (c CreateListenerAPIResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(marshalComposedID([]string{d.Get(SidecarIdKey).(string), c.ListenerId}, "/"))
	_ = d.Set(ListenerIdKey, c.ListenerId)
	return nil
}

func (data ReadSidecarListenersAPIResponse) WriteToSchema(d *schema.ResourceData) error {
	if data.ListenerConfig != nil {
		_ = d.Set(ListenerIdKey, data.ListenerConfig.ListenerId)
		_ = d.Set(RepoTypesKey, data.ListenerConfig.RepoTypesAsInterface())
		_ = d.Set(UnixListenerKey, data.ListenerConfig.UnixSettingsAsInterface())
		_ = d.Set(TcpListenerKey, data.ListenerConfig.TcpSettingsAsInterface())
		_ = d.Set(MultiplexedKey, data.ListenerConfig.Multiplexed)
		_ = d.Set(S3SettingsKey, data.ListenerConfig.S3SettingsAsInterface())
		_ = d.Set(DynamoDbSettingsKey, data.ListenerConfig.DynamoDbSettingsAsInterface())
	}
	return nil
}
func (l *SidecarListener) RepoTypesAsInterface() *[]interface{} {
	if l.RepoTypes == nil {
		return nil
	}
	result := make([]interface{}, len(l.RepoTypes))
	for i, v := range l.RepoTypes {
		result[i] = v
	}
	return &result
}
func (l *SidecarListener) RepoTypesFromInterface(anInterface []interface{}) {
	repoTypes := make([]string, len(anInterface))
	for i, v := range anInterface {
		repoTypes[i] = v.(string)
	}
	l.RepoTypes = repoTypes
}
func (l *SidecarListener) TcpSettingsAsInterface() []interface{} {
	if l.TcpListener != nil {
		result := make([]interface{}, 1)
		result[0] = map[string]interface{}{
			HostKey: l.TcpListener.Host,
			PortKey: l.TcpListener.Port,
		}
		return result
	}
	return nil
}
func (l *SidecarListener) TcpSettingsFromInterface(anInterface []interface{}) {
	if len(anInterface) == 0 {
		return
	}
	if anInterface != nil {
		l.TcpListener = &TcpListener{
			Host: anInterface[0].(map[string]interface{})[HostKey].(string),
			Port: anInterface[0].(map[string]interface{})[PortKey].(int),
		}
	}
	log.Printf("[DEBUG] LOLOLO: %s %d", l.TcpListener.Host, l.TcpListener.Port)
}
func (l *SidecarListener) UnixSettingsAsInterface() []interface{} {
	if l.UnixListener != nil {
		result := make([]interface{}, 1)
		result[0] = map[string]interface{}{
			FileKey: l.UnixListener.File,
		}
		return result
	}
	return nil
}
func (l *SidecarListener) UnixSettingsFromInterface(anInterface []interface{}) {
	if len(anInterface) == 0 {
		return
	}
	if anInterface != nil {
		l.UnixListener = &UnixListener{
			File: anInterface[0].(map[string]interface{})[FileKey].(string),
		}
	}
}
func (l *SidecarListener) MysqlSettingsAsInterface() []interface{} {
	if l.MysqlSettings == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		DbVersionKey:    l.MysqlSettings.DbVersion,
		CharacterSetKey: l.MysqlSettings.CharacterSet,
	}}
}
func (l *SidecarListener) MysqlSettingsFromInterface(anInterface []interface{}) {
	if len(anInterface) == 0 {
		return
	}
	if anInterface != nil {
		l.MysqlSettings = &MysqlSettings{
			DbVersion:    anInterface[0].(map[string]interface{})[DbVersionKey].(string),
			CharacterSet: anInterface[0].(map[string]interface{})[CharacterSetKey].(string),
		}
	}
}
func (l *SidecarListener) S3SettingsAsInterface() []interface{} {
	if l.S3Settings == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		ProxyModeKey: l.S3Settings.ProxyMode,
	}}
}
func (l *SidecarListener) S3SettingsFromInterface(anInterface []interface{}) {
	if len(anInterface) == 0 {
		return
	}
	if anInterface != nil {
		l.S3Settings = &S3Settings{
			ProxyMode: anInterface[0].(map[string]interface{})[ProxyModeKey].(bool),
		}
	}
}
func (l *SidecarListener) DynamoDbSettingsAsInterface() []interface{} {
	if l.DynamoDbSettings == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		ProxyModeKey: l.DynamoDbSettings.ProxyMode,
	}}
}
func (l *SidecarListener) DynamoDbSettingsFromInterface(anInterface []interface{}) {
	if len(anInterface) == 0 {
		return
	}
	if anInterface != nil {
		l.DynamoDbSettings = &DynamoDbSettings{
			ProxyMode: anInterface[0].(map[string]interface{})[ProxyModeKey].(bool),
		}
	}
}

// SidecarListenerResource represents the payload of a create or update a listener request
type SidecarListenerResource struct {
	ListenerConfig SidecarListener `json:"listenerConfig"`
}

// ReadFromSchema populates the SidecarListenerResource from the schema
func (s *SidecarListenerResource) ReadFromSchema(d *schema.ResourceData) error {
	s.ListenerConfig = SidecarListener{
		SidecarId:   d.Get(SidecarIdKey).(string),
		ListenerId:  d.Get(ListenerIdKey).(string),
		Multiplexed: d.Get(MultiplexedKey).(bool),
	}
	s.ListenerConfig.RepoTypesFromInterface(d.Get(RepoTypesKey).([]interface{}))
	s.ListenerConfig.TcpSettingsFromInterface(d.Get(TcpListenerKey).(*schema.Set).List())
	s.ListenerConfig.UnixSettingsFromInterface(d.Get(UnixListenerKey).(*schema.Set).List())
	s.ListenerConfig.MysqlSettingsFromInterface(d.Get(MysqlSettingsKey).(*schema.Set).List())
	s.ListenerConfig.S3SettingsFromInterface(d.Get(S3SettingsKey).(*schema.Set).List())
	s.ListenerConfig.DynamoDbSettingsFromInterface(d.Get(DynamoDbSettingsKey).(*schema.Set).List())
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
						d.Get(SidecarIdKey).(string))

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
						d.Get(SidecarIdKey).(string),
						d.Get(ListenerIdKey).(string))

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
						d.Get(SidecarIdKey).(string),
						d.Get(ListenerIdKey).(string))
				},
			},
		),
		Schema: map[string]*schema.Schema{
			ListenerIdKey: {
				Description: "ID of the listener that will be bound to the sidecar.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
			},
			SidecarIdKey: {
				Description: "ID of the sidecar that the listener will be bound to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			RepoTypesKey: {
				Description: "List of repository types that the listener supports. Currently limited to one repo type, eg [\"mysql\"]",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			UnixListenerKey: {
				Description: "Unix listener settings.",
				Type:        schema.TypeSet,
				Optional:    true,
				// Notice the MaxItems: 1 here. This ensures that the user can only specify one this block.
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						FileKey: {
							Description: "File in which the sidecar will listen for the given repository.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			MultiplexedKey: {
				Description: "Multiplexed listener, defaults to not multiplexing (false). Not supported for all repository types.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			TcpListenerKey: {
				Description: "tcp listener settings.",
				Type:        schema.TypeSet,
				Optional:    true,
				// Notice the MaxItems: 1 here. This ensures that the user can only specify one this block.
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						HostKey: {
							Description: "Host in which the sidecar will listen for the given repository. Omit to listen on all interfaces.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						PortKey: {
							Description: "Port in which the sidecar will listen for the given repository.",
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
			},
			MysqlSettingsKey: {
				Description: "MysqlSettings represents the listener settings for a mysql data repository.",
				Type:        schema.TypeSet,
				Optional:    true,
				// Notice the MaxItems: 1 here. This ensures that the user can only specify one this block.
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						DbVersionKey: {
							Description: "MySQL DB version. Required (and only relevant) for multiplexed listeners of type mysql",
							Type:        schema.TypeString,
							Optional:    true,
						},
						CharacterSetKey: {
							Description: "MySQL character set. Optional and only relevant for multiplexed listeners of type mysql.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			S3SettingsKey: {
				Description: "S3 settings.",
				Type:        schema.TypeSet,
				Optional:    true,
				// Notice the MaxItems: 1 here. This ensures that the user can only specify one this block.
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ProxyModeKey: {
							Description: "S3 proxy mode, only relevant for S3 listeners. Defaults to false.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
			DynamoDbSettingsKey: {
				Description: "DynamoDB settings.",
				Type:        schema.TypeSet,
				Optional:    true,
				// Notice the MaxItems: 1 here. This ensures that the user can only specify one this block.
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ProxyModeKey: {
							Description: "DynamoDB proxy mode, only relevant for DynamoDB listeners. Defaults to false.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				ids, err := unmarshalComposedID(d.Id(), "/", 2)
				if err != nil {
					return nil, err
				}
				_ = d.Set(SidecarIdKey, ids[0])
				_ = d.Set(ListenerIdKey, ids[1])
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
