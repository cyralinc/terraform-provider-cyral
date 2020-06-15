package cyral

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestEnvVarAuth0ClientID(t *testing.T) {
	c := Config{}
	if _, err := c.Client(); err != nil {
		if _, err2 := c.getEnv("AUTH0_CLIENT_ID"); err.Error() != err2.Error() {
			t.Error("Expected error in env var: AUTH0_CLIENT_ID")
		}
	}
}
func TestEnvVarAuth0ClientSecret(t *testing.T) {
	c := Config{}
	os.Setenv("AUTH0_CLIENT_ID", "bla")
	if _, err := c.Client(); err != nil {
		if _, err2 := c.getEnv("AUTH0_CLIENT_SECRET"); err.Error() != err2.Error() {
			t.Error("Expected error in env var: AUTH0_CLIENT_SECRET")
		}
	}
}

// FALTA: MOCK CLIENT (E ONDE ELE APARECE)

/*
//rascunho
type HTTPClient interface {
	Do(req *http.Request)(*http.Response, error)
}

type MockClient http.Client

func (mock *MockClient) Do(req *http.Request) (respAddr *http.Response, error) {
	auth0Token := auth0TokenResponse{AccessToken: "", TokenType:""}
	body, err := json.Marshal(auth0Token)

	respAddr := &http.Response{}
	// erro ?
	//tratar req: erro ou não
	respAddr := New(http.Response)
	*respAddr = http.Response{}
	// erro = ?
	json.Marshal()

}
//fimrascunho
//**************** implementar a interface Do mock
*/

func (confAddr *Config) discoverTokenError(clientID, clientSecret string) error {
	url := fmt.Sprintf("https://%s/oauth/token", confAddr.Auth0Domain)
	audienceURL := fmt.Sprintf("https://%s", confAddr.Auth0Audience)
	tokenReq := auth0TokenRequest{
		Audience:     audienceURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	payloadBytes, err := json.Marshal(tokenReq)
	if err != nil {
		return fmt.Errorf("Token error, failed to encode readToken payload: %v.\n", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return fmt.Errorf("Token error, unable to create auth0 request; err: %v", err)
	}
	req.Header.Add("content-type", "application/json")
	//***********************
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Token error, unable execute auth0 request; err: %v", err)
	}
	defer res.Body.Close() ///*****************
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Token error, unable to read data from request body; err: %v", err)
	}
	token := auth0TokenResponse{}
	err = json.Unmarshal(body, &token)
	//if err != nil {
	return fmt.Errorf("Token error, unable to get access token from json; err: %v", err)
	//}
}

func (confAddr *Config) checkEnvVar(envVar string) (string, error) {
	varValue, err := confAddr.getEnv(envVar)
	if err != nil {
		return "", err
	}
	return varValue, nil
}

func (confAddr *Config) traceProblem() error {
	auth0ClientID, err := confAddr.checkEnvVar("AUTH0_CLIENT_ID")
	if err != nil {
		return err
	}
	auth0ClientSecret, err := confAddr.checkEnvVar("AUTH0_CLIENT_SECRET")
	if err != nil {
		return err
	}
	return confAddr.discoverTokenError(auth0ClientID, auth0ClientSecret)
}

func TestClient(t *testing.T) {

	mockConfig := Config{}

	/* //teste
	t.Logf("For test purposes:\n")
	t.Logf("Auth0Domain: %s\n", mockConfig.Auth0Domain)
	t.Logf("Auth0Audience: %s\n", mockConfig.Auth0Audience)
	t.Logf("controlPlane: %s\n", mockConfig.controlPlane)
	t.Logf("controlPlaneAPIPort: %d\n", mockConfig.controlPlaneAPIPort)
	t.Logf("terraformVersion: %s\n", mockConfig.terraformVersion)
	//fimteste */

	_, err := (&mockConfig).Client()

	if err != nil {
		// ARRUMAR PARA MSGS SIGNIFICATIVAS!!!
		// exemplo: "erro ao inicialiazar abjshas. + msg erro comum."
		//t.Log ( (&mockConfig).traceProblem() )
		t.Error((&mockConfig).traceProblem())
	}
}

//func Test
