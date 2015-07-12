// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tg

type UpdateResponse struct {
	Ok      bool     `json:"ok"`
	Updates []Update `json:"result"`
}

type Response struct {
	Ok      bool    `json:"ok"`
	Message Message `json:"result"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	ID       int         `json:"message_id"`
	From     User        `json:"from"`
	Date     int         `json:"date"`
	Chat     Chat        `json:"chat"`
	Text     string      `json:"text"`
	Photo    []PhotoSize `json:"photo"`
	Document Document    `json:"document"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID int
}

type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

type Document struct {
	FileID   string    `json:"file_id"`
	Thumb    PhotoSize `json:"thumb"`
	FileName string    `json:"file_name"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
}

type ReplyKeyboardMarkup struct {
	Keyboard  [][]string `json:"keyboard"`
	Resize    bool       `json:"resize_keyboard"`
	OneTime   bool       `json:"one_time_keyboard"`
	Selective bool       `json:"selective"`
}
