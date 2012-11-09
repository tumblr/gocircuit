// Package util implements commonly used Zookeeper patterns and recipes
package zutil

import (
	"errors"
	"path"
	"strings"
	"circuit/kit/zookeeper"
)

func CreateRecursive(z *zookeeper.Conn, leafPath string, aclv []zookeeper.ACL) error {
	leafPath = path.Clean(leafPath)
	if len(leafPath) == 0 || leafPath[0] != '/' {
		return errors.New("zookeeper util path syntax")
	}
	parts := strings.Split(leafPath, "/")
	if len(parts) < 1 {
		return errors.New("creating zookeeper root")
	}
	prefix := "/"
	for i := 0; i < len(parts); i++ {
		prefix = path.Join(prefix, parts[i])
		if _, err := z.Create(prefix, "", 0, aclv); err != nil && !IsNodeExists(err) {
			return err
		}
	}
	return nil
}
