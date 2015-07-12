// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/jroimartin/tgbot-ng/tg"
)

const picsURL = "http://ano.lolcathost.org/pics/"

type cmdAno struct {
	re     *regexp.Regexp
	config CmdConfigAno
	cli    *tg.Client
}

type CmdConfigAno struct {
	Enabled     bool `json:"enabled"`
	SearchLimit int  `json:"search_limit"`
}

func NewCmdAno(config CmdConfigAno, cli *tg.Client) Command {
	if config.SearchLimit < 1 {
		config.SearchLimit = 10
	}
	return &cmdAno{
		re:     regexp.MustCompile(`^/ano($| .+$)`),
		config: config,
		cli:    cli,
	}
}

func (cmd *cmdAno) Match(text string) bool {
	return cmd.re.MatchString(text)
}

func (cmd *cmdAno) Run(chatID int, text string) error {
	var (
		filename string
		data     []byte
		err      error
	)

	tags := strings.TrimSpace(strings.TrimPrefix(text, "/ano"))
	if tags == "" {
		filename, data, err = cmd.randomPic()
	} else {
		filename, data, err = cmd.searchTag(strings.Split(tags, ","))
	}
	if err != nil {
		return err
	}

	if _, err := cmd.cli.SendMessage(chatID, "What has been seen cannot be unseen...\n"); err != nil {
		return err
	}

	if path.Ext(filename) == ".gif" {
		_, err = cmd.cli.SendDocument(chatID, filename, data)
	} else {
		_, err = cmd.cli.SendPhoto(chatID, filename, data)
	}
	return err
}

func (cmd *cmdAno) randomPic() (filename string, data []byte, err error) {
	var respData struct {
		Pic struct {
			ID string
		}
	}

	reqData := struct {
		Method string `json:"method"`
	}{
		"random",
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", nil, err
	}

	resp, err := http.Post("http://ano.lolcathost.org/json/pic.json",
		"application/json", bytes.NewReader(reqBody))
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("HTTP error: %v (%v)", resp.Status, resp.StatusCode)
	}

	repBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	err = json.Unmarshal(repBody, &respData)
	if err != nil {
		return "", nil, err
	}

	imgData, err := download(picsURL + respData.Pic.ID)
	if err != nil {
		return "", nil, err
	}
	return respData.Pic.ID, imgData, nil
}

func (cmd *cmdAno) searchTag(tags []string) (filename string, data []byte, err error) {
	var respData struct {
		Pics []struct {
			ID string
		}
	}

	reqData := struct {
		Method string   `json:"method"`
		Tags   []string `json:"tags"`
		Limit  int      `json:"limit"`
	}{
		"searchRelated",
		tags,
		cmd.config.SearchLimit,
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", nil, err
	}

	resp, err := http.Post("http://ano.lolcathost.org/json/tag.json",
		"application/json", bytes.NewReader(reqBody))
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("HTTP error: %v (%v)", resp.Status, resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return "", nil, err
	}
	if len(respData.Pics) <= 1 {
		return "", nil, errors.New("no pics")
	}

	rndInt := rand.Intn(len(respData.Pics) - 1)
	rndData := respData.Pics[rndInt]

	imgData, err := download(picsURL + rndData.ID)
	if err != nil {
		return "", nil, err
	}
	return rndData.ID, imgData, nil
}
