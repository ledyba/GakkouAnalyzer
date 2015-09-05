package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func (cl *Client) GetCommentKey(info *ResInfo, tid string) (url.Values, error) {
	resp, err := cl.cl.Get(fmt.Sprintf("http://flapi.nicovideo.jp/api/getthreadkey?thread=%s", tid))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (cl *Client) GetWaybackKey(info *ResInfo, tid string) (string, error) {
	resp, err := cl.cl.Get(fmt.Sprintf("http://flapi.nicovideo.jp/api/getwaybackkey?thread=%s", tid))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", err
	}
	return values.Get("waybackkey"), nil
}

type Chat struct {
	Thread    int64  `json:"thread"`
	No        int64  `json:"no"`
	Vpos      int64  `json:"vpos"`
	Date      int64  `json:"date"`
	Mail      string `json:"mail"`
	UserID    string `json:"user_id"`
	Anonymity int64  `json:"anonymity"`
	Leaf      int64  `json:"leaf"`
	Content   string `json:"content"`
}

type ChatList []Chat

func (p ChatList) Len() int {
	return len(p)
}

func (p ChatList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ChatList) Less(i, j int) bool {
	return p[i].Date < p[j].Date
}

func (cl *Client) GetComment(info *ResInfo, tid string, when int64) ([]Chat, error) {
	values := &url.Values{}
	values.Add("version", "20090904")
	values.Add("thread", info.ThreadID)
	values.Add("res_from", "-1000")
	values.Add("user_id", info.UserID)
	if info.NeedsKey {
		key, err := cl.GetCommentKey(info, tid)
		if err != nil {
			return nil, err
		}
		for k := range key {
			values.Add(k, key.Get(k))
		}
	}
	if when >= 0 {
		key, err := cl.GetWaybackKey(info, tid)
		if err != nil {
			return nil, err
		}
		values.Add("when", strconv.FormatInt(when, 10))
		values.Add("waybackkey", key)
	}
	u := strings.Replace(info.MessageURL, "/api/", "/api.json/", -1)
	resp, err := cl.cl.Get(fmt.Sprintf("%sthread?%s", u, values.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var els []map[string]*Chat
	decoder.Decode(&els)

	var chats []Chat
	for _, v := range els {
		chat, has := v["chat"]
		if has {
			chats = append(chats, *chat)
		}
	}
	sort.Sort(ChatList(chats))

	return chats, nil
}
