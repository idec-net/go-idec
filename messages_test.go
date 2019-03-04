package idec

import (
	"testing"
)

func TestCollectTags(t *testing.T) {
	tags := Tags{
		II:    "ok",
		Repto: "hXzRNEzmMuzKkT1HCxUb",
	}
	collected, err := tags.CollectTags()
	if err != nil {
		t.Error(err)
	}
	// With repto
	if collected != "ii/ok/repto/hXzRNEzmMuzKkT1HCxUb" {
		t.Error("Wrong tags collection")
	}
	// Without repto
	tags.Repto = ""
	collected, err = tags.CollectTags()
	if collected != "ii/ok" {
		t.Error("Wrong tags collection")
	}
	// Wrong ii/ok
	tags.II = ""
	_, err = tags.CollectTags()
	if err == nil {
		t.Error("Wrong tags collection")
	}
}

func TestPrepareMessageForSend(t *testing.T) {
	p := &PointMessage{
		Echo: "ii.test.14",
		To:   "All",
		Subg: "Test message",
		Body: "This is a message body.",
	}

	result := p.PrepareMessageForSend()
	if result == "" {
		t.Error("Prepared message is empty")
	}
	// With repto
	p.Repto = "hXzRNEzmMuzKkT1HCxUb"
	result2 := p.PrepareMessageForSend()
	if result2 == "" {
		t.Error("Prepared message is empty")
	}

	if result == result2 {
		t.Error("Messages with and without repto is equal!")
	}
}
