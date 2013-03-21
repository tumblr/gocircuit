package kafka

import (
	"fmt"
	"testing"
)

func TestClientConn(t *testing.T) {
	c, err := Dial("192.168.128.121:9092")
	if err != nil {
		t.Fatalf("connect to broker (%s)", err)
	}

	// Produce
	fmt.Printf("Producing\n")
	err = c.Produce(&ProduceArg{
		Topic:     "_hello",
		Partition: 0,
		Messages:  [][]byte{{1, 2, 3}, {4, 5, 6}},
	})
	if err != nil {
		t.Errorf("produce (%s)", err)
	}

	// Fetch
	fmt.Printf("Fetching\n")
	returns, err := c.Fetch(&FetchArg{
		Topic:     "_hello",
		Partition: 0,
		Offset:    0,
		MaxSize:   100, // Intentionally an off-boundary max size
	})
	if err != nil {
		t.Errorf("fetch (%s)", err)
	}
	fmt.Printf("%v\n", returns)

	// Offsets
	fmt.Printf("Offsets\n")
	offsets, err := c.Offsets(&OffsetsArg{
		Topic:      "_hello",
		Partition:  0,
		Time:       Latest,
		MaxOffsets: 13,
	})
	if err != nil {
		t.Errorf("offsets (%s)", err)
	}
	if len(offsets) != 2 || offsets[1] != 0 {
		t.Fatalf("unexpected offsets")
	}
	fmt.Printf("%v\n", offsets)

	// Multi-produce
	fmt.Printf("Multi-producing\n")
	err = c.Produce(
		&ProduceArg{
			Topic:     "_hello",
			Partition: 0,
			Messages:  [][]byte{{7, 8, 9}, {1, 1, 1}},
		},
		&ProduceArg{
			Topic:     "_hello",
			Partition: 0,
			Messages:  [][]byte{{3, 3, 3}, {2, 2, 2}},
		},
	)
	if err != nil {
		t.Errorf("multi-produce (%s)", err)
	}

	// Multi-fetch
	fmt.Printf("Multi-fetching\n")
	returns, err = c.Fetch(
		&FetchArg{
			Topic:     "_hello",
			Partition: 0,
			Offset:    offsets[0],
			MaxSize:   100, // Intentionally an off-boundary max size
		},
		&FetchArg{
			Topic:     "_hello",
			Partition: 0,
			Offset:    offsets[1],
			MaxSize:   100, // Intentionally an off-boundary max size
		},
	)
	if err != nil {
		t.Errorf("multi-fetch (%s)", err)
	}
	fmt.Printf("%v\n", returns)

	// Close
	fmt.Printf("Closing\n")
	err = c.Close()
	if err != nil {
		t.Errorf("connect to broker (%s)", err)
	}
}
