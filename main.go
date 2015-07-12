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
	Ano  CmdConfigAno  `json:"ano"`
	Bing CmdConfigBing `json:"bing"`
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
	return cfg, nil
}
