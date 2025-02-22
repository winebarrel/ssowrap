package ssowrap

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type SSOTokenCache struct {
	*Options
	SSOStartUrl string
}

func NewSSOTokenCache(options *Options, ssoStartUrl string) *SSOTokenCache {
	ssoTokenCache := &SSOTokenCache{
		Options:     options,
		SSOStartUrl: ssoStartUrl,
	}

	return ssoTokenCache
}

type Token struct {
	StartUrl    string
	AccessToken string
	Region      string
	ExpiresAt   time.Time
}

func (ssoTokenCache *SSOTokenCache) getCaches() ([]Token, error) {
	pat := filepath.Join(ssoTokenCache.AWSSSOCacheDir, "*.json")
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

func (ssoTokenCache *SSOTokenCache) LastToken() (*Token, error) {
	caches, err := ssoTokenCache.getCaches()

	if err != nil {
		return nil, err
	}

	var lastToken *Token
	maxExpiresAt := time.Now()

	for _, token := range caches {
		if token.ExpiresAt.Before(maxExpiresAt) {
			continue
		}

		if token.StartUrl != ssoTokenCache.SSOStartUrl {
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
