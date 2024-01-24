// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package up

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/zeebo/errs"

	"storj.io/common/uuid"
	"storj.io/storj/satellite/console/consoleauth"
)

var errDecodeResponse = errors.New("unable to decode json response")

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

	return getToken(resp)
}

func getToken(resp *http.Response) (string, error) {
	var tokenInfo struct {
		Token string `json:"token"`
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errs.Wrap(err)
	}
	err = json.Unmarshal(responseBody, &tokenInfo)
	if err != nil {
		// nolint
		// before https://review.dev.storj.io/c/storj/storj/+/8033
		return strings.Trim(string(responseBody), "\n\""), nil
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
// this method will try to get/create the project using the old graphql endpoint,
// and if it fails, it will try to get/create the project using the new http endpoint.
// see the changes in https://github.com/storj/storj/commit/516241e406923dedcc66df06b7e7c1479dc98b91
// for the update from old API to new.
// TODO: update this method when the old API is no longer needed.
func (ce *ConsoleEndpoint) GetOrCreateProject(ctx context.Context) (string, string, error) {
	projectID, token, err := ce.getGraphqlProject(ctx)
	if errors.Is(err, io.EOF) || errors.Is(err, errDecodeResponse) {
		projectID, token, err = ce.getHttpProject(ctx)
	}
	if err == nil {
		return projectID, token, nil
	}
	projectID, token, err = ce.createGraphqlProject(ctx)
	if errors.Is(err, io.EOF) || errors.Is(err, errDecodeResponse) {
		projectID, token, err = ce.createHttpProject(ctx)
	}
	if err == nil {
		return projectID, token, nil
	}
	return "", "", err
}

func (ce *ConsoleEndpoint) getHttpProject(ctx context.Context) (string, string, error) {
	var projects []struct {
		ID string `json:"id"`
	}
	err := ce.projectQuery(ctx, &projects)
	if err != nil {
		return "", "", err
	}
	if len(projects) == 0 {
		return "", "", errs.New("No project exists")
	}
	return projects[0].ID, ce.token, nil
}

func (ce *ConsoleEndpoint) createHttpProject(ctx context.Context) (string, string, error) {
	rng := rand.NewSource(time.Now().UnixNano())
	body := fmt.Sprintf(`{"name":"TestProject-%d","description":""}`, rng.Int63())

	var createdProject struct {
		ID string `json:"id"`
	}
	err := ce.projectMutation(ctx, body, &createdProject)
	if err != nil {
		return "", "", err
	}
	return createdProject.ID, ce.token, nil
}

func (ce *ConsoleEndpoint) getGraphqlProject(ctx context.Context) (string, string, error) {
	query := `query {myProjects{id}}`
	var getProjects struct {
		MyProjects []struct {
			ID string
		}
	}
	err := ce.graphqlQuery(ctx, query, &getProjects)
	if err != nil {
		return "", "", err
	}
	if len(getProjects.MyProjects) == 0 {
		return "", "", errs.New("No project exists")
	}
	return getProjects.MyProjects[0].ID, ce.token, nil
}

func (ce *ConsoleEndpoint) createGraphqlProject(ctx context.Context) (string, string, error) {
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
	if err != nil {
		return "", "", err
	}
	return createProject.CreateProject.ID, ce.token, nil
}

// CreateAPIKey creates new API key to access Storj services.
// this method will try to create the key using the old graphql endpoint,
// and if it fails, it will try to get the project using the new http endpoint.
// see the changes in https://github.com/storj/storj/commit/516241e406923dedcc66df06b7e7c1479dc98b91
// for the update from old API to new.
// TODO: update this method when the old API is no longer needed.
func (ce *ConsoleEndpoint) CreateAPIKey(ctx context.Context, projectID string) (string, error) {
	key, err := ce.createGraphqlAPIKey(ctx, projectID)
	if errors.Is(err, io.EOF) || errors.Is(err, errDecodeResponse) {
		key, err = ce.createAPIKey(ctx, projectID)
	}
	return key, err
}

func (ce *ConsoleEndpoint) createGraphqlAPIKey(ctx context.Context, projectID string) (string, error) {
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
	if err != nil {
		return "", err
	}
	return createAPIKey.CreateAPIKey.Key, nil
}

func (ce *ConsoleEndpoint) createAPIKey(ctx context.Context, projectID string) (string, error) {
	rng := rand.NewSource(time.Now().UnixNano())
	apiKeyName := fmt.Sprintf("TestKey-%d", rng.Int63())

	var createdKey struct {
		Key string `json:"key"`
	}
	err := ce.apiKeyMutation(ctx, apiKeyName, projectID, &createdKey)
	if err != nil {
		return "", err
	}
	return createdKey.Key, nil
}

func (ce *ConsoleEndpoint) projectQuery(ctx context.Context, response interface{}) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		ce.projectEndpointPath(),
		nil)
	if err != nil {
		return errs.Wrap(err)
	}

	request.AddCookie(&http.Cookie{
		Name:  ce.cookieName,
		Value: ce.token,
	})

	request.Header.Add("Content-Type", "application/json")

	return ce.httpDo(request, response)
}

func (ce *ConsoleEndpoint) projectMutation(ctx context.Context, query string, response interface{}) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ce.projectEndpointPath(),
		bytes.NewReader([]byte(query)))
	if err != nil {
		return errs.Wrap(err)
	}

	request.AddCookie(&http.Cookie{
		Name:  ce.cookieName,
		Value: ce.token,
	})

	request.Header.Add("Content-Type", "application/json")

	return ce.httpDo(request, response)
}

func (ce *ConsoleEndpoint) apiKeyMutation(ctx context.Context, apiKeyName, projectID string, response interface{}) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ce.apiKeyEndpointPath()+"/create/"+projectID,
		bytes.NewReader([]byte(apiKeyName)))
	if err != nil {
		return errs.Wrap(err)
	}

	request.AddCookie(&http.Cookie{
		Name:  ce.cookieName,
		Value: ce.token,
	})

	request.Header.Add("Content-Type", "application/json")

	return ce.httpDo(request, response)
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

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response struct {
		Data   json.RawMessage
		Errors []interface{}
	}

	if err = json.NewDecoder(bytes.NewReader(b)).Decode(&response); err != nil {
		return errDecodeResponse
	}

	if response.Errors != nil {
		return errs.New("inner graphql error: %v", response.Errors)
	}

	if jsonResponse == nil {
		return errs.New("empty response: %q", b)
	}

	return json.NewDecoder(bytes.NewReader(response.Data)).Decode(jsonResponse)
}

func (ce *ConsoleEndpoint) httpDo(request *http.Request, jsonResponse interface{}) error {
	resp, err := ce.client.Do(request)
	if err != nil {
		return err
	}
	defer func() { err = errs.Combine(err, resp.Body.Close()) }()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if jsonResponse == nil {
		return errs.New("empty response: %q", b)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return json.NewDecoder(bytes.NewReader(b)).Decode(jsonResponse)
	}

	var errResponse struct {
		Error string `json:"error"`
	}

	err = json.NewDecoder(bytes.NewReader(b)).Decode(&errResponse)
	if err != nil {
		return err
	}

	return errs.New("request failed with status %d: %s", resp.StatusCode, errResponse.Error)
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

func (ce *ConsoleEndpoint) projectEndpointPath() string {
	return ce.appendPath("/api/v0/projects")
}

func (ce *ConsoleEndpoint) apiKeyEndpointPath() string {
	return ce.appendPath("/api/v0/api-keys")
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

	body, err := io.ReadAll(resp.Body)
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
