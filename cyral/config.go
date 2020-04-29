package cyral

// Config contains the provider configuration parameters.
type Config struct {
	Auth0Domain       string
	Auth0ClientID     string
	Auth0ClientSecret string
	JWTToken          string
	terraformVersion  string
}
