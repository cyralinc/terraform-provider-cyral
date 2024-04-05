package listener

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"

	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

// SidecarListener struct for sidecar listener.
type SidecarListener struct {
	SidecarId         string             `json:"-"`
	ListenerId        string             `json:"id"`
	RepoTypes         []string           `json:"repoTypes"`
	NetworkAddress    *NetworkAddress    `json:"address,omitempty"`
	MySQLSettings     *MySQLSettings     `json:"mysqlSettings,omitempty"`
	S3Settings        *S3Settings        `json:"s3Settings,omitempty"`
	DynamoDbSettings  *DynamoDbSettings  `json:"dynamoDbSettings,omitempty"`
	SQLServerSettings *SQLServerSettings `json:"sqlServerSettings,omitempty"`
}

type NetworkAddress struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port"`
}

type MySQLSettings struct {
	DbVersion    string `json:"dbVersion,omitempty"`
	CharacterSet string `json:"characterSet,omitempty"`
}

type S3Settings struct {
	ProxyMode bool `json:"proxyMode,omitempty"`
}

type DynamoDbSettings struct {
	ProxyMode bool `json:"proxyMode,omitempty"`
}

type SQLServerSettings struct {
	Version string `json:"version,omitempty"`
}

type ReadSidecarListenerAPIResponse struct {
	ListenerConfig *SidecarListener `json:"listenerConfig"`
}
type CreateListenerAPIResponse struct {
	ListenerId string `json:"listenerId"`
}

func (c CreateListenerAPIResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(utils.MarshalComposedID([]string{d.Get(utils.SidecarIDKey).(string), c.ListenerId}, "/"))
	return d.Set(utils.ListenerIDKey, c.ListenerId)
}

func (data ReadSidecarListenerAPIResponse) WriteToSchema(d *schema.ResourceData) error {
	ctx := context.Background()
	tflog.Debug(ctx, "Init ReadSidecarListenerAPIResponse.WriteToSchema")
	if data.ListenerConfig != nil {
		_ = d.Set(utils.ListenerIDKey, data.ListenerConfig.ListenerId)
		_ = d.Set(RepoTypesKey, data.ListenerConfig.RepoTypesAsInterface())
		_ = d.Set(NetworkAddressKey, data.ListenerConfig.NetworkAddressAsInterface())
		_ = d.Set(S3SettingsKey, data.ListenerConfig.S3SettingsAsInterface())
		_ = d.Set(MySQLSettingsKey, data.ListenerConfig.MySQLSettingsAsInterface())
		_ = d.Set(DynamoDbSettingsKey, data.ListenerConfig.DynamoDbSettingsAsInterface())
		_ = d.Set(SQLServerSettingsKey, data.ListenerConfig.SQLServerSettingsAsInterface())
	}
	tflog.Debug(ctx, "End ReadSidecarListenerAPIResponse.WriteToSchema")
	return nil
}

func (l *SidecarListener) RepoTypesAsInterface() []interface{} {
	if l.RepoTypes == nil {
		return nil
	}
	result := make([]interface{}, len(l.RepoTypes))
	for i, v := range l.RepoTypes {
		result[i] = v
	}
	return result
}
func (l *SidecarListener) RepoTypesFromInterface(anInterface []interface{}) {
	repoTypes := make([]string, len(anInterface))
	for i, v := range anInterface {
		repoTypes[i] = v.(string)
	}
	l.RepoTypes = repoTypes
}
func (l *SidecarListener) NetworkAddressAsInterface() []interface{} {
	if l.NetworkAddress == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			utils.HostKey: l.NetworkAddress.Host,
			utils.PortKey: l.NetworkAddress.Port,
		},
	}
}
func (l *SidecarListener) NetworkAddressFromInterface(anInterface []interface{}) {
	if len(anInterface) == 0 {
		return
	}
	l.NetworkAddress = &NetworkAddress{
		Host: anInterface[0].(map[string]interface{})[utils.HostKey].(string),
		Port: anInterface[0].(map[string]interface{})[utils.PortKey].(int),
	}
}
func (l *SidecarListener) MySQLSettingsAsInterface() []interface{} {
	if l.MySQLSettings == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		DbVersionKey:    l.MySQLSettings.DbVersion,
		CharacterSetKey: l.MySQLSettings.CharacterSet,
	}}
}
func (l *SidecarListener) MySQLSettingsFromInterface(anInterface []interface{}) {
	if len(anInterface) == 0 {
		return
	}
	l.MySQLSettings = &MySQLSettings{
		DbVersion:    anInterface[0].(map[string]interface{})[DbVersionKey].(string),
		CharacterSet: anInterface[0].(map[string]interface{})[CharacterSetKey].(string),
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
	l.S3Settings = &S3Settings{
		ProxyMode: anInterface[0].(map[string]interface{})[ProxyModeKey].(bool),
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
	l.DynamoDbSettings = &DynamoDbSettings{
		ProxyMode: anInterface[0].(map[string]interface{})[ProxyModeKey].(bool),
	}
}
func (l *SidecarListener) SQLServerSettingsAsInterface() []interface{} {
	if l.SQLServerSettings == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		SQLServerVersionKey: l.SQLServerSettings.Version,
	}}
}
func (l *SidecarListener) SQLServerSettingsFromInterface(anInterface []interface{}) {
	if len(anInterface) == 0 {
		return
	}
	l.SQLServerSettings = &SQLServerSettings{
		Version: anInterface[0].(map[string]interface{})[SQLServerVersionKey].(string),
	}
}

