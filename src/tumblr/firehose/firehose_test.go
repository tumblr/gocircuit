package firehose

import (
	"fmt"
	"testing"
)

var testFreq = &Request{
	HostPort:      "",
	Username:      "",
	Password:      "",
	ApplicationID: "",
	ClientID:      "",
	Offset:        "",
}

var testPrivFreq = &Request{
	HostPort:      "",
	Username:      "",
	Password:      "",
	ApplicationID: "",
	ClientID:      "",
	Offset:        "",
}

func validateRaw(s string) {
	bb := []byte(s)
	for _, b := range bb {
		if b == 0 {
			fmt.Printf("0-byte\n")
		}
	}
}

type fishOutActivity struct {
	Activity string `json:"activity"`
}

func TestActivity(t *testing.T) {
	conn, err := Dial(testFreq)
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}
	for {
		fa := &fishOutActivity{}
		if err := conn.ReadInterface(fa); err != nil {
			t.Errorf("read interface (%s)", err)
		} else {
			fmt.Printf("[%s]\n", fa.Activity)
		}
	}
	if err = conn.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}
}

func TestReadRaw(t *testing.T) {
	conn, err := Dial(testFreq)
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}
	for i := 0; i < 4; i++ {
		if line, err := conn.ReadRaw(); err != nil {
			t.Errorf("read raw (%s)", err)
		} else {
			validateRaw(line)
			fmt.Printf("`%s`\n———\n", line)
		}
	}
	if err = conn.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}
}

func TestReadEvent(t *testing.T) {
	conn, err := Dial(testFreq)
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}
	for i := 0; i < 100; i++ {
		if ev, err := conn.Read(); err != nil {
			t.Errorf("read (%s)", err)
		} else {
			//fmt.Printf("PrivateData: %#v\n", ev.PrivateData)
			//fmt.Printf("a=%s\n", ev.Activity.String())
			if ev.Post != nil {
				if ev.Post.BlogID == 0 {
					fmt.Printf("WOA\n")
				}
				//fmt.Printf("BlogID: %#v\n", ev.Post.BlogID)
				//fmt.Printf("BlogName=%s\n", ev.Post.BlogName)
				//fmt.Printf("Tags=%#v\n", ev.Post.Tags)
			}
		}
	}
	if err = conn.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}
}

func TestFirehose(t *testing.T) {
	conn, err := Dial(testFreq)
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}

	for i := 0; i < 20; i++ {
		if line, err := conn.ReadRaw(); err != nil {
			t.Errorf("read raw (%s)", err)
		} else {
			fmt.Printf("%s\n———\n", line)
		}
	}

	for i := 0; i < 20; i++ {
		v := make(map[string]interface{})
		if err = conn.ReadInterface(&v); err != nil {
			t.Errorf("read interface (%s)", err)
		} else {
			fmt.Printf("%v\n———\n", v)
		}
	}

	for i := 0; i < 20; i++ {
		if ev, err := conn.Read(); err != nil {
			t.Errorf("read (%s)", err)
		} else {
			fmt.Printf("%v\n———\n", ev)
		}
	}

	if err = conn.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}
}
