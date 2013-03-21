package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"tumblr/firehose"
	"tumblr/balkan/proto"
)

func StreamFirehose(freq *firehose.Request) <-chan *createRequest {
	ch := make(chan *createRequest)
	go func() {
		conn := firehose.Redial(freq)
		for {
			q := filter(conn.Read())
			if q == nil {
				continue
			}
			println(fmt.Sprintf("CREATE blogID=%d postID=%d", q.TimelineID, q.PostID))
			ch <- &createRequest{
				Forwarded:      false,
				Post:           q,
				ReturnResponse: func(err error) {
					if err != nil {
						println("Firehose->XCreatePost error:", err.Error())
						return
					}
				},
			}
		}
	}()
	return ch
}

func filter(e *firehose.Event) *proto.XCreatePost {
	if e.Activity != firehose.CreatePost {
		return nil
	}
	return &proto.XCreatePost{TimelineID: e.Post.BlogID, PostID: e.Post.ID}
}

func sendmail(recipient, subject, body string) error {
	cmd := exec.Command("sendmail", recipient)
	var w bytes.Buffer
	w.WriteString("Subject: ")
	w.WriteString(subject)
	w.WriteByte('\n')
	w.Write([]byte(body))
	cmd.Stdin = &w
	_, err := cmd.CombinedOutput()
	return err
}
