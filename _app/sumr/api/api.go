package api

import (
	"circuit/use/circuit"
	"circuit/kit/sched/limiter"
	"tumblr/app/sumr/client"
)

var (
	ErrMode      = circuit.NewError("write operation on read-only API")
	ErrBackend   = circuit.NewError("backend")
	ErrFormat    = circuit.NewError("format")
	ErrFields    = circuit.NewError("bad fields")
	ErrNoValue   = circuit.NewError("missing value")
	ErrNoFeature = circuit.NewError("missing feature")
	ErrFieldType = circuit.NewError("field type not string")
	ErrTime      = circuit.NewError("time format")
)

type API struct {
	server    *httpServer
	client    *client.Client
	lmtr      limiter.Limiter
}

func init() {
	circuit.RegisterType(&API{}) // Register as circuit value
}

func New(dfile string, port int, readOnly bool) (api *API, err error) {
	api = &API{}
	api.client, err = client.New(dfile, readOnly)
	if err != nil {
		return nil, err
	}
	api.lmtr.Init(200)
	api.server, err = startServer(
		port,
		func(req []interface{}) []interface{} { 
			return api.respondAdd(req)
		},
		func(req []interface{}) []interface{} { 
			return api.respondSum(req)
		},
	)
	return api, err
}

// Given slice of AddRequests, fire a batch query to client and fetch responses as slice of Response
// respondAdd will panic if the underlying SUMR client panics.
func (api *API) respondAdd(req []interface{}) []interface{} {
	api.lmtr.Open()
	defer api.lmtr.Close()

	q := make([]client.AddRequest, len(req))
	for i, a_ := range req {
		a := a_.(*AddRequest)
		q[i].UpdateTime = a.Change.Time
		q[i].Key = a.Key()
		q[i].Value = a.Change.Value
	}
	r := api.client.AddBatch(q)
	s := make([]interface{}, len(req))
	for i, _ := range s {
		s[i] = &Response{Sum: r[i]}
	}
	return s
}

func (api *API) respondSum(req []interface{}) []interface{} {
	api.lmtr.Open()
	defer api.lmtr.Close()

	q := make([]client.SumRequest, len(req))
	for i, a_ := range req {
		a := a_.(*SumRequest)
		q[i].Key = a.Key()
	}
	r := api.client.SumBatch(q)
	s := make([]interface{}, len(req))
	for i, _ := range s {
		s[i] = &Response{Sum: r[i]}
	}
	return s
}
