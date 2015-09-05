package client

import "testing"

const TestUser = "saccubus@gmail.com"
const TestPass = "test1234"

func TestLogin(t *testing.T) {

	_, err := LoginWithPassword(TestUser, TestPass)
	if err != nil {
		t.Error(err)
	}

}
