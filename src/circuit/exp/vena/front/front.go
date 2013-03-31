// Copyright 2013 Tumblr, Inc.
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

package front

import (
	"circuit/exp/vena"
	"circuit/exp/vena/client"
	"circuit/kit/sched/limiter"
	"circuit/use/circuit"
)

var (
	ErrMode      = circuit.NewError("write operation on read-only API")
)

type Front struct {
	http   *httpServer
	tsdb   *tsdbServer
	client *client.Client
	lmtr   limiter.Limiter
}

func init() {
	circuit.RegisterValue(&Front{})
}

func New(c *vena.Config, httpPort, tsdbPort int) *Front {
	front := &Front{}
	front.client, err = client.New(c)
	if err != nil {
		panic(err)
	}
	front.lmtr.Init(200)
	front.http, err = startHTTP(httpPort)
	if err != nil {
		panic(err)
	}
	front.tsdb, err = startTSDB(tsdbPort)
	if err != nil {
		panic(err)
	}
	return front
}

// Given slice of AddRequests, fire a batch query to client and fetch responses as slice of Response
// respondAdd will panic if the underlying sumr client panics.
func (api *API) respondAdd(req []interface{}) []interface{} {
	api.lmtr.Open()
	defer api.lmtr.Close()

	q := make([]client.AddRequest, len(req))
	for i, a_ := range req {
		a := a_.(*addRequest)
		q[i].UpdateTime = a.change.Time
		q[i].Key = a.Key()
		q[i].Value = a.change.Value
	}
	r := api.client.AddBatch(q)
	s := make([]interface{}, len(req))
	for i, _ := range s {
		s[i] = &response{Sum: r[i]}
	}
	return s
}

func (api *API) respondSum(req []interface{}) []interface{} {
	api.lmtr.Open()
	defer api.lmtr.Close()

	q := make([]client.SumRequest, len(req))
	for i, a_ := range req {
		a := a_.(*sumRequest)
		q[i].Key = a.Key()
	}
	r := api.client.SumBatch(q)
	s := make([]interface{}, len(req))
	for i, _ := range s {
		s[i] = &response{Sum: r[i]}
	}
	return s
}
