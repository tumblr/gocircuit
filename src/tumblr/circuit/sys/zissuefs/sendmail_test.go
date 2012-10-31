package zissuefs

import (
	"testing"
)

func TestSendmail(t *testing.T) {
	if err := sendmail("user@test.com", "hi subject", "some body\n\naha\nx\n"); err != nil {
		t.Fatalf("sendmail (%s)", err)
	}
}
