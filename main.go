// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
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
	tgbot.loop()
}
