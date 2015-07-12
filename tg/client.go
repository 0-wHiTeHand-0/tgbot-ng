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

func (c *Client) GetUpdates() (UpdateResponse, error) {
	resp, err := http.PostForm(baseURL+c.Token+"/getUpdates",
		url.Values{"offset": {strconv.Itoa(c.lastUpdateID + 1)}})
	if err != nil {
		return UpdateResponse{}, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return UpdateResponse{}, err
	}
	var ur UpdateResponse
	if err := json.Unmarshal(b, &ur); err != nil {
		return UpdateResponse{}, err
	}
	if !ur.Ok {
		return UpdateResponse{}, errors.New("update is not OK")
	}
	for _, u := range ur.Updates {
		if u.UpdateID > c.lastUpdateID {
			c.lastUpdateID = u.UpdateID
		}
	}
	return ur, nil
}

func (c *Client) SendMessage(chatID int, text string) (Response, error) {
	resp, err := http.PostForm(baseURL+c.Token+"/sendMessage",
		url.Values{"chat_id": {strconv.Itoa(chatID)}, "text": {text}})
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	var r Response
	if err := json.Unmarshal(body, &r); err != nil {
		return Response{}, err
	}
	if !r.Ok {
		return Response{}, errors.New("response is not OK")
	}

	return r, nil
}

func (c *Client) SendPhoto(chatID int, filename string, data []byte) (Response, error) {
	params := map[string]string{"chat_id": strconv.Itoa(chatID)}
	return c.uploadFile("sendPhoto", "photo", filename, data, params)
}

func (c *Client) SendDocument(chatID int, filename string, data []byte) (Response, error) {
	params := map[string]string{"chat_id": strconv.Itoa(chatID)}
	return c.uploadFile("sendDocument", "document", filename, data, params)
}

func (c *Client) uploadFile(endpoint, fieldname, filename string, data []byte, params map[string]string) (Response, error) {
	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)

	fw, err := w.CreateFormFile(fieldname, filename)
	if err != nil {
		return Response{}, err
	}
	if _, err := fw.Write(data); err != nil {
		return Response{}, err
	}

	for param, value := range params {
		fw, err = w.CreateFormField(param)
		if err != nil {
			return Response{}, err
		}
		if _, err := fw.Write([]byte(value)); err != nil {
			return Response{}, err
		}
	}

	w.Close()

	resp, err := http.Post(baseURL+c.Token+"/"+endpoint, w.FormDataContentType(), b)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var r Response
	if err := json.Unmarshal(body, &r); err != nil {
		return Response{}, err
	}

	if !r.Ok {
		return Response{}, errors.New("response is not OK")
	}

	return r, nil
}
