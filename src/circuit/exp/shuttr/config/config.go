package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"tumblr/firehose"
	"tumblr/balkan/shard"
	"tumblr/balkan/x"
	"circuit/kit/xor"
)

type Config struct {
	InstallDir string
	Firehose   *firehose.Request	// Firehose credentials
	Timeline   []*shard.Shard
	Dashboard  []*shard.Shard
	PushMap    string		// File name of push map
}

type configSource struct {
	InstallDir string
	Firehose   *firehose.Request
	Timeline   []*shardSource
	Dashboard  []*shardSource
	PushMap    string
}

type shardSource struct {
	Pivot string
	Addr  string
	HTTP  int
}

func Read(name string) (*Config, error) {
	raw, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	csrc := &configSource{}
	if err = json.Unmarshal(raw, csrc); err != nil {
		return nil, err
	}

	c := &Config{
		InstallDir: csrc.InstallDir,
		Firehose:   csrc.Firehose,
		PushMap:    csrc.PushMap,
	}

	if c.Timeline, err = makeShards(csrc.Timeline); err != nil {
		return nil, err
	}
	if c.Dashboard, err = makeShards(csrc.Dashboard); err != nil {
		return nil, err
	}

	return c, nil
}

func makeShards(src []*shardSource) ([]*shard.Shard, error) {
	out := make([]*shard.Shard, len(src))
	for i, sh := range src {
		if strings.Index(sh.Pivot, "0x") != 0 {
			return nil, errors.New("invalid pivot format")
		}
		pivot, err := strconv.ParseUint(sh.Pivot[2:], 16, 64)
		if err != nil {
			return nil, err
		}
		out[i] = &shard.Shard{
			Pivot: xor.Key(pivot),
			Addr:  x.Addr(sh.Addr),
			HTTP: sh.HTTP,
		}
	}
	return out, nil
}
