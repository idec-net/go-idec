package idec

import (
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestNewExtensions(t *testing.T) {
	e := NewExtensions()
	if e.ListTXT != "list.txt" {
		t.Error("list.txt not found")
	}
}

func TestGetMessagesIDS(t *testing.T) {
	httpmock.Activate()
	fc := FetchConfig{
		Node:   "http://localhost/idec/",
		Echoes: []string{"ii.test.14"},
		Num:    5,
		Offset: -5,
		Limit:  5,
	}
	httpmock.RegisterResponder("GET", "http://localhost/idec/u/e/ii.test.14/-5:5", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, `ii.test.14
hXzRNEzmMuzKkT1HCxUb
JN3ylpxjaNofxgPy6NhL
xF3kkmrZYld330BO7qaA
3uS3uij0Y4AUnSxhf4WB
zi9YpQGddLW5WQKi9GMf`)
		return resp, nil
	})
	id, err := fc.GetMessagesIDS()
	if err != nil {
		t.Error(err)
	}
	if len(id) == 0 {
		t.Error("Message not fetched")
	}
}

func TestGetAllMessagesIDS(t *testing.T) {
	httpmock.Activate()
	fc := FetchConfig{
		Node:   "http://localhost/idec/",
		Echoes: []string{"ii.test.14"},
	}
	httpmock.RegisterResponder("GET", "http://localhost/idec/u/e/ii.test.14", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, `ii.test.14
hXzRNEzmMuzKkT1HCxUb
JN3ylpxjaNofxgPy6NhL
xF3kkmrZYld330BO7qaA
3uS3uij0Y4AUnSxhf4WB
zi9YpQGddLW5WQKi9GMf`)
		return resp, nil
	})
	id, err := fc.GetAllMessagesIDS()
	if err != nil {
		t.Error(err)
	}
	if len(id) == 0 {
		t.Error("Message not fetched")
	}
}

func TestGetRawMessages(t *testing.T) {
	httpmock.Activate()
	fc := FetchConfig{
		Node:   "http://localhost/idec/",
		Echoes: []string{"ii.test.14"},
	}
	httpmock.RegisterResponder("GET", "http://localhost/idec/u/m/hXzRNEzmMuzKkT1HCxUb/JN3ylpxjaNofxgPy6NhL", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, `hXzRNEzmMuzKkT1HCxUb:aWkvb2svcmVwdG8vSk4zeWxweGphTm9meGdQeTZOaEwKaWkudGVzdC4xNAoxNTUxNjk5NjE0CkRpZnJleApkeW5hbWljLDEKRGlmcmV4ClJlOiBpZGVjCgpzZGZnc2ZkZyBzZGdmCgpmZGcgc2dmIHMKCmZkZyBzZGZnIHMKc2RmZyBzZGZnIHMKCgoKc2ZkZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dn
JN3ylpxjaNofxgPy6NhL:aWkvb2svcmVwdG8veEYza2ttclpZbGQzMzBCTzdxYUEKaWkudGVzdC4xNAoxNTUxNjk5NTk1CkRpZnJleApkeW5hbWljLDEKRGlmcmV4ClJlOiBpZGVjCgpKS0hKS0hLSlNGSCBscwo9PT0gZHMKc2RmZyBzZGZnCmdmc2QgZGZnc2dmIGYKPT09PQpzaCAtYyAnZWNobyBPSycKPT09PQ==`)
		return resp, nil
	})

	ids := []ID{ID{"ii.test.14", "hXzRNEzmMuzKkT1HCxUb"}, ID{"ii.test.14", "JN3ylpxjaNofxgPy6NhL"}}

	msg, err := fc.GetRawMessages(ids)
	if err != nil {
		t.Error(err)
	}
	if len(msg) == 0 {
		t.Error("Messages not fetched")
	}
}

func TestGetEchoList(t *testing.T) {
	httpmock.Activate()
	fc := FetchConfig{
		Node:   "http://localhost/idec/",
		Echoes: []string{"ii.test.14"},
	}

	httpmock.RegisterResponder("GET", "http://localhost/idec/x/features", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, `list.txt
u/e
u/m
x/c`)
		return resp, nil
	})

	httpmock.RegisterResponder("GET", "http://localhost/idec/list.txt", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, `bash.rss:14573:RSS с сайта bash.im
creepy.14:334:Страшные истории
develop.16:402:Обсуждение вопросов программирования
file.wishes:10:Поиск файлов`)
		return resp, nil
	})

	echoes, err := fc.GetEchoList()
	if err != nil {
		t.Error(err)
	}
	if len(echoes) == 0 {
		t.Error("Wrong echoes list")
	}

	httpmock.RegisterResponder("GET", "http://localhost/idec/x/features", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, `u/e
u/m
x/c`)
		return resp, nil
	})
	_, err = fc.GetEchoList()
	if err == nil {
		t.Error(err)
	}
}

func TestPostMessage(t *testing.T) {
	httpmock.Activate()
	fc := FetchConfig{
		Node:   "http://localhost/idec/",
		Echoes: []string{"ii.test.14"},
	}

	message := "aWkvb2svcmVwdG8veEYza2ttclpZbGQzMzBCTzdxYUEKaWkudGVzdC4xNAoxNTUxNjk5NTk1CkRpZnJleApkeW5hbWljLDEKRGlmcmV4ClJlOiBpZGVjCgpKS0hKS0hLSlNGSCBscwo9PT0gZHMKc2RmZyBzZGZnCmdmc2QgZGZnc2dmIGYKPT09PQpzaCAtYyAnZWNobyBPSycKPT09PQ=="

	httpmock.RegisterResponder("POST", "http://localhost/idec/u/point", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, `msg ok`)
		return resp, nil
	})

	err := fc.PostMessage("auth", message)
	if err != nil {
		t.Error(err)
	}

	httpmock.RegisterResponder("POST", "http://localhost/idec/u/point", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(403, `error: wrong authstring`)
		return resp, nil
	})

	err = fc.PostMessage("auth", message)
	if err == nil {
		t.Error("Errors not precessed")
	}
}
