package zissuefs

import (
	"circuit/kit/zookeeper"
	"circuit/kit/zookeeper/zutil"
)

func Dial(zookeepers []string) *zookeeper.Conn {
	c, err := zutil.DialUntilReady(zutil.ZookeeperString(zookeepers))
	if err != nil {
		panic(err)
	}
	return c
}
