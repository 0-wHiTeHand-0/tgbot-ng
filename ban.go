package main

import (
	"github.com/jroimartin/tgbot-ng/tg"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type cmdBan struct {
	re     *regexp.Regexp
	config CmdConfigBan
	cli    *tg.Client
}

type CmdConfigBan struct {
	Enabled		bool	`json:"enabled"`
	Allowed		[]int	`json:"allowed"`
	Pre_Ban_ids	[]int	`json:"pre_banned_ids"`
	Pre_Ban_time	int	`json:"pre_banned_time"`
}

func NewCmdBan(config CmdConfigBan, cli *tg.Client) Command {
	for _, i := range config.Pre_Ban_ids{
		cli.BannedIDs[i] = time.Now()
		cli.BannedIDs_min[i] = config.Pre_Ban_time
	}
	return &cmdBan{
		re:     regexp.MustCompile(`^/ban(?: [0-9]+$| -$)`),
		config: config,
		cli:    cli,
	}
}

func (cmd *cmdBan) Match(text string) bool {
	return cmd.re.MatchString(text)
}

func (cmd *cmdBan) Run(chatID, replyID int, text string, from tg.User, reply_ID *tg.Message) error {
	//Compruebo que el ID este permitido
	flag := false
	for i := 0; i < len(cmd.config.Allowed); i++ {
		if cmd.config.Allowed[i] == from.ID {
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
	if m[1] == "-" {
		cmd.cli.BannedIDs = make(map[int]time.Time)
		cmd.cli.BannedIDs_min = make(map[int]int)
		message = "Ban list cleared! Trolls are free again!"
	} else {
		if reply_ID == nil {
			cmd.cli.SendText(chatID, "Error. You should reply a message.")
			return nil
		}
		cmd.cli.BannedIDs[reply_ID.From.ID] = time.Now()
		cmd.cli.BannedIDs_min[reply_ID.From.ID], _ = strconv.Atoi(m[1])
		message = "Ban set to " + reply_ID.From.FirstName + " for " + m[1] + " seconds. Keep calm and relax your boobies, little troll."
	}
	cmd.cli.SendText(chatID, message)
	return nil
}
