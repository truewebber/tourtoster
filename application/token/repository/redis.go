package repository

import (
	"github.com/go-redis/redis"

	"tourtoster/token"
)

type (
	Redis struct {
		conn *redis.Client
	}
)

func NewRedis(addr string, db int) (*Redis, error) {
	conn := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
	if err := conn.Ping().Err(); err != nil {
		return nil, err
	}

	return &Redis{
		conn: conn,
	}, nil
}

func (r *Redis) Token(token string) (*token.Token, error) {
	return nil, nil
}

func (r *Redis) Save(token *token.Token) error {
	return nil
}

func (r *Redis) Delete(token string) error {
	return nil
}
