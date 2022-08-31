// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package up

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/zeebo/errs"

	"storj.io/common/uuid"
	"storj.io/storj/satellite/console/consoleauth"
)

// ConsoleEndpoint represents a user session to the web console.
type ConsoleEndpoint struct {
	client     *http.Client
	base       string
	cookieName string
	email      string
	token      string
}

// NewConsoleEndpoints creates a new client which connects to running web console.
func NewConsoleEndpoints(address string, email string) *ConsoleEndpoint {
	return &ConsoleEndpoint{
		client:     http.DefaultClient,
		base:       "http://" + address,
		cookieName: "_tokenKey",
		email:      email,
	}
}

// Login logins in to the web console (and creates use if it's necessary).
func (ce *ConsoleEndpoint) Login(ctx context.Context) (err error) {
	ce.token, err = ce.tryLogin(ctx, ce.email)
	if err != nil {
		_ = ce.tryCreateAndActivateUser(ctx, ce.email)
		ce.token, err = ce.tryLogin(ctx, ce.email)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func (ce *ConsoleEndpoint) tryLogin(ctx context.Context, email string) (string, error) {
	var authToken struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	authToken.Password = "123a123"
	authToken.Email = email

	res, err := json.Marshal(authToken)
	if err != nil {
		return "", errs.Wrap(err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ce.tokenEndpointPath(),
		bytes.NewReader(res))
	if err != nil {
		return "", errs.Wrap(err)
	}

	request.Header.Add("Content-Type", "application/json")

	resp, err := ce.client.Do(request)
	if err != nil {
		return "", errs.Wrap(err)
	}
	defer func() { err = errs.Combine(err, resp.Body.Close()) }()

	if resp.StatusCode != http.StatusOK {
		return "", errs.New("unexpected status code: %d (%q)",
			resp.StatusCode, tryReadLine(resp.Body))
	}

	var tokenInfo struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&tokenInfo)
	if err != nil {
		return "", errs.Wrap(err)
	}

	return tokenInfo.Token, nil
}

func (ce *ConsoleEndpoint) tryCreateAndActivateUser(ctx context.Context, email string) error {
	regToken, err := ce.createRegistrationToken(ctx)
	if err != nil {
		return errs.Wrap(err)
	}
	userID, err := ce.createUser(ctx, regToken, email)
	if err != nil {
		return errs.Wrap(err)
	}
	return errs.Wrap(ce.activateUser(ctx, userID, email))
}

func (ce *ConsoleEndpoint) createRegistrationToken(ctx context.Context) (string, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		ce.registrationTokenPath(),
		nil)
	if err != nil {
		return "", errs.Wrap(err)
	}

	resp, err := ce.client.Do(request)
	if err != nil {
		return "", errs.Wrap(err)
	}
	defer func() { err = errs.Combine(err, resp.Body.Close()) }()

	if resp.StatusCode != http.StatusOK {
		return "", errs.New("unexpected status code: %d (%q)",
			resp.StatusCode, tryReadLine(resp.Body))
	}

	var createTokenResponse struct {
		Secret string
		Error  string
	}
	if err = json.NewDecoder(resp.Body).Decode(&createTokenResponse); err != nil {
		return "", errs.Wrap(err)
	}
	if createTokenResponse.Error != "" {
		return "", errs.New("unable to create registration token: %s", createTokenResponse.Error)
	}

	return createTokenResponse.Secret, nil
}

func (ce *ConsoleEndpoint) createUser(ctx context.Context, regToken string, email string) (string, error) {
	var registerData struct {
		FullName  string `json:"fullName"`
		ShortName string `json:"shortName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Secret    string `json:"secret"`
	}

	registerData.FullName = "Alice"
	registerData.Email = email
	registerData.Password = "123a123"
	registerData.ShortName = "al"
	registerData.Secret = regToken

	res, err := json.Marshal(registerData)
	if err != nil {
		return "", errs.Wrap(err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ce.registerEndpointPath(),
		bytes.NewReader(res))
	if err != nil {
		return "", errs.Wrap(err)
	}
	request.Header.Add("Content-Type", "application/json")

	resp, err := ce.client.Do(request)
	if err != nil {
		return "", errs.Wrap(err)
	}
	defer func() { err = errs.Combine(err, resp.Body.Close()) }()

	if resp.StatusCode != http.StatusOK {
		return "", errs.New("unexpected status code: %d (%q)",
			resp.StatusCode, tryReadLine(resp.Body))
	}

	var userID string
	if err = json.NewDecoder(resp.Body).Decode(&userID); err != nil {
		return "", errs.Wrap(err)
	}

	return userID, nil
}

func (ce *ConsoleEndpoint) activateUser(ctx context.Context, userID string, email string) error {
	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return errs.Wrap(err)
	}

	activationToken, err := generateActivationKey(userUUID, email, time.Now())
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		ce.activationEndpointPath(activationToken),
		nil)
	if err != nil {
		return err
	}

	resp, err := ce.client.Do(request)
	if err != nil {
		return err
	}
	defer func() { err = errs.Combine(err, resp.Body.Close()) }()

	if resp.StatusCode != http.StatusOK {
		return errs.New("unexpected status code: %d (%q)",
			resp.StatusCode, tryReadLine(resp.Body))
	}

	return nil
}

// GetOrCreateProject return with project to use in this session.
func (ce *ConsoleEndpoint) GetOrCreateProject(ctx context.Context) (string, string, error) {
	projectID, token, err := ce.getProject(ctx)
	if err == nil {
		return projectID, token, nil
	}
	projectID, token, err = ce.createProject(ctx)
	if err == nil {
		return projectID, token, nil
	}
	return ce.getProject(ctx)
}

func (ce *ConsoleEndpoint) getProject(ctx context.Context) (string, string, error) {
	query := `query {myProjects{id}}`
	var getProjects struct {
		MyProjects []struct {
			ID string
		}
	}
	err := ce.graphqlQuery(ctx, query, &getProjects)
	if len(getProjects.MyProjects) == 0 {
		return "", "", errs.New("No project exists")
	}
	return getProjects.MyProjects[0].ID, ce.token, err
}

func (ce *ConsoleEndpoint) createProject(ctx context.Context) (string, string, error) {
	rng := rand.NewSource(time.Now().UnixNano())
	createProjectQuery := fmt.Sprintf(
		`mutation {createProject(input:{name:"TestProject-%d",description:""}){id}}`,
		rng.Int63())

	var createProject struct {
		CreateProject struct {
			ID string
		}
	}
	err := ce.graphqlMutation(ctx, createProjectQuery, &createProject)
	return createProject.CreateProject.ID, ce.token, err
}

// CreateAPIKey creates new API key to access Storj services.
func (ce *ConsoleEndpoint) CreateAPIKey(ctx context.Context, projectID string) (string, error) {
	rng := rand.NewSource(time.Now().UnixNano())
	createAPIKeyQuery := fmt.Sprintf(
		`mutation {createAPIKey(projectID:%q,name:"TestKey-%d"){key}}`,
		projectID, rng.Int63())

	var createAPIKey struct {
		CreateAPIKey struct {
			Key string
		}
	}
	err := ce.graphqlMutation(ctx, createAPIKeyQuery, &createAPIKey)
	return createAPIKey.CreateAPIKey.Key, err
}

func (ce *ConsoleEndpoint) graphqlQuery(ctx context.Context, createAPIKeyQuery string, response interface{}) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		ce.graphQLEndpointPath(),
		nil)
	if err != nil {
		return errs.Wrap(err)
	}

	q := request.URL.Query()
	q.Add("query", createAPIKeyQuery)
	request.URL.RawQuery = q.Encode()

	request.AddCookie(&http.Cookie{
		Name:  ce.cookieName,
		Value: ce.token,
	})

	request.Header.Add("Content-Type", "application/graphql")

	return ce.graphqlDo(request, response)
}

func (ce *ConsoleEndpoint) graphqlMutation(ctx context.Context, query string, response interface{}) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ce.graphQLEndpointPath(),
		bytes.NewReader([]byte(query)))
	if err != nil {
		return errs.Wrap(err)
	}

	request.AddCookie(&http.Cookie{
		Name:  ce.cookieName,
		Value: ce.token,
	})

	request.Header.Add("Content-Type", "application/graphql")

	return ce.graphqlDo(request, response)
}

func generateActivationKey(userID uuid.UUID, email string, createdAt time.Time) (string, error) {
	claims := consoleauth.Claims{
		ID:         userID,
		Email:      email,
		Expiration: createdAt.Add(24 * time.Hour),
	}

	// TODO: change it in future, when satellite/console secret will be changed
	signer := &consoleauth.Hmac{Secret: []byte("my-suppa-secret-key")}

	resJSON, err := claims.JSON()
	if err != nil {
		return "", err
	}

	token := consoleauth.Token{Payload: resJSON}
	encoded := base64.URLEncoding.EncodeToString(token.Payload)

	signature, err := signer.Sign([]byte(encoded))
	if err != nil {
		return "", err
	}

	token.Signature = signature

	return token.String(), nil
}

func tryReadLine(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	return scanner.Text()
}

func (ce *ConsoleEndpoint) graphqlDo(request *http.Request, jsonResponse interface{}) error {
	resp, err := ce.client.Do(request)
	if err != nil {
		return err
	}
	defer func() { err = errs.Combine(err, resp.Body.Close()) }()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response struct {
		Data   json.RawMessage
		Errors []interface{}
	}

	if err = json.NewDecoder(bytes.NewReader(b)).Decode(&response); err != nil {
		return err
	}

	if response.Errors != nil {
		return errs.New("inner graphql error: %v", response.Errors)
	}

	if jsonResponse == nil {
		return errs.New("empty response: %q", b)
	}

	return json.NewDecoder(bytes.NewReader(response.Data)).Decode(jsonResponse)
}

func (ce *ConsoleEndpoint) appendPath(suffix string) string {
	return ce.base + suffix
}

func (ce *ConsoleEndpoint) registrationTokenPath() string {
	return ce.appendPath("/registrationToken/?projectsLimit=1")
}

func (ce *ConsoleEndpoint) registerEndpointPath() string {
	return ce.appendPath("/api/v0/auth/register")
}

func (ce *ConsoleEndpoint) activationEndpointPath(token string) string {
	return ce.appendPath("/activation/?token=" + token)
}

func (ce *ConsoleEndpoint) tokenEndpointPath() string {
	return ce.appendPath("/api/v0/auth/token")
}

func (ce *ConsoleEndpoint) graphQLEndpointPath() string {
	return ce.appendPath("/api/v0/graphql")
}

// RegisterAccess creates new access registered to linksharing.
func RegisterAccess(ctx context.Context, authService string, accessSerialized string) (accessKey, secretKey, endpoint string, err error) {
	if authService == "" {
		return "", "", "", errs.New("no auth service address provided")
	}

	postData, err := json.Marshal(map[string]interface{}{
		"access_grant": accessSerialized,
		"public":       false,
	})
	if err != nil {
		return accessKey, "", "", errs.Wrap(err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/v1/access", authService), bytes.NewReader(postData))
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer func() { err = errs.Combine(err, resp.Body.Close()) }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	respBody := make(map[string]string)
	if err := json.Unmarshal(body, &respBody); err != nil {
		return "", "", "", errs.New("unexpected response from auth service: %s", string(body))
	}

	accessKey, ok := respBody["access_key_id"]
	if !ok {
		return "", "", "", errs.New("access_key_id missing in response")
	}
	secretKey, ok = respBody["secret_key"]
	if !ok {
		return "", "", "", errs.New("secret_key missing in response")
	}
	return accessKey, secretKey, respBody["endpoint"], nil
}
