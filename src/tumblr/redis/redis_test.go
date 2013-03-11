package redis

import (
	"fmt"
	"testing"
)

func TestLowLevel(t *testing.T) {
	c, err := Dial("localhost:6300")
	if err != nil {
		fmt.Printf("err (%s)\n", err)
		return
	}
	err = c.WriteMultiBulk("get", "chris")
	if err != nil {
		fmt.Printf("err2 (%s)\n", err)
		return
	}
	resp, err := c.ReadResponse()
	if err != nil {
		fmt.Printf("read resp (%s)\n", err)
		return
	}
	fmt.Println(ResponseString(resp))
}

func TestSetGet(t *testing.T) {
	c, err := Dial("")
	if err != nil {
		fmt.Printf("err (%s)\n", err)
		return
	}
	if err = c.SetInt("oOOo", 345); err != nil {
		t.Fatalf("set (%s)", err)
	}
	i, err := c.GetInt("oOOo")
	if err != nil {
		t.Fatalf("get (%s)", err)
	}
	if i != 345 {
		t.Errorf("mismatch")
	}
}
