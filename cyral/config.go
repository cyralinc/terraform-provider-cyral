package cyral

// Config contains the provider configuration parameters.
type Config struct {
	auth0Domain       string
	auth0ClientID     string
	auth0ClientSecret string
	auth0Audience     string
	token             string
	tokenType         string
	controlPlane      string
	terraformVersion  string
	repoID            string
}
