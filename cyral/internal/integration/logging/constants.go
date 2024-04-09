package logging

const (
	resourceName   = "cyral_integration_logging"
	dataSourceName = "cyral_integration_logging"
)

const (
	CloudWatchKey = "cloudwatch"
	DatadogKey    = "datadog"
	ElkKey        = "elk"
	SplunkKey     = "splunk"
	SumoLogicKey  = "sumo_logic"
	FluentbitKey  = "fluent_bit"
)

var allLogIntegrationConfigs = []string{
	CloudWatchKey,
	DatadogKey,
	ElkKey,
	SplunkKey,
	SumoLogicKey,
	FluentbitKey,
}
