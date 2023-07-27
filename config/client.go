package config

import "time"

const (
	PONG_WAIT     = 10 * time.Second
	PING_INTERVAL = (PONG_WAIT * 9) / 10
)
