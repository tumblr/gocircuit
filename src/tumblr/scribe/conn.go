package scribe

import (
	"errors"
	"net"
	"sync"
	"tumblr/encoding/thrift"
	"tumblr/net/scribe/thrift/scribe"
)

type Conn struct {
	sync.Mutex
	transport thrift.TTransport
	client    *scribe.ScribeClient
}

func Dial(hostport string) (*Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp", hostport)
	if err != nil {
		return nil, err
	}
	conn := &Conn{}
	if conn.transport, err = thrift.NewTNonblockingSocketAddr(addr); err != nil {
		return nil, err
	}
	conn.transport = thrift.NewTFramedTransport(conn.transport)
	protocol := thrift.NewTBinaryProtocolFactoryDefault()
	conn.client = scribe.NewScribeClientFactory(conn.transport, protocol)
	if err = conn.transport.Open(); err != nil {
		return nil, err
	}
	return conn, nil
}

type Message struct {
	Category string
	Payload  string
}

func (conn *Conn) E(category, payload string) error {
	return conn.Emit(Message{category, payload})
}

func (conn *Conn) Emit(msgs ...Message) error {
	tlist := thrift.NewTList(thrift.TypeFromValue(scribe.NewLogEntry()), len(msgs))
	for _, msg := range msgs {
		tlog := scribe.NewLogEntry()
		tlog.Category = msg.Category
		tlog.Message = msg.Payload
		tlist.Push(tlog)
	}

	conn.Lock()
	defer conn.Unlock()

	result, err := conn.client.Log(tlist)
	if err != nil {
		return err
	}
	err = resultCodeToError(result)
	if err != nil {
		return err
	}
	return nil
}

var (
	ErrUnknown = errors.New("thrift unknown result code")
	ErrTryLater = errors.New("thrift try later")
)

func resultCodeToError(resultCode scribe.ResultCode) error {
	switch resultCode {
	case scribe.OK:
		return nil
	case scribe.TRY_LATER:
		return ErrTryLater
	}
	return ErrUnknown
}

func (conn *Conn) Close() error {
	conn.Lock()
	defer conn.Unlock()
	return conn.transport.Close()
}
