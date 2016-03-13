// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: tgbot config")
		os.Exit(2)
	}
	configFile := os.Args[1]

	cfg, err := parseConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	tgbot := newBot(cfg.Name, cfg.Token)
	if len(cfg.AllowedIDs) > 0 {
		tgbot.setAllowedIDs(cfg.AllowedIDs)
	}
	if cfg.UpdateInterval > 0 {
		tgbot.setUpdateInterval(cfg.UpdateInterval)
	}
	// Add enabled commands
	if cfg.Commands.Ano.Enabled {
		tgbot.addCommand(NewCmdAno(cfg.Commands.Ano, tgbot.cli))
	}
	if cfg.Commands.Bing.Enabled {
		tgbot.addCommand(NewCmdBing(cfg.Commands.Bing, tgbot.cli))
	}

	if cfg.Commands.breakfast.Enabled {
		tgbot.addCommand(NewCmdBreak(cfg.Commands.breakfast, tgbot.cli))
	}
	if cfg.Commands.fcdg.Enabled {
		tgbot.addCommand(NewCmd4cdg(cfg.Commands.fcdg, tgbot.cli))
	}
	if cfg.Commands.quote.Enabled {
		tgbot.addCommand(NewCmdQuote(cfg.Commands.quote, tgbot.cli))
	}
	if cfg.Commands.voice.Enabled {
		tgbot.addCommand(NewCmdVoice(cfg.Commands.voice, tgbot.cli))
	}
	if cfg.Commands.chive.Enabled {
		tgbot.addCommand(NewCmdChive(cfg.Commands.chive, tgbot.cli))
	}
	if cfg.Commands.ban.Enabled {
		tgbot.addCommand(NewCmdBan(cfg.Commands.ban, tgbot.cli))
	}
	tgbot.loop()
}

type config struct {
	Name           string     `json:"name"`
	Token          string     `json:"token"`
	UpdateInterval int        `json:"update_interval"`
	AllowedIDs     []int      `json:"allowed_ids"`
	Commands       cmdConfigs `json:"commands"`
}

type cmdConfigs struct {
	Ano       CmdConfigAno   `json:"ano"`
	Bing      CmdConfigBing  `json:"bing"`
	fcdg      CmdConfig4cdg  `json:"fcdg"`
	quote     CmdConfigQuote `json:"quote"`
	voice     CmdConfigVoice `json:"voice"`
	breakfast CmdConfigBreak `json:"breakfast"`
	chive     CmdConfigChive `json:"chive"`
	ban       CmdConfigBan   `json:"ban"`
}

func parseConfig(file string) (config, error) {
	f, err := os.Open(file)
	if err != nil {
		return config{}, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return config{}, err
	}

	var cfg config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return config{}, err
	}
	log.Println(cfg)
	cfg.Commands.fcdg.Enabled = false
	cfg.Commands.fcdg.Path = "cards"
	cfg.Commands.quote.Enabled = false
	cfg.Commands.quote.Path = "quotes.txt"
	cfg.Commands.quote.Allowed = []int{0, 0}
	cfg.Commands.voice.Enabled = false
	cfg.Commands.breakfast.Allowed = []int{0, 0}
	cfg.Commands.breakfast.Enabled = false
	cfg.Commands.chive.Enabled = true
	cfg.Commands.ban.Enabled = true
	cfg.Commands.ban.Allowed = []int{0}
	return cfg, nil
}
