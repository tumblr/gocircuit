package front

import (
	"bufio"
	"circuit/exp/vena"
	"circuit/kit/sched/limiter"
	"errors"
	"net"
	"strconv"
	"strings"
)

type Replier interface {
	Put(string, vena.Time, []*tag, float64)
	Quit()
}

func listenTSDB(addr string, reply Replier) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	// Accept incoming requests
	go func() {
		lmtr := limiter.New(100)	// At most 100 concurrent connections
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
				defer recover()	// Recover from panics in reply logic
				// Read request, send reply
				r := bufio.NewReader(conn)
				for {
					line, err := r.ReadString('\n')
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
					switch p := cmd.(type) {
					case quit:
						reply.Quit()
					case *put:
						reply.Put(p.Metric, p.Time, p.Tags, p.Value)
					}
				}
			}()
		}
	}()
}

type quit struct{}

type tag struct {
	Name  string
	Value string
}

type put struct{
	Metric string
	Time   vena.Time
	Tags   []*tag
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
		if len(q) != 2 {
			return nil, errors.New("parse tag")
		}
		a.Tags = append(a.Tags, &tag{Name: q[0], Value: q[1]})
	}
	return a, nil
}
