package api

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"strings"
	"time"
)

type Credentials struct {
	KeyID            string `json:"key_id"`
	ServiceAccountID string `json:"service_account_id"`
	PrivateKey       string `json:"private_key"`
}

type Token struct {
	IAMToken  string `json:"iamToken"`
	ExpiresAt string `json:"expiresAt"`
}

func GetIAMToken(credentials Credentials) (Token, error) {
	var token Token
	jot, err := signedToken(credentials)
	if err != nil {
		return token, err
	}
	resp, err := http.Post(
		"https://iam.api.cloud.yandex.net/iam/v1/tokens",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"jwt":"%s"}`, jot)),
	)
	if err != nil {
		return token, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return token, fmt.Errorf("%s: %s", resp.Status, body)
	}

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return token, err
	}

	return token, nil
}

// signedToken формирование JWT
func signedToken(credentials Credentials) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    credentials.ServiceAccountID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Audience:  []string{"https://iam.api.cloud.yandex.net/iam/v1/tokens"},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["kid"] = credentials.KeyID

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(credentials.PrivateKey))
	if err != nil {
		return "", err
	}
	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

type InstancesList struct {
	NextPageToken string     `json:"nextPageToken"`
	Instances     []Instance `json:"instances"`
}

type Instance struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type InstancesListOptions struct {
	FolderId  string `json:"folderId"`
	PageSize  int    `json:"pageSize"`
	PageToken string `json:"pageToken"`
	Filter    string `json:"filter"`
}
