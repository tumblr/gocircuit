// Copyright 2012 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package util implements commonly used Zookeeper patterns and recipes
package zutil

import (
	"circuit/kit/zookeeper"
	"errors"
	"path"
	"strings"
)

// CreateRecursive creates the directory leafPath and any parent subdirectories if necessary
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
