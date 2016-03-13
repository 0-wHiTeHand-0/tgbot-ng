package main

import (
	"github.com/jroimartin/tgbot-ng/tg"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

type cmdQuote struct {
	re       *regexp.Regexp
	config   CmdConfigQuote
	f_quotes []string
	cli      *tg.Client
}

type CmdConfigQuote struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
	Allowed []int  `json:"allowed"`
}

func NewCmdQuote(config CmdConfigQuote, cli *tg.Client) Command {
	if config.Path == "" {
		config.Path = "quotes.txt"
	}
	f, err := ioutil.ReadFile(config.Path) //Si es un fichero grande, se puede liar. En este caso es de mas o menos 3k. No problemo.
	if err != nil {
		log.Fatalln(err)
	}
	return &cmdQuote{
		re:       regexp.MustCompile(`^/quote(?:$|@[a-zA-Z0-9_]+bot$| .+$)`),
		config:   config,
		f_quotes: strings.Split(string(f), "\n"),
		cli:      cli,
	}
}

func (cmd *cmdQuote) Match(text string) bool {
	return cmd.re.MatchString(text)
}

func (cmd *cmdQuote) Run(chatID, replyID int, text string, from string, reply_ID *tg.Message) error {
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
	m = strings.SplitN(m[0], " ", 3)
	if len(m) == 1 {
		rndInt := rand.Intn(len(cmd.f_quotes))
		cmd.cli.SendText(chatID, "<-- Random quote -->\r\n"+cmd.f_quotes[rndInt])
	} else if len(m) == 3 && m[1] == ">" {
		re := regexp.MustCompile(`[\r\n]+`)
		quote := re.ReplaceAllString(m[2], " ")
		linesFiltered := make([]string, 0)
		for _, line := range cmd.f_quotes {
			if strings.Contains(strings.ToLower(line), strings.ToLower(quote)) {
				linesFiltered = append(linesFiltered, line)
			}
		}
		if len(linesFiltered) == 0 {
			cmd.cli.SendText(chatID, "No quote found.")
		} else {
			rndInt := rand.Intn(len(linesFiltered))
			cmd.cli.SendText(chatID, "<-- Match! -->\r\n"+linesFiltered[rndInt])
		}
	} else if len(m) == 3 && m[1] == "<" {
		re := regexp.MustCompile(`[\r\n]+`)
		quote := re.ReplaceAllString(m[2], " ")

		f, err := os.OpenFile(cmd.config.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
		if _, err = f.WriteString(quote + "\n"); err != nil {
			return err
		}
		f.Close()

		cmd.f_quotes = append(cmd.f_quotes, quote)
		cmd.cli.SendText(chatID, "<-- We have a new quote! -->\r\n"+quote)
	} else {
		cmd.cli.SendText(chatID, "Nope. Command error.")
	}
	return nil
}
