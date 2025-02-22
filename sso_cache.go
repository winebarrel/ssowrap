package ssowrap

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type SSOCache struct {
	*Options
	SSOStartUrl string
}

func NewSSOCache(options *Options, ssoStartUrl string) *SSOCache {
	ssoCache := &SSOCache{
		Options:     options,
		SSOStartUrl: ssoStartUrl,
	}

	return ssoCache
}

type Token struct {
	StartUrl    string
	AccessToken string
	Region      string
	ExpiresAt   time.Time
}

func (ssoCache *SSOCache) getTokens() ([]Token, error) {
	pat := filepath.Join(ssoCache.AWSSSOCacheDir, "*.json")
	files, err := filepath.Glob(pat)

	if err != nil {
		return nil, err
	}

	caches := []Token{}

	for _, f := range files {
		raw, err := os.ReadFile(f)

		if err != nil {
			return nil, err
		}

		var token Token
		err = json.Unmarshal(raw, &token)

		if err != nil {
			panic(err)
		}

		if token.StartUrl == "" || token.AccessToken == "" || token.Region == "" || token.ExpiresAt.IsZero() {
			continue
		}

		caches = append(caches, token)
	}

	return caches, nil
}

func (ssoCache *SSOCache) LastToken() (*Token, error) {
	caches, err := ssoCache.getTokens()

	if err != nil {
		return nil, err
	}

	var lastToken *Token
	maxExpiresAt := time.Now()

	for _, token := range caches {
		if token.StartUrl != ssoCache.SSOStartUrl {
			continue
		}

		if token.ExpiresAt.Before(maxExpiresAt) {
			continue
		}

		lastToken = &token
		maxExpiresAt = lastToken.ExpiresAt
	}

	if lastToken == nil {
		return nil, fmt.Errorf("SSO token not found, try `aws sso login`")
	}

	return lastToken, nil
}
