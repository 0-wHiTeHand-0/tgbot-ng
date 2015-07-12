// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"math/rand"
	"path"
	"regexp"
	"net/http"
	"io/ioutil"

	"github.com/jroimartin/tgbot-ng/tg"
	"github.com/jroimartin/tgbot-ng/bing"
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

func (cmd *cmdBing) Run(chatID int, text string) error {
	var (
		filename string
		data     []byte
		err      error
		query    string
	)

	m := cmd.re.FindStringSubmatch(text)
	if len(m) == 2 {
		query = m[1]
	}

	filename, data, err = cmd.search(query)
	if err != nil {
		return err
	}

	if path.Ext(filename) == ".gif" {
		_, err = cmd.cli.SendDocument(chatID, filename, data)
	} else {
		_, err = cmd.cli.SendPhoto(chatID, filename, data)
	}
	return err
}

func (cmd *cmdBing) search(query string) (filename string, data []byte, err error) {
	c := bing.NewClient(cmd.config.ApiKey)
	if cmd.config.SearchLimit > 0 {
		c.Limit = cmd.config.SearchLimit
	}

	results, err := c.Query(bing.Image, query)
	if err != nil {
		return "", nil, err
	}
	if len(results) == 0 {
		return "", nil, errors.New("no pics")
	}
	rndInt := rand.Intn(len(results))

	resp, err := http.Get(results[rndInt].MediaUrl)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	imgData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	return results[rndInt].MediaUrl, imgData, nil
}
