package client

import (
	"errors"
	"log"
	"net/url"
)

const LoginURL = "https://secure.nicovideo.jp/secure/login?show_button_twitter=1&site=niconico&show_button_facebook=1&next_url="
const CookieName = "user_session"

var LoginFail = errors.New("Login Failed")

func LoginWithPassword(user, pass string) (*Client, error) {
	cl := NewClient()
	data := url.Values{}
	data.Add("mail_tel", user)
	data.Add("password", pass)
	resp, err := cl.cl.PostForm(LoginURL, data)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	u, err := url.Parse("http://www.nicovideo.jp/")
	if err != nil {
		log.Fatal(err)
	}
	cookies := cl.cl.Jar.Cookies(u)
	if cookies == nil {
		log.Fatal("!!BUG!!")
	}

	for _, cookie := range cl.cl.Jar.Cookies(u) {
		name := cookie.Name

		if name == CookieName {
			return cl, nil
		}
	}
	return nil, LoginFail
}
