package passkey

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/redis/go-redis/v9"
	"github.com/soxft/openid-go/config"
	"github.com/soxft/openid-go/process/redisutil"
)

func registrationSessionKey(userID int) string {
	return fmt.Sprintf("%s:passkey:register:%d", config.RedisPrefix, userID)
}

func loginSessionKey(userID int) string {
	return fmt.Sprintf("%s:passkey:login:%d", config.RedisPrefix, userID)
}

func storeRegistrationSession(ctx context.Context, userID int, data *webauthn.SessionData, ttl time.Duration) error {
	return storeSession(ctx, registrationSessionKey(userID), data, ttl)
}

func loadRegistrationSession(ctx context.Context, userID int) (*webauthn.SessionData, error) {
	return loadSession(ctx, registrationSessionKey(userID))
}

func deleteRegistrationSession(ctx context.Context, userID int) error {
	return deleteSession(ctx, registrationSessionKey(userID))
}

func storeLoginSession(ctx context.Context, userID int, data *webauthn.SessionData, ttl time.Duration) error {
	return storeSession(ctx, loginSessionKey(userID), data, ttl)
}

func loadLoginSession(ctx context.Context, userID int) (*webauthn.SessionData, error) {
	return loadSession(ctx, loginSessionKey(userID))
}

func deleteLoginSession(ctx context.Context, userID int) error {
	return deleteSession(ctx, loginSessionKey(userID))
}

func storeSession(ctx context.Context, key string, data *webauthn.SessionData, ttl time.Duration) error {
	if redisutil.RDB == nil {
		return fmt.Errorf("redis not initialized")
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return redisutil.RDB.Set(ctx, key, payload, ttl).Err()
}

func loadSession(ctx context.Context, key string) (*webauthn.SessionData, error) {
	if redisutil.RDB == nil {
		return nil, fmt.Errorf("redis not initialized")
	}
	raw, err := redisutil.RDB.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	var session webauthn.SessionData
	if err := json.Unmarshal(raw, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func deleteSession(ctx context.Context, key string) error {
	if redisutil.RDB == nil {
		return fmt.Errorf("redis not initialized")
	}
	return redisutil.RDB.Del(ctx, key).Err()
}