// SidecarListenerResource represents the payload of a create or update a listener request
type SidecarListenerResource struct {
	ListenerConfig SidecarListener `json:"listenerConfig"`
}

// ReadFromSchema populates the SidecarListenerResource from the schema
func (s *SidecarListenerResource) ReadFromSchema(d *schema.ResourceData) error {
	s.ListenerConfig = SidecarListener{
		SidecarId:  d.Get(utils.SidecarIDKey).(string),
		ListenerId: d.Get(utils.ListenerIDKey).(string),
	}
	s.ListenerConfig.RepoTypesFromInterface(d.Get(RepoTypesKey).([]interface{}))
	s.ListenerConfig.NetworkAddressFromInterface(d.Get(NetworkAddressKey).(*schema.Set).List())
	s.ListenerConfig.MySQLSettingsFromInterface(d.Get(MySQLSettingsKey).(*schema.Set).List())
	s.ListenerConfig.S3SettingsFromInterface(d.Get(S3SettingsKey).(*schema.Set).List())
	s.ListenerConfig.DynamoDbSettingsFromInterface(d.Get(DynamoDbSettingsKey).(*schema.Set).List())
	s.ListenerConfig.SQLServerSettingsFromInterface(d.Get(SQLServerSettingsKey).(*schema.Set).List())

	return nil
}

type ReadDataSourceSidecarListenerAPIResponse struct {
	ListenerConfigs []SidecarListener `json:"listenerConfigs"`
}

func (data ReadDataSourceSidecarListenerAPIResponse) WriteToSchema(d *schema.ResourceData) error {
	ctx := context.Background()
	tflog.Debug(ctx, "Init ReadDataSourceSidecarListenerAPIResponse.WriteToSchema")
	var listenersList []any
	tflog.Debug(ctx, fmt.Sprintf("data.ListenerConfig: %+v", data.ListenerConfigs))
	tflog.Debug(ctx, "Init for _, l := range data.ListenerConfig")
	repoTypeFilter := d.Get(DSRepoTypeKey).(string)
	portFilter := d.Get(utils.PortKey).(int)
	for _, listenerConfig := range data.ListenerConfigs {
		// Check if either the repo filter or the port filter is provided and matches the listener
		if (repoTypeFilter == "" || slices.Contains(listenerConfig.RepoTypes, repoTypeFilter)) &&
			(portFilter == 0 || listenerConfig.NetworkAddress.Port == portFilter) {
			listener := map[string]any{
				utils.ListenerIDKey:  listenerConfig.ListenerId,
				utils.SidecarIDKey:   d.Get(utils.SidecarIDKey).(string),
				RepoTypesKey:         listenerConfig.RepoTypes,
				NetworkAddressKey:    listenerConfig.NetworkAddressAsInterface(),
				MySQLSettingsKey:     listenerConfig.MySQLSettingsAsInterface(),
				S3SettingsKey:        listenerConfig.S3SettingsAsInterface(),
				DynamoDbSettingsKey:  listenerConfig.DynamoDbSettingsAsInterface(),
				SQLServerSettingsKey: listenerConfig.SQLServerSettingsAsInterface(),
			}
			tflog.Debug(ctx, fmt.Sprintf("listener: %q", listener))
			listenersList = append(listenersList, listener)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("listenersList: %q", listenersList))
	tflog.Debug(ctx, "End for _, l := range data.ListenerConfig")

	if err := d.Set(SidecarListenerListKey, listenersList); err != nil {
		return err
	}

	d.SetId(uuid.New().String())

	tflog.Debug(ctx, "End ReadDataSourceSidecarListenerAPIResponse.WriteToSchema")

	return nil
}
