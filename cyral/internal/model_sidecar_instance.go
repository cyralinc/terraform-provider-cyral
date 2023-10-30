package internal

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SidecarInstances struct {
	Instances []SidecarInstance `json:"instances"`
}

func (wrapper *SidecarInstances) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(uuid.New().String())
	d.Set(SidecarInstanceListKey, wrapper.InstancesToInterfaceList())
	return nil
}

func (wrapper *SidecarInstances) InstancesToInterfaceList() []any {
	instancesInterfaceList := make([]any, len(wrapper.Instances))
	for index, instance := range wrapper.Instances {
		instancesInterfaceList[index] = instance.ToMap()
	}
	return instancesInterfaceList
}

type SidecarInstance struct {
	ID         string                    `json:"id"`
	Metadata   SidecarInstanceMetadata   `json:"metadata"`
	Monitoring SidecarInstanceMonitoring `json:"monitoring"`
}

func (instance *SidecarInstance) ToMap() map[string]any {
	return map[string]any{
		IDKey:         instance.ID,
		MetadataKey:   instance.Metadata.ToInterfaceList(),
		MonitoringKey: instance.Monitoring.ToInterfaceList(),
	}
}

type SidecarInstanceMetadata struct {
	Version             string               `json:"version"`
	IsDynamicVersion    bool                 `json:"isDynamicVersion"`
	SidecarCapabilities *SidecarCapabilities `json:"capabilities"`
	StartTimestamp      string               `json:"startTimestamp"`
	LastRegistration    string               `json:"lastRegistration"`
	IsRecycling         bool                 `json:"isRecycling"`
}

func (metadata *SidecarInstanceMetadata) ToInterfaceList() []any {
	return []any{
		map[string]any{
			VersionKey:          metadata.Version,
			DynamicVersionKey:   metadata.IsDynamicVersion,
			CapabilitiesKey:     metadata.SidecarCapabilities.ToInterfaceList(),
			StartTimestampKey:   metadata.StartTimestamp,
			LastRegistrationKey: metadata.LastRegistration,
			RecyclingKey:        metadata.IsRecycling,
		},
	}
}

type SidecarCapabilities struct {
	Recyclable bool `json:"recyclable"`
}

func (capabilities *SidecarCapabilities) ToInterfaceList() []any {
	if capabilities == nil {
		return nil
	}
	return []any{
		map[string]any{
			RecyclableKey: capabilities.Recyclable,
		},
	}
}

type SidecarInstanceMonitoring struct {
	Status   string                    `json:"status"`
	Services map[string]SidecarService `json:"services"`
}

func (monitoring *SidecarInstanceMonitoring) ToInterfaceList() []any {
	var services map[string]any
	if monitoring.Services != nil {
		services = make(map[string]any, len(monitoring.Services))
	}
	for serviceKey, service := range monitoring.Services {
		services[serviceKey] = service.ToMap()
	}
	return []any{
		map[string]any{
			StatusKey:   monitoring.Status,
			ServicesKey: services,
		},
	}
}

type SidecarService struct {
	Status      string                             `json:"status"`
	MetricsPort uint32                             `json:"metricsPort"`
	Components  map[string]SidecarServiceComponent `json:"components"`
	Host        string                             `json:"host"`
}

func (service *SidecarService) ToMap() map[string]any {
	var components map[string]any
	if service.Components != nil {
		components = make(map[string]any, len(service.Components))
	}
	for componentKey, component := range service.Components {
		components[componentKey] = component.ToMap()
	}
	return map[string]any{
		StatusKey:      service.Status,
		MetricsPortKey: service.MetricsPort,
		ComponentsKey:  components,
		HostKey:        service.Host,
	}
}

type SidecarServiceComponent struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Error       string `json:"error"`
}

func (component *SidecarServiceComponent) ToMap() map[string]any {
	return map[string]any{
		StatusKey:            component.Status,
		utils.DescriptionKey: component.Description,
		ErrorKey:             component.Error,
	}
}
