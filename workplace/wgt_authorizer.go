// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Lifted from https://github.com/hashicorp/go-azure-sdk/blob/f6aaaa9a/sdk/auth/azure_cli_authorizer.go

package workplace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"golang.org/x/oauth2"
)

// NewWgtAuthorizer returns an Authorizer which authenticates using tools/wpGetToken
func NewWgtAuthorizer(ctx context.Context) (auth.Authorizer, error) {
	// Cache access tokens internally to avoid unnecessary `az` invocations
	return auth.NewCachedAuthorizer(&WgtAuthorizer{})
}

var _ auth.Authorizer = &WgtAuthorizer{}

// WgtAuthorizer is an Authorizer which supports tools/wpGetToken
type WgtAuthorizer struct {
}

// Token returns an access token using wpGetToken as an authentication mechanism.
func (a *WgtAuthorizer) Token(_ context.Context, _ *http.Request) (*oauth2.Token, error) {
	var tokenResult wgtTokenResult
	if err := jsonUnmarshalWgtCmd(&tokenResult); err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken: tokenResult.AccessToken,
		TokenType:   tokenResult.TokenType,
	}, nil
}

// AuxiliaryTokens returns additional tokens for auxiliary tenant IDs, for use in multi-tenant scenarios
func (a *WgtAuthorizer) AuxiliaryTokens(_ context.Context, _ *http.Request) ([]*oauth2.Token, error) {
	panic("WgtAuthorizer does not support method AuxiliaryTokens()")
}

type wgtTokenResult struct {
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	ExpiresOn   int64  `json:"expires_on"`
}

// jsonUnmarshalWgtCmd executes wpGetToken and unmarshalls the JSON output.
func jsonUnmarshalWgtCmd(i interface{}, arg ...string) error {
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	cmd := exec.Command("wpGetToken", arg...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Start(); err != nil {
		err := fmt.Errorf("launching wpGetToken: %+v", err)
		if stdErrStr := stderr.String(); stdErrStr != "" {
			err = fmt.Errorf("%s: %s", err, strings.TrimSpace(stdErrStr))
		}
		return err
	}

	if err := cmd.Wait(); err != nil {
		err := fmt.Errorf("running wpGetToken: %+v", err)
		if stdErrStr := stderr.String(); stdErrStr != "" {
			err = fmt.Errorf("%s: %s", err, strings.TrimSpace(stdErrStr))
		}
		return err
	}

	if err := json.Unmarshal(stdout.Bytes(), &i); err != nil {
		return fmt.Errorf("unmarshaling the output of wpGetToken: %v", err)
	}

	return nil
}
