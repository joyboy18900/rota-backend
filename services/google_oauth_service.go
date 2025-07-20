package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuthService interface {
	GetAuthURL(state string) string
	GetAuthURLWithRedirect(state, redirectURI string) string
	ExchangeCode(code, redirectURI string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*GoogleUserInfo, error)
	GenerateState() string
	ValidateState(state string) bool
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type stateStore struct {
	mu     sync.RWMutex
	states map[string]time.Time
}

func (s *stateStore) add(state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[state] = time.Now().Add(10 * time.Minute)
}

func (s *stateStore) validate(state string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	expiry, exists := s.states[state]
	if !exists {
		return false
	}
	return time.Now().Before(expiry)
}

func (s *stateStore) remove(state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, state)
}

type googleOAuthService struct {
	config *oauth2.Config
	states *stateStore
}

func NewGoogleOAuthService() GoogleOAuthService {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	stateStore := &stateStore{
		states: make(map[string]time.Time),
	}
	
	service := &googleOAuthService{
		config: config,
		states: stateStore,
	}
	
	go service.cleanupExpiredStates()
	
	return service
}

func (s *googleOAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("redirect_uri", s.config.RedirectURL))
}

func (s *googleOAuthService) GetAuthURLWithRedirect(state, redirectURI string) string {
	// Create a temporary config with the custom redirect URI
	tempConfig := *s.config
	tempConfig.RedirectURL = redirectURI
	return tempConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *googleOAuthService) ExchangeCode(code, redirectURI string) (*oauth2.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Create a temporary config with the custom redirect URI
	tempConfig := *s.config
	tempConfig.RedirectURL = redirectURI
	
	return tempConfig.Exchange(ctx, code)
}

func (s *googleOAuthService) GetUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client := s.config.Client(ctx, token)
	
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}
	
	return &userInfo, nil
}

func (s *googleOAuthService) GenerateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	s.states.add(state)
	return state
}

func (s *googleOAuthService) ValidateState(state string) bool {
	if s.states.validate(state) {
		s.states.remove(state)
		return true
	}
	return false
}

func (s *googleOAuthService) cleanupExpiredStates() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			s.states.mu.Lock()
			now := time.Now()
			for state, expiry := range s.states.states {
				if now.After(expiry) {
					delete(s.states.states, state)
				}
			}
			s.states.mu.Unlock()
		}
	}
}
