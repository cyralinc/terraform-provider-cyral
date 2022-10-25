package cyral

type CreateSidecarResponse struct {
	ID string `json:"ID"`
}

type IdentifiedSidecarInfo struct {
	ID      string      `json:"id"`
	Sidecar SidecarData `json:"sidecar"`
}

type SidecarData struct {
	ID                       string                   `json:"id"`
	Name                     string                   `json:"name"`
	Labels                   []string                 `json:"labels"`
	Properties               *SidecarProperties       `json:"properties"`
	Services                 SidecarServicesConfig    `json:"services"`
	UserEndpoint             string                   `json:"userEndpoint"`
	CertificateBundleSecrets CertificateBundleSecrets `json:"certificateBundleSecrets,omitempty"`
}

func (sd *SidecarData) BypassMode() string {
	if sd.Services != nil {
		if dispConfig, ok := sd.Services["dispatcher"]; ok {
			if bypass_mode, ok := dispConfig["bypass"]; ok {
				return bypass_mode
			}
		}
	}
	return ""
}

type SidecarProperties struct {
	DeploymentMethod string `json:"deploymentMethod"`
}

func NewSidecarProperties(deploymentMethod string) *SidecarProperties {
	return &SidecarProperties{
		DeploymentMethod: deploymentMethod,
	}
}

type SidecarServicesConfig map[string]map[string]string

type CertificateBundleSecrets map[string]*CertificateBundleSecret

type CertificateBundleSecret struct {
	Engine   string `json:"engine,omitempty"`
	SecretId string `json:"secretId,omitempty"`
	Type     string `json:"type,omitempty"`
}
