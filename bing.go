// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"regexp"

	"github.com/jroimartin/tgbot-ng/bing"
	"github.com/jroimartin/tgbot-ng/tg"
)

type cmdBing struct {
	re     *regexp.Regexp
	config CmdConfigBing
	cli    *tg.Client
}

type CmdConfigBing struct {
	Enabled     bool   `json:"enabled"`
	ApiKey      string `json:"api_key"`
	SearchLimit int    `json:"search_limit"`
}

func NewCmdBing(config CmdConfigBing, cli *tg.Client) Command {
	if config.SearchLimit < 1 {
		config.SearchLimit = 10
	}
	return &cmdBing{
		re:     regexp.MustCompile(`^/bing(?:@[^ ]+?)?(?:$| +(.+)$)`),
		config: config,
		cli:    cli,
	}
}

func (cmd *cmdBing) Match(text string) bool {
	return cmd.re.MatchString(text)
}

func (cmd *cmdBing) Run(chatID, replyID int, text string, from string) error {
	var query string

	m := cmd.re.FindStringSubmatch(text)
	if len(m) == 2 {
		query = m[1]
	}

	if query == "" {
		cmd.cli.SendText(chatID, "Nope. Try again.")
	} else {
		if _, err := cmd.cli.SendText(chatID, "Don't byte off more than you can view."); err != nil {
			return err
		}
		img, err := cmd.search(query)
		if err != nil {
			return err
		}
		if path.Ext(img.Name) == ".gif" {
			_, err = cmd.cli.SendDocument(chatID, img)
		} else {
			_, err = cmd.cli.SendPhoto(chatID, img)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *cmdBing) search(query string) (img tg.File, err error) {
	c := bing.NewClient(cmd.config.ApiKey)
	if cmd.config.SearchLimit > 0 {
		c.Limit = cmd.config.SearchLimit
	}

	results, err := c.Query(bing.Image, query)
	if err != nil {
		return tg.File{}, err
	}
	if len(results) == 0 {
		return tg.File{}, errors.New("no pics")
	}
	rndInt := rand.Intn(len(results))

	resp, err := http.Get(results[rndInt].MediaUrl)
	if err != nil {
		return tg.File{}, err
	}
	defer resp.Body.Close()
	imgData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tg.File{}, err
	}

	return tg.File{Name: results[rndInt].MediaUrl, Data: imgData}, nil
}
