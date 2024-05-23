package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionManager struct {
	rdb        *redis.Client
	expiration int
}

type UserSession struct {
	UserId string
}

func (s *SessionManager) Close() error {
	return s.rdb.Close()
}

func NewSessionManager(
	addr string,
	password string,
	expiration int,
) (*SessionManager, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	return &SessionManager{
		rdb:        rdb,
		expiration: expiration,
	}, nil
}

func (s *SessionManager) NewSession(
	ctx context.Context,
	userId int32,
) (string, error) {
	key := uuid.NewString()
	status := s.rdb.Set(
		ctx,
		key,
		userId,
		time.Duration(time.Minute*time.Duration(s.expiration)),
	)
	if status.Err() != nil {
		return "", status.Err()
	}

	return key, nil
}

func (s *SessionManager) UserIdForSession(
	ctx context.Context,
	sessionKey string,
) (string, error) {
	res := s.rdb.Get(ctx, sessionKey)
	if res.Err() != nil {
		return "", res.Err()
	}
	return res.Val(), nil
}
