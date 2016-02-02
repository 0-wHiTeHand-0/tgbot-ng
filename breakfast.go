package main

import (
	"github.com/jroimartin/tgbot-ng/tg"
	"regexp"
	"strings"
)

type cmdBreak struct {
	re      *regexp.Regexp
	config  CmdConfigBreak
	f_break []string
	cli     *tg.Client
}

type CmdConfigBreak struct {
	Enabled bool  `json:"enabled"`
	Allowed []int `json:"allowed"`
}

func NewCmdBreak(config CmdConfigBreak, cli *tg.Client) Command {
	return &cmdBreak{
		re:     regexp.MustCompile(`^/break(?:$|@[a-zA-Z0-9_]+bot$| .+$)`),
		config: config,
		cli:    cli,
	}
}

func (cmd *cmdBreak) Match(text string) bool {
	return cmd.re.MatchString(text)
}

func (cmd *cmdBreak) Run(chatID, replyID int, text string, from string) error {
	//Compruebo que chatID este permitido
	flag := false
	for i := 0; i < len(cmd.config.Allowed); i++ {
		if cmd.config.Allowed[i] == chatID {
			flag = true
			break
		}
	}
	if flag == false {
		cmd.cli.SendText(chatID, "You don't have power here, motherfucker!")
		return nil
	}

	m := cmd.re.FindStringSubmatch(text)
	m = strings.SplitN(m[0], " ", 2)
	var message string
	if len(m) == 1 {
		message = "<-- Today's breakfast! -->\r\n"
		for i := range cmd.f_break {
			message += cmd.f_break[i] + "\r\n"
		}
		if len(cmd.f_break) > 0 {
			message = message[:len(message)-2]
		} else {
			message = "Nobody wants breakfast for now :("
		}
	} else if len(m) == 2 && m[1] == "-" {
		cmd.f_break = cmd.f_break[:0]
		message = "Breakfast list cleared!"
	} else {
		cmd.f_break = append(cmd.f_break, from+": "+m[1])
		message = from + " wants " + m[1]
	}
	cmd.cli.SendText(chatID, message)
	return nil
}
