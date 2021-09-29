// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/zeebo/errs"

	"storj.io/common/uuid"
	"storj.io/storj/satellite/console/consoleauth"
	"storj.io/storj/satellite/payments"
)

type consoleEndpoints struct {
	client     *http.Client
	base       string
	cookieName string
	email      string
	token      string
}

func newConsoleEndpoints(address string, email string) *consoleEndpoints {
	return &consoleEndpoints{
		client:     http.DefaultClient,
		base:       "http://" + address,
		cookieName: "_tokenKey",
		email:      email,
	}
}

func (ce *consoleEndpoints) login(ctx context.Context) (err error) {
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

func (ce *consoleEndpoints) checkLogin() error {
	if ce.token == "" {
		return errors.New("user is not yet logged in. please call login() first")
	}
	return nil
}

func (ce *consoleEndpoints) tryLogin(ctx context.Context, email string) (string, error) {
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
		ce.Token(),
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

	var token string
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return "", errs.Wrap(err)
	}

	return token, nil
}

func (ce *consoleEndpoints) tryCreateAndActivateUser(ctx context.Context, email string) error {
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

func (ce *consoleEndpoints) createRegistrationToken(ctx context.Context) (string, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		ce.RegToken(),
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

func (ce *consoleEndpoints) createUser(ctx context.Context, regToken string, email string) (string, error) {
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
		ce.Register(),
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

func (ce *consoleEndpoints) activateUser(ctx context.Context, userID string, email string) error {
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
		ce.Activation(activationToken),
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

func (ce *consoleEndpoints) setupAccount(ctx context.Context, token string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ce.SetupAccount(),
		nil)
	if err != nil {
		return err
	}

	request.AddCookie(&http.Cookie{
		Name:  ce.cookieName,
		Value: token,
	})

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

func (ce *consoleEndpoints) addCreditCard(ctx context.Context, token, cctoken string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ce.CreditCards(),
		strings.NewReader(cctoken))
	if err != nil {
		return err
	}

	request.AddCookie(&http.Cookie{
		Name:  ce.cookieName,
		Value: token,
	})

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

func (ce *consoleEndpoints) listCreditCards(ctx context.Context, token string) ([]payments.CreditCard, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		ce.CreditCards(),
		nil)
	if err != nil {
		return nil, err
	}

	request.AddCookie(&http.Cookie{
		Name:  ce.cookieName,
		Value: token,
	})

	resp, err := ce.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { err = errs.Combine(err, resp.Body.Close()) }()

	if resp.StatusCode != http.StatusOK {
		return nil, errs.New("unexpected status code: %d (%q)",
			resp.StatusCode, tryReadLine(resp.Body))
	}

	var list []payments.CreditCard

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (ce *consoleEndpoints) makeCreditCardDefault(ctx context.Context, token, ccID string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		ce.CreditCards(),
		strings.NewReader(ccID))
	if err != nil {
		return err
	}

	request.AddCookie(&http.Cookie{
		Name:  ce.cookieName,
		Value: token,
	})

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

func (ce *consoleEndpoints) getOrCreateProject(ctx context.Context) (string, error) {
	projectID, err := ce.getProject(ctx)
	if err == nil {
		return projectID, nil
	}
	projectID, err = ce.createProject(ctx)
	if err == nil {
		return projectID, nil
	}
	return ce.getProject(ctx)
}

func (ce *consoleEndpoints) getProject(ctx context.Context) (string, error) {
	query := `query {myProjects{id}}`
	var getProjects struct {
		MyProjects []struct {
			ID string
		}
	}
	err := ce.graphqlQuery(ctx, query, &getProjects)
	if len(getProjects.MyProjects) == 0 {
		return "", errs.New("No project exists")
	}
	return getProjects.MyProjects[0].ID, err
}

func (ce *consoleEndpoints) createProject(ctx context.Context) (string, error) {
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
	return createProject.CreateProject.ID, err
}

func (ce *consoleEndpoints) createAPIKey(ctx context.Context, projectID string) (string, error) {
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

func (ce *consoleEndpoints) graphqlQuery(ctx context.Context, createAPIKeyQuery string, response interface{}) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		ce.GraphQL(),
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

func (ce *consoleEndpoints) graphqlMutation(ctx context.Context, query string, response interface{}) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ce.GraphQL(),
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

func (ce *consoleEndpoints) graphqlDo(request *http.Request, jsonResponse interface{}) error {
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

func (ce *consoleEndpoints) appendPath(suffix string) string {
	return ce.base + suffix
}

func (ce *consoleEndpoints) RegToken() string {
	return ce.appendPath("/registrationToken/?projectsLimit=1")
}

func (ce *consoleEndpoints) Register() string {
	return ce.appendPath("/api/v0/auth/register")
}

func (ce *consoleEndpoints) SetupAccount() string {
	return ce.appendPath("/api/v0/payments/account")
}

func (ce *consoleEndpoints) CreditCards() string {
	return ce.appendPath("/api/v0/payments/cards")
}

func (ce *consoleEndpoints) Activation(token string) string {
	return ce.appendPath("/activation/?token=" + token)
}

func (ce *consoleEndpoints) Token() string {
	return ce.appendPath("/api/v0/auth/token")
}

func (ce *consoleEndpoints) GraphQL() string {
	return ce.appendPath("/api/v0/graphql")
}
