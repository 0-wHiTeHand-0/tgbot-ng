// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/jroimartin/tgbot-ng/tg"
)

type Command interface {
	Match(text string) bool
	Run(chatID, replyID int, text string, from string, reply_id *tg.Message) error
}
