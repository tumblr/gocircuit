package redis

import (
	"errors"
	"strconv"
)

var (
	ErrNotInt64 = errors.New("type did not assert against int64")
)

// increment a redis key by 1
func (r *Conn) Incr(key string) (int64, error) {
	id := r.Next()
	r.StartRequest(id)
	err := r.WriteMultiBulk("INCR", key)
	r.EndRequest(id)

	if err != nil {
		return 0, err
	}

	return r.readResponseInt64(id)
}

// decrement a redis key by 1
func (r *Conn) Decr(key string) (int64, error) {
	id := r.Next()

	r.StartRequest(id)
	err := r.WriteMultiBulk("DECR", key)
	r.EndRequest(id)

	if err != nil {
		return 0, err
	}

	return r.readResponseInt64(id)
}

func (r *Conn) SetInt(key string, value int64) error {
	id := r.Next()

	r.StartRequest(id)
	err := r.WriteMultiBulk("SET", key, strconv.FormatInt(value, 10))
	r.EndRequest(id)

	if err != nil {
		return err
	}

	return r.readResponse(id)
}

type KeyIntValue struct {
	Key   string
	Value int64
}

/*
func (r *Conn) SetIntBulk(pairs ...KeyInt64Value) error {
	?
}
*/

// get the value of a redis key as an
func (r *Conn) GetInt(key string) (int64, error) {
	id := r.Next()

	r.StartRequest(id)
	err := r.WriteMultiBulk("GET", key)
	r.EndRequest(id)

	if err != nil {
		return 0, err
	}

	return r.readResponseInt64(id)
}

func (r *Conn) readResponse(id uint) error {
	r.StartResponse(id)
	_, err := r.ReadResponse()
	r.EndResponse(id)
	return err
}

func (r *Conn) readResponseInt64(id uint) (int64, error) {
	r.StartResponse(id)
	resp, err := r.ReadResponse()
	r.EndResponse(id)
	if err != nil {
		return 0, err
	}

	respString, ok := resp.(Bulk)
	if !ok {
		return 0, errors.New("unknown remote response")
	}
	return strconv.ParseInt(string(respString), 10, 64)
}
