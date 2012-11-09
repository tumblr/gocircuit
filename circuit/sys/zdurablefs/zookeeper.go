package zdurablefs

import (
	"tumblr/circuit/kit/zookeeper"
	"tumblr/circuit/kit/zookeeper/zutil"
)

func Dial(zookeepers []string) *zookeeper.Conn {
	c, err := zutil.DialUntilReady(zutil.ZookeeperString(zookeepers))
	if err != nil {
		panic(err)
	}
	return c
}
