package main

import (
	"github.com/omarhachach/bear"
)

// Config is extended from Bear config.
type Config struct {
	*bear.Config
	SupportChannelID string `json:"support_channel_id"`
	DB               string `json:"db"`
}
