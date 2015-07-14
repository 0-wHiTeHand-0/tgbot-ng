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

func (c *Client) SendText(chatID int, text string) (Response, error) {
	params := SendMessageParams{
		ChatID: chatID,
		Text: text,
	}
	return c.SendMessage(params)
}

func (c *Client) SendTextReply(chatID, replyID int, text string) (Response, error) {
	params := SendMessageParams{
		ChatID: chatID,
		ReplyID: replyID,
		Text: text,
	}
	return c.SendMessage(params)
}

func (c *Client) SendKbd(chatID, replyID int, text string, kbd ReplyKeyboardMarkup) (Response, error) {
	params := SendMessageParams{
		ChatID: chatID,
		ReplyID: replyID,
		Text: text,
		ReplyMarkup: kbd,
	}
	return c.SendMessage(params)
}

func (c *Client) SendMessage(params SendMessageParams) (Response, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(params.ChatID))
	v.Add("text", params.Text)
	if params.ReplyID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(params.ReplyID))
	}
	if params.ReplyMarkup != nil {
		b, err := json.Marshal(params.ReplyMarkup)
		if err != nil {
			return Response{}, err
		}
		v.Add("reply_markup", string(b))
	}
	
	resp, err := http.PostForm(baseURL+c.Token+"/sendMessage", v)
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

func (c *Client) SendPhoto(chatID int, photo File) (Response, error) {
	params := map[string]string{"chat_id": strconv.Itoa(chatID)}
	return c.uploadFile("sendPhoto", "photo", photo, params)
}

func (c *Client) SendDocument(chatID int, doc File) (Response, error) {
	params := map[string]string{"chat_id": strconv.Itoa(chatID)}
	return c.uploadFile("sendDocument", "document", doc, params)
}

func (c *Client) uploadFile(endpoint, fieldname string, file File, params map[string]string) (Response, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile(fieldname, file.Name)
	if err != nil {
		return Response{}, err
	}
	if _, err := fw.Write(file.Data); err != nil {
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

	resp, err := http.Post(baseURL+c.Token+"/"+endpoint, w.FormDataContentType(), &b)
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
