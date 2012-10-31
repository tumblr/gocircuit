package issuefs

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"tumblr/circuit/kit/join"
	"tumblr/circuit/use/lang"
)

type Issue struct {
	ID       int64
	Time     time.Time
	Reporter lang.Addr
	Affected lang.Addr
	Anchor   []string
	Msg      string
}

func ChooseID() int64 {
	return rand.Int63()
}

func IDString(id int64) string {
	return strconv.FormatInt(id, 10)
}

func ParseID(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func (i *Issue) String() string {
	if i == nil {
		return "nil issue"
	}
	var w bytes.Buffer
	fmt.Fprintf(&w, "MSG:      %s\n", i.Msg)
	fmt.Fprintf(&w, "ID:       %d\n", i.ID)
	fmt.Fprintf(&w, "Time:     %s\n", i.Time.Format(time.RFC1123))
	fmt.Fprintf(&w, "Reporter: %s\n", i.Reporter.String())
	if i.Affected != nil {
		fmt.Fprintf(&w, "Affected: %s\n", i.Affected.String())
	}
	fmt.Fprintf(&w, "Anchor:   ")
	for _, a := range i.Anchor {
		w.WriteString(a)
		w.WriteString(", ")
	}
	w.WriteByte('\n')
	return string(w.Bytes())
}

type fs interface {
	Add(msg string) int64
	Resolve(id int64) error
	List() []*Issue
	Subscribe(email string) error
	Unsubscribe(email string) error
}

// Bindings
var link = join.SetThenGet{Name: "issue file system"}

func Bind(v fs) {
	link.Set(v)
}

func get() fs {
	return link.Get().(fs)
}

func Add(msg string) int64 {
	return get().Add(msg)
}

func List() []*Issue {
	return get().List()
}

func Resolve(id int64) error {
	return get().Resolve(id)
}

func Subscribe(email string) error {
	return get().Subscribe(email)
}

func Unsubscribe(email string) error {
	return get().Unsubscribe(email)
}
