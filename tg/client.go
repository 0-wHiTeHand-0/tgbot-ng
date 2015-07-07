// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const baseURL = "https://api.telegram.org/bot"

type Client struct {
	Name string
	Token string
	
	lastUpdateID int
}

func NewClient(name, token string) *Client {
	return &Client{Name:name, Token:token}
}

func (c *Client) GetUpdates() (Update, error) {
	b, err := c.doRequest("getUpdates",
		url.Values{"offset": {strconv.Itoa(c.lastUpdateID +1)}})
	if err != nil {
		return Update{}, err
	}	
	var update Update
	if err := json.Unmarshal(b, &update); err != nil {
		return Update{}, err
	}
	for _, r := range(update.Results) {
		if r.UpdateID > c.lastUpdateID {
			c.lastUpdateID = r.UpdateID
		}
	}
	return update, nil
}

func (c *Client) SendMessage(chatID int, t string) error {
	_, err := c.doRequest("sendMessage",
		url.Values{"chat_id": {strconv.Itoa(chatID)}, "text": {t}})
	return err
}

func (c *Client) doRequest(command string, data url.Values) (body []byte, err error) {
	resp, err := http.PostForm(baseURL + c.Token + "/" + command, data)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}