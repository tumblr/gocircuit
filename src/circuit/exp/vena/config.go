package vena

import (
	"circuit/kit/xor"
	"path"
)

type Config struct {
	Shard  []xor.Key
	Anchor string // Root anchor for the shards
}

func (c *Config) ShardAnchor(key xor.Key) string {
	return path.Join(c.Anchor, key.String())
}
