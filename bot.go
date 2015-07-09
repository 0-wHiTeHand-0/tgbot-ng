// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"log"
	"time"

	"github.com/jroimartin/tgbot-ng/tg"
)

type bot struct {
	cli            *tg.Client
	allowedIDs     []int
	updateInterval time.Duration
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

func (b *bot) loop() {
	for {
		update, err := b.cli.GetUpdates()
		if err != nil {
			log.Println("error:", err)
		}
		for _, r := range update.Results {
			if err := b.handleResult(r); err != nil {
				log.Println("error:", err)
			}
		}
		time.Sleep(b.updateInterval)
	}
}

func (b *bot) handleResult(r tg.Result) error {
	log.Printf("result: %+v\n", r)
	if !b.isAllowed(r) {
		return errors.New("not allowed")
	}
	// TODO
	return nil
}

func (b *bot) isAllowed(r tg.Result) bool {
	for _, aid := range b.allowedIDs {
		if r.Message.Chat.ID == aid {
			return true
		}
	}
	return false
}
