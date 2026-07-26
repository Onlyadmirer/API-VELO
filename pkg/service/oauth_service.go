package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/helper"
	"VELO-backend/pkg/repository"

	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	ID      string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type OAuthService interface {
	GetGoogleLoginURL() (string, string, error)
	HandleGoogleCallback(code string, state string) (*http.Cookie, *entity.User, error)
}

type oauthService struct {
	repo  repository.UserRepository
	redis *redis.Client
	oauth *oauth2.Config
}

func NewOAuthService(repo repository.UserRepository, redisClient *redis.Client) OAuthService {
	return &oauthService{
		repo:  repo,
		redis: redisClient,
		oauth: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"openid",
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
	}
}

func (s *oauthService) generateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *oauthService) GetGoogleLoginURL() (string, string, error) {
	state, err := s.generateState()
	if err != nil {
		return "", "", err
	}

	err = s.redis.Set(context.Background(), "oauth_state:"+state, "1", 10*time.Minute).Err()
	if err != nil {
		return "", "", err
	}

	url := s.oauth.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url, state, nil
}

func (s *oauthService) verifyState(state string) error {
	key := "oauth_state:" + state
	_, err := s.redis.Get(context.Background(), key).Result()
	if err != nil {
		return errors.New("invalid state parameter")
	}
	s.redis.Del(context.Background(), key)
	return nil
}

func (s *oauthService) HandleGoogleCallback(code string, state string) (*http.Cookie, *entity.User, error) {
	if err := s.verifyState(state); err != nil {
		return nil, nil, err
	}

	token, err := s.oauth.Exchange(context.Background(), code)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal exchange code: %v", err)
	}

	userInfo, err := s.getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return nil, nil, err
	}

	existingUser, err := s.repo.GetUserByOAuth("google", userInfo.ID)
	if err != nil {
		return nil, nil, err
	}

	var user *entity.User

	if existingUser != nil {
		user = existingUser
	} else {
		existingByEmail, _ := s.repo.GetUserByEmail(userInfo.Email)
		if existingByEmail != nil {
			err = s.repo.LinkOAuthToUser(existingByEmail.ID, "google", userInfo.ID)
			if err != nil {
				return nil, nil, err
			}
			user, err = s.repo.GetUserByID(existingByEmail.ID)
			if err != nil {
				return nil, nil, err
			}
		} else {
			newUser := &entity.User{
				Name:          userInfo.Name,
				Email:         userInfo.Email,
				OAuthID:       &userInfo.ID,
				OAuthProvider: strPtr("google"),
				AvatarURL:     &userInfo.Picture,
				Role:          "customer",
				IsVerified:    true,
			}
			user, err = s.repo.CreateOAuthUser(newUser)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	jwtToken, err := helper.GenerateJWTToken(user.ID, user.Role)
	if err != nil {
		return nil, nil, err
	}

	cookie := &http.Cookie{
		Name:     "jwt_token",
		Value:    jwtToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}

	return cookie, user, nil
}

func (s *oauthService) getGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("gagal get user info: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal baca response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google API error: %s", string(body))
	}

	var info GoogleUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("gagal parse user info: %v", err)
	}

	return &info, nil
}

func strPtr(s string) *string {
	return &s
}
