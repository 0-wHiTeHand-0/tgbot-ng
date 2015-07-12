// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"time"

	"github.com/jroimartin/tgbot-ng/tg"
)

type bot struct {
	cli            *tg.Client
	allowedIDs     []int
	updateInterval time.Duration
	commands       []Command
}

func newBot(name, token string) *bot {
	return &bot{
		cli:            tg.NewClient(name, token),
		updateInterval: 2 * time.Second,
	}
}

func (b *bot) setAllowedIDs(ids []int) {
	if len(ids) == 0 {
		return
	}
	b.allowedIDs = make([]int, len(ids))
	copy(b.allowedIDs, ids)
}

func (b *bot) setUpdateInterval(nseg int) {
	b.updateInterval = time.Duration(nseg) * time.Second
}

func (b *bot) addCommand(cmd Command) {
	b.commands = append(b.commands, cmd)
}

func (b *bot) loop() {
	for {
		ur, err := b.cli.GetUpdates()
		if err != nil {
			log.Println("error:", err)
		}
		for _, u := range ur.Updates {
			go b.handleUpdate(u)
		}
		time.Sleep(b.updateInterval)
	}
}

func (b *bot) handleUpdate(u tg.Update) {
	log.Printf("update: %+v\n", u)
	if !b.isAllowed(u) {
		log.Println("error: not allowed")
	}
	for _, cmd := range b.commands {
		if cmd.Match(u.Message.Text) {
			if err := cmd.Run(u.Message.Chat.ID, u.Message.Text); err != nil {
				log.Printf("error: %v\n", err)
			}
			break
		}
	}
}

func (b *bot) isAllowed(u tg.Update) bool {
	for _, aid := range b.allowedIDs {
		if u.Message.Chat.ID == aid {
			return true
		}
	}
	return false
}
