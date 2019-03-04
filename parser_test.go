package idec

import (
	"encoding/base64"
	"testing"
)

func TestParseMessage(t *testing.T) {
	m := `aWkvb2svcmVwdG8vdFU3SjBienVMMnI0RG9zRGtOUE8KcGlwZS4yMDMyCjE1NTE2ODk3NjYKRGlmcmV4CmR5bmFtaWMsMQpEaWZyZXgKUmU6IGlkZWMKCtCY0LvQuCDQtNCw0LbQtSDRgtCw0Lo6Cj09PT0KY3VybCAtWFBPU1QgLUggIlgtSWRlYy1QYXV0aDogc2Rsa2RzZmprbHNkZiIgLVQgL2V0Yy9wYXNzd2QgaWRlYy5ub2RlL3gvZC9tc2dpZAo9PT09`
	msg, err := ParseMessage(m)
	if err != nil {
		t.Error(err)
	}
	if msg.Tags.II != "ok" {
		t.Error("Bad message tags")
	}

	// Wrong unixtime
	m = `aWkvb2svcmVwdG8vdFU3SjBienVMMnI0RG9zRGtOUE8KcGlwZS4yMDMyCgpEaWZyZXgKZHluYW1p
YywxCkRpZnJleApSZTogaWRlYwoK0JjQu9C4INC00LDQttC1INGC0LDQujoKPT09PQpjdXJsIC1Y
UE9TVCAtSCAiWC1JZGVjLVBhdXRoOiBzZGxrZHNmamtsc2RmIiAtVCAvZXRjL3Bhc3N3ZCBpZGVj
Lm5vZGUveC9kL21zZ2lkCj09PT0K`
	_, err = ParseMessage(m)
	if err == nil {
		t.Error("Wrong time parsing")
	}

	// Test bad message
	m = `aWkvb2svcmVwdG8vdFU3SjBienVMMnI0RG9zRGtOUE8KcGlwZS4yMDMyCjE1NTE2ODk3NjYKRGlmcmV4CmR5bmFtaWMsMQpEaWZgKUmU6IGlkZWMKCtCY0LvQuCDQtNCw0LbQtSDRgtCw0Lo6Cj09PT0KY3VybCAtWFBPU1QgLUggIlgtSWRlYy1QYXV0aDogc2Rsa2RzZmprbHNkZiIgLVQgL2V0Yy9wYXNzd2QgaWRlYy5ub2RlL3gvZC9tc2dpZAo9P`
	_, err = ParseMessage(m)
	if err == nil {
		t.Error(err)
	}
}

func TestParsePointMessage(t *testing.T) {
	m := `ii.test.14
Difrex
Test message

@repto:EviyYJSFrnubg0DvckW9
This is a message body string.`
	mb64 := base64.StdEncoding.EncodeToString([]byte(m))
	pmsg, err := ParsePointMessage(mb64)
	if err != nil {
		t.Error(err)
	}
	if pmsg.Echo != "ii.test.14" {
		t.Error("Wrong echo parsing")
	}
	if pmsg.Body != "\nThis is a message body string." {
		t.Errorf("Wrong body parsing, b: %s", pmsg.Body)
	}
	if pmsg.Repto != "EviyYJSFrnubg0DvckW9" {
		t.Error("Wrong repto parsing")
	}

	// Test point message string
	s := pmsg.String()
	if s == "" {
		t.Errorf("Wrong message string")
	}

	// Without repto
	m = `ii.test.14
All
Test message


This is a message body string.`
	mb64 = base64.StdEncoding.EncodeToString([]byte(m))
	pmsg, err = ParsePointMessage(mb64)
	if err != nil {
		t.Error(err)
	}

	// Bad message
	m = `ii.test.14
All
Test message
This is a message body string.`
	mb64 = base64.StdEncoding.EncodeToString([]byte(m))
	pmsg, err = ParsePointMessage(mb64)
	if err == nil {
		t.Error("Wrong bad message detection")
	}

	pmsg, err = ParsePointMessage("BlaBlaBla")
	if err == nil {
		t.Error("Wrong base64 decryption")
	}

	// Wrong unsafe
	pmsg, err = ParsePointMessage("aWl0ZXN0CkRpZnJleApSZTogaWRlYwoKQHJlcHRvOmhYelJORXptTXV6S2tUMUhDeFViCnNhZGZhZiBhcwpmZCBhc2ZkCmRmIGFzZiBhCgpmYWQgc2Y9PT09PT09PT09PT09PT09PT09PT09PT09Cj09PT09PT09PT09PT09PT0KCj09PT0KCXBvaW50TWVzc2FnZSA9ICZQb2ludE1lc3NhZ2V7CgkJRWNobzogICAgICBzdHJpbmdzLlRyaW0odHh0TWVzc2FnZVswXSwgIiAiKSwKCQlUbzogICAgICAgIHR4dE1lc3NhZ2VbMV0sCgkJU3ViZzogICAgICB0eHRNZXNzYWdlWzJdLAoJCUVtcHR5TGluZTogdHh0TWVzc2FnZVszXSwKCQlCb2R5OiAgICAgIGJvZHksCgl9Cj09PT0%rt0")
	if err == nil {
		t.Error("Wrong unsafing")
	}
}

