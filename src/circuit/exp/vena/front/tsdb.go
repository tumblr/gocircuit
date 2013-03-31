package front

import (
	"circuit/kit/sched/limiter"
	"net"
	"strconv"
	"strings"
)

type Replier interface {
	Reply(??)
}

func listenTSDB(addr string, reply Replier) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	// Accept incoming requests
	go func() {
		lmtr := l.New(100)	// At most 100 concurrent connections
		for {
			lmtr.Open()
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}
			// Serve individual connection
			go func() {
				defer lmtr.Close()
				defer conn.Close()
				// Read request, send reply
				r := bufio.NewReader(conn)
				for {
					line, err := bufio.ReadString("\n")
					if err != nil {
						println("read line", err.Error())
						break
					}
					cmd, err := parse(line)
					if err != nil {
						println("parse", err.Error())
						break
					}
					if cmd == nil {
						continue
					}
					switch cmd.(type) {
					case quit:
						??
					case put:
						??
					}
					reply.Reply()
				}
			}()
		}
	}()
}

type quit struct{}

type put struct{
	Metric string
	Time   vena.Time
	Tags   []*struct{Tag, Value string}
	Value  float64
}

// put proc.loadavg.1min 1234567890 1.35 host=A
func parse(l string) (interface{}, error) {
	t := strings.Split(l, " ")
	if len(t) == 0 {
		return nil, nil
	}
	if t[0] == "diediedie" {
		return quit{}, nil
	}
	if t[0] != "put" {
		return nil, errors.New("unrecognized command")
	}
	t = t[1:]
	if len(t) < 3 {
		return nil, errors.New("too few")
	}
	a := &put{Metric: t[0]}
	// Time
	sec, err := strconv.Atoi(t[1])
	if err != nil {
		return nil, err
	}
	a.Time = vena.Time(sec)
	// Value
	a.Value, err = strconv.ParseFloat(t[2], 64)
	if err != nil {
		return nil, err
	}
	t = t[3:]
	// Tags
	for _, tv := range t {
		q := strings.SplitN(tv, ":", 2)
		a.Tags = append(a.Tags, &struct{Tag, Value string}{Tag: , Value: })
	}
	return a, nil
}
