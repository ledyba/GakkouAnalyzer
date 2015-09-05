package client

import "testing"

func TestResInfo(t *testing.T) {

	cl, err := LoginWithPassword(TestUser, TestPass)
	if err != nil {
		t.Fatal(err)
	}
	res, err := cl.GetResInfo("sm60")
	if err != nil {
		t.Fatal(err)
	}
	if res.ThreadID == "" {
		t.Fatal("Empty Thread ID")
	}
}
