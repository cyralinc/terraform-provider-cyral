package cyral

type CreateSidecarResponse struct {
	ID string `json:"ID"`
}

type SidecarData struct {
	ID                       string                   `json:"id"`
	Name                     string                   `json:"name"`
	Labels                   []string                 `json:"labels"`
	SidecarProperties        *SidecarProperties       `json:"properties"`
	ServiceConfigs           SidecarServiceConfigs    `json:"services"`
	UserEndpoint             string                   `json:"userEndpoint"`
	CertificateBundleSecrets CertificateBundleSecrets `json:"certificateBundleSecrets,omitempty"`
}

type SidecarProperties struct {
	DeploymentMethod           string `json:"deploymentMethod"`
	LogIntegrationID           string `json:"logIntegrationID,omitempty"`
	DiagnosticLogIntegrationID string `json:"diagnosticLogIntegrationID,omitempty"`
}

func NewSidecarProperties(deploymentMethod, activityLogIntegrationID, diagnosticLogIntegrationID string) *SidecarProperties {
	return &SidecarProperties{
		DeploymentMethod:           deploymentMethod,
		LogIntegrationID:           activityLogIntegrationID,
		DiagnosticLogIntegrationID: diagnosticLogIntegrationID,
	}
}

type SidecarServiceConfigs map[string]map[string]string

func (serviceConfigs *SidecarServiceConfigs) SidecarServiceConfigsAsInterfaceList() []any {
	if serviceConfigs == nil {
		return nil
	}
	serviceConfigsInterfaceList := []any{}
	for serviceName, serviceConfig := range *serviceConfigs {
		serviceConfigMap := map[string]any{
			"service_name": serviceName,
			"config":       serviceConfig,
		}
		serviceConfigsInterfaceList = append(serviceConfigsInterfaceList, serviceConfigMap)
	}
	return serviceConfigsInterfaceList
}

func (serviceConfigs *SidecarServiceConfigs) getBypassMode() string {
	if serviceConfigs != nil {
		if dispatcherConfigs, ok := (*serviceConfigs)["dispatcher"]; ok {
			if bypassMode, ok := dispatcherConfigs["bypass"]; ok {
				return bypassMode
			}
		}
	}
	return ""
}

func getSidecarServiceConfigsDefault() SidecarServiceConfigs {
	return SidecarServiceConfigs{
		"certificate-manager": nil,
		"dispatcher": {
			"bypass": "failover",
		},
		"oracle-wire": {
			"command-queue-size":       "10",
			"command-queue-timeout-ms": "200",
		},
		"pg-wire": {
			"memory-budget-enabled":         "false",
			"memory-budget-per-connection":  "8388608",
			"memory-budget-request-factor":  "512",
			"memory-budget-response-factor": "2",
		},
	}
}

type CertificateBundleSecrets map[string]*CertificateBundleSecret

type CertificateBundleSecret struct {
	Engine   string `json:"engine,omitempty"`
	SecretId string `json:"secretId,omitempty"`
	Type     string `json:"type,omitempty"`
}
