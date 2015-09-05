package client

import "testing"

func TestComment(t *testing.T) {

	cl, err := LoginWithPassword(TestUser, TestPass)
	if err != nil {
		t.Fatal(err)
	}
	res, err := cl.GetResInfo("1436342441")
	if err != nil {
		t.Fatal(err)
	}
	_, err = cl.GetComment(res, res.OptionalThreadID, 1441441656)
	if err != nil {
		t.Fatal(err)
	}
	chats, err := cl.GetComment(res, res.ThreadID, 1441441656)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(len(chats))
	last := &chats[len(chats)-1]
	chats, err = cl.GetComment(res, res.ThreadID, last.Date)
	t.Fatal(chats)

}
