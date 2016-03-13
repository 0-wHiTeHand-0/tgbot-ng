package main

import (
	"github.com/jroimartin/tgbot-ng/tg"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"regexp"
	"strings"
)

type cmd4cdg struct {
	re      *regexp.Regexp
	f_slice []string
	config  CmdConfig4cdg
	cli     *tg.Client
}

type CmdConfig4cdg struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

func NewCmd4cdg(config CmdConfig4cdg, cli *tg.Client) Command {
	if config.Path == "" {
		config.Path = "cards"
	}
	return &cmd4cdg{
		re:      regexp.MustCompile(`^/4cdg(?:$|@[a-zA-Z0-9_]+bot$| rules$)`),
		f_slice: []string{},
		config:  config,
		cli:     cli,
	}
}

func (cmd *cmd4cdg) Match(text string) bool {
	return cmd.re.MatchString(text)
}

func (cmd *cmd4cdg) Run(chatID, replyID int, text string, from string, reply_ID *tg.Message) error {
	var (
		err error
	)
	if strings.Contains(text, " rules") {
		_, err = cmd.cli.SendText(chatID, "'4chan drinking card game' rules\n\n1. The left of the phone owner starts.\n2. Players take a card when it's their turn. They must do what the card says.\n3. You win the game when everyone else pass' out.\n\nCard types:\nAction: This is a standard 'do what it says' card.\nInstant: This card may be kept and used at anytime in the game.\nMandatory: Everyone must play this card.\nStatus: This is constant for the whole game or the timeframe indicated on the card.")
		return err
	}
	if len(cmd.f_slice) < 3 {
		cmd.f_slice, err = filepath.Glob(cmd.config.Path + "/*.jpg")
		if cmd.f_slice == nil {
			return err
		}
		temp, err := filepath.Glob(cmd.config.Path + "/*.png") //Cochinada maxima, pero funciona.
		if cmd.f_slice == nil {
			return err
		}
		cmd.f_slice = append(cmd.f_slice, temp...)
		_, err = cmd.cli.SendText(chatID, "Cards shuffled. A new game has begun...")
	}
	rndInt := rand.Intn(len(cmd.f_slice))
	img, err := ioutil.ReadFile(cmd.f_slice[rndInt])
	if err != nil {
		return err
	}
	imgName := strings.Split(cmd.f_slice[rndInt], "/")
	log.Println("SENT -> " + cmd.f_slice[rndInt])
	_, err = cmd.cli.SendPhoto(chatID, tg.File{Name: imgName[len(imgName)-1], Data: img})
	cmd.f_slice = append(cmd.f_slice[:rndInt], cmd.f_slice[rndInt+1:]...)
	return err
}