func TestValidate(t *testing.T) {
	m := `ii.test.14
Difrex
Test message

@repto:EviyYJSFrnubg0DvckW9
This is a message body string.`
	mb64 := base64.StdEncoding.EncodeToString([]byte(m))
	pmsg, err := ParsePointMessage(mb64)
	if err != nil {
		t.Error(err)
	}
	if err := pmsg.Validate(); err != nil {
		t.Error(err)
	}
	// Test wrong echo
	pmsg.Echo = "invalid"
	if err := pmsg.Validate(); err != nil && err.Error() != "Wrong Echo name" {
		t.Error("Validating echo(wrong) field is broken")
	} else {
		pmsg.Echo = "ii.test.14"
	}
	pmsg.Echo = ""
	if err := pmsg.Validate(); err != nil && err.Error() != "Wrong Echo name" {
		t.Error("Validating echo(empty) field is broken: ", err)
	} else {
		pmsg.Echo = "ii.test.14"
	}
	pmsg.Echo = "ii.test.14"
	// Test empty to
	pmsg.To = ""
	if err := pmsg.Validate(); err != nil && err.Error() != "`To' field is empty" {
		t.Error("Validating to field is broken")
	} else {
		pmsg.To = "Difrex"
	}
	// Test subg
	pmsg.Subg = ""
	if err := pmsg.Validate(); err != nil && err.Error() != "`Subg' field is empty" {
		t.Error("Validating subg field is broken")
	} else {
		pmsg.Subg = "Test message"
	}
	// Test empty line
	pmsg.EmptyLine = "not empty"
	if err := pmsg.Validate(); err != nil && err.Error() != "EmptyLine is not empty" {
		t.Error("Validating empty line field is broken")
	} else {
		pmsg.EmptyLine = ""
	}
	// Test body
	pmsg.Body = ""
	if err := pmsg.Validate(); err != nil && err.Error() != "`Body' field is empty" {
		t.Error("Validating body field is broken")
	} else {
		pmsg.Body = "\nThis is a message body string."
	}
	// Test repto
	pmsg.Repto = "EviyYJSFrnubg0DvckW"
	if err := pmsg.Validate(); err != nil && err.Error() != "Wrong @repto field length" {
		t.Error("Validating repto field is broken")
	}
}

func TestParseReptoField(t *testing.T) {
	repto := "@repto:EviyYJSFrnubg0DvckWa"
	if ParseReptoField(repto) != "EviyYJSFrnubg0DvckWa" {
		t.Error("Can't parse repto field")
	}
}

func TestMakeBundledMessage(t *testing.T) {
	m := `ii.test.14
Difrex
Test message

@repto:EviyYJSFrnubg0DvckW9
This is a message body string.`
	mb64 := base64.StdEncoding.EncodeToString([]byte(m))
	pmsg, err := ParsePointMessage(mb64)
	if err != nil {
		t.Error(err)
	}
	msg, err := MakeBundledMessage(pmsg)
	if err != nil {
		t.Error(err)
	}
	if msg.Echo != "ii.test.14" {
		t.Error("Wrong bundle creation")
	}
}

func TestParseTags(t *testing.T) {
	tg := "ii/ok"
	tags, err := ParseTags(tg)
	if err != nil {
		t.Error(err)
	}
	if tags.II != "ok" {
		t.Error("Wrong ii tag")
	}
	tg = "wrong/tags"
	tags, err = ParseTags(tg)
	if err != nil && err.Error() != "Bad tagstring" {
		t.Error("Wrong tags parsing")
	}
}

func TestParseEchoList(t *testing.T) {
	list := `bash.rss:14573:RSS с сайта bash.im
creepy.14:334:Страшные истории
develop.16:402:Обсуждение вопросов программирования
file.wishes:10:Поиск файлов
`
	echoes, err := ParseEchoList(list)
	if err != nil {
		t.Error(err)
	}
	if len(echoes) == 0 {
		t.Error("Wrong echoes list")
	}
	if echoes[0].Description != "RSS с сайта bash.im" {
		t.Error("Wrong description parsing")
	}

	// Wrong echo list
	list = `bash.rss:broken:RSS с сайта bash.im
creepy.14:sasdads:Страшные истории
develop.16:402:Обсуждение вопросов программирования
file.wishes:10:Поиск файлов
`
	echoes, err = ParseEchoList(list)
	if err == nil {
		t.Error("Wrong echoes list parsing")
	}
}

func TestMakeMsgID(t *testing.T) {
	m := `ii/ok/repto/tU7J0bzuL2r4DosDkNPO
pipe.2032
1551689766
Difrex
dynamic,1
Difrex
Re: idec

Или даже так:
====
curl -XPOST -H "X-Idec-Pauth: sdlkdsfjklsdf" -T /etc/passwd idec.node/x/d/msgid
====
`
	id := MakeMsgID(m)
	if id != "Jc0StQZltt2EoHV9fLee" {
		t.Errorf("id %s not equal %s", id, "Jc0StQZltt2EoHV9fLee")
	}
}
