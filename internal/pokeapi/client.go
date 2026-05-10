package pokeapi

import (
	"github.com/kelvinjrosado/pokedex/internal/logger"
	"github.com/kelvinjrosado/pokedex/internal/pokecache"
)

type Client struct {
	Cache  *pokecache.Cache
	Logger *logger.CustomLogger
}

func NewClient(cache *pokecache.Cache, lgr *logger.CustomLogger) *Client {
	return &Client{Cache: cache, Logger: lgr}
}
