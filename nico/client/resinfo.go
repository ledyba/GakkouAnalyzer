package client

import (
	"fmt"
	"io/ioutil"
	"net/url"
)

const resInfoUrl = "http://flapi.nicovideo.jp/api/getflv/%s?watch_harmful=1&as3=1"

type ResInfo struct {
	UserID           string
	MessageURL       string
	ThreadID         string
	OptionalThreadID string
	NeedsKey         bool
}

func (cl *Client) GetResInfo(vid string) (*ResInfo, error) {
	resp, err := cl.cl.Get(fmt.Sprintf(resInfoUrl, vid))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := string(bbody)
	values, err := url.ParseQuery(body)
	if err != nil {
		return nil, err
	}
	res := &ResInfo{}
	res.UserID = values.Get("user_id")
	res.MessageURL = values.Get("ms")
	res.ThreadID = values.Get("thread_id")
	res.OptionalThreadID = values.Get("optional_thread_id")
	res.NeedsKey = values.Get("needs_key") != ""
	return res, err
}
