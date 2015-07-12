// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

const baseURL = "https://api.telegram.org/bot"

type Client struct {
	Name  string
	Token string

	lastUpdateID int
}

func NewClient(name, token string) *Client {
	return &Client{Name: name, Token: token}
}

func (c *Client) GetUpdates() (Update, error) {
	resp, err := http.PostForm(baseURL+c.Token+"/getUpdates",
		url.Values{"offset": {strconv.Itoa(c.lastUpdateID + 1)}})
	if err != nil {
		return Update{}, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Update{}, err
	}
	var update Update
	if err := json.Unmarshal(b, &update); err != nil {
		return Update{}, err
	}
	if !update.Ok {
		return Update{}, errors.New("update is not OK")
	}
	for _, r := range update.Results {
		if r.UpdateID > c.lastUpdateID {
			c.lastUpdateID = r.UpdateID
		}
	}
	return update, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	resp, err := http.PostForm(baseURL+c.Token+"/sendMessage",
		url.Values{"chat_id": {strconv.Itoa(chatID)}, "text": {text}})
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) SendPhoto(chatID int, filename string, data []byte) error {
	params := map[string]string{"chat_id": strconv.Itoa(chatID)}
	return c.uploadFile("sendPhoto", "photo", filename, data, params)
}

func (c *Client) SendDocument(chatID int, filename string, data []byte) error {
	params := map[string]string{"chat_id": strconv.Itoa(chatID)}
	return c.uploadFile("sendDocument", "document", filename, data, params)
}

func (c *Client) uploadFile(endpoint, fieldname, filename string, data []byte, params map[string]string) error {
	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)

	fw, err := w.CreateFormFile(fieldname, filename)
	if err != nil {
		return err
	}
	if _, err := fw.Write(data); err != nil {
		return err
	}

	for param, value := range params {
		fw, err = w.CreateFormField(param)
		if err != nil {
			return err
		}
		if _, err := fw.Write([]byte(value)); err != nil {
			return err
		}
	}

	w.Close()

	resp, err := http.Post(baseURL+c.Token+"/"+endpoint, w.FormDataContentType(), b)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
