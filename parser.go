package idec

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseMessage ...
func ParseMessage(message string) (Message, error) {
	var m Message
	plainMessage, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return m, err
	}

	txtMessage := strings.Split(string(plainMessage), "\n")

	var body string
	for i := 8; i < len(txtMessage); i++ {
		body = strings.Join([]string{body, txtMessage[i]}, "\n")
	}

	ts, err := strconv.Atoi(txtMessage[2])
	if err != nil {
		return m, err
	}

	tags, err := ParseTags(txtMessage[0])

	m.Tags = tags
	m.Echo = txtMessage[1]
	m.Timestamp = ts
	m.From = txtMessage[3]
	m.Address = txtMessage[4]
	m.To = txtMessage[5]
	m.Subg = txtMessage[6]
	m.Body = body

	return m, err
}

// ParsePointMessage ...
func ParsePointMessage(message string) (*PointMessage, error) {
	var pointMessage *PointMessage
	plainMessage, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return pointMessage, err
	}

	txtMessage := strings.Split(string(plainMessage), "\n")
	if len(txtMessage) < 6 {
		e := errors.New("Bad message")
		return pointMessage, e
	}

	var body string
	for i := 5; i < len(txtMessage); i++ {
		body = strings.Join([]string{body, txtMessage[i]}, "\n")
	}

	pointMessage = &PointMessage{
		Echo:      txtMessage[0],
		To:        txtMessage[1],
		Subg:      txtMessage[2],
		EmptyLine: txtMessage[3],
		Repto:     txtMessage[4],
		Body:      body,
	}

	return pointMessage, nil
}

// MakeBundledMessage from point message.
// Returns Message with empty From and Address fields
// you must set this somewhere outside
func MakeBundledMessage(pointMessage *PointMessage) (Message, error) {
	var msg Message
	t := "ii/ok"
	if pointMessage.Repto != "" {
		t = fmt.Sprintf("%s/repto/%s", t, pointMessage.Repto)
	}
	tags, err := ParseTags(t)
	if err != nil {
		return msg, err
	}
	msg = Message{
		Tags:      tags,
		Echo:      pointMessage.Echo,
		Timestamp: int(time.Now().Unix()),
		To:        pointMessage.To,
		Subg:      pointMessage.Subg,
		Repto:     pointMessage.Repto,
		Body:      pointMessage.Body,
	}

	return msg, nil
}

// parseTags parse message tags and return Tags struct
func ParseTags(tags string) (Tags, error) {
	var t Tags

	if !strings.Contains(tags, "ii/") {
		e := errors.New("Bad tagstring")
		return t, e
	}

	tagsSlice := strings.Split(tags, "/")
	if len(tagsSlice) < 4 {
		t.II = tagsSlice[1]
		return t, nil
	}
	t.II = tagsSlice[1]
	t.Repto = tagsSlice[3]
	return t, nil
}

// ParseEchoList parse /list.txt
func ParseEchoList(list string) ([]Echo, error) {
	var echoes []Echo
	for _, e := range strings.Split(list, "\n") {
		desc := strings.Split(e, ":")
		if len(desc) <= 1 {
			break
		}
		count, err := strconv.Atoi(desc[1])
		if err != nil {
			return echoes, err
		}
		echoes = append(echoes, Echo{desc[0], count, desc[2]})
	}

	return echoes, nil
}

// MakeMsgID from provided plain bundled message
func MakeMsgID(msg string) string {
	id := string(
		sha256.New().Sum(
			[]byte(base64.StdEncoding.EncodeToString(
				[]byte(msg)))))[:20] // LISP style, LOL
	id = strings.Replace(id, "+", "A", -1)
	id = strings.Replace(id, "/", "Z", -1)
	return id
}

// String from PointMessage
func (p *PointMessage) String() string {
	return strings.Join([]string{
		p.Echo,
		p.To,
		p.Subg,
		"",
		p.Repto,
		p.Body,
	}, "\n")
}
