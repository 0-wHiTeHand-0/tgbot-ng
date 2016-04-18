package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/jroimartin/tgbot-ng/tg"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type cmdVoice struct {
	re     *regexp.Regexp
	config CmdConfigVoice
	cli    *tg.Client
}

type CmdConfigVoice struct {
	Enabled      bool   `json:"enabled"`
	Espeak_param string `json:"espeak_param"`
}

func NewCmdVoice(config CmdConfigVoice, cli *tg.Client) Command {
	return &cmdVoice{
		re:     regexp.MustCompile(`^/voice(?:$|@[a-zA-Z0-9_]+bot$| [ a-zÁÉÍÓÚáéíóúñÑA-Z0-9.,?!]+$)`),
		config: config,
		cli:    cli,
	}
}

func (cmd *cmdVoice) Match(text string) bool {
	return cmd.re.MatchString(text)
}

func (cmd *cmdVoice) Run(chatID, replyID int, text string, from tg.User, reply_ID *tg.Message) error {
	var (
		err error
	)
	if len(text) > 145 {
		_, err = cmd.cli.SendText(chatID, "Text too large. Don't touch my cyberbowls.")
	} else {
		m := cmd.re.FindStringSubmatch(text)
		m = strings.SplitN(m[0], " ", 2)
		if len(m) < 2 {
			_, err = cmd.cli.SendText(chatID, "Sometimes silence is golden. Now it's not.")
			return err
		}
		re := regexp.MustCompile(`[\r\n]+`)
		texto := re.ReplaceAllString(m[1], " ")

		var stdout1 []byte
		stdout1, err = exec.Command("espeak", cmd.config.Espeak_param, "--stdout", "-s125", texto).Output() //No he encontrado la manera de hacer un pipe multiple en go
		if err != nil {
			log.Println("Espeak error. If you want to speak, you must install it.")
			return err
		}
		ex := exec.Command("opusenc", "-", "-")
		stdin2, _ := ex.StdinPipe()
		stdout2, _ := ex.StdoutPipe()
		err = ex.Start()
		if err != nil {
			log.Println("Opusenc error. If you want to speak, you must install it.")
			return err
		}
		stdin2.Write(stdout1)
		stdin2.Close()
		grepbytes, _ := ioutil.ReadAll(stdout2)
		ex.Wait()
		hasher := md5.New()
		hasher.Write(grepbytes)
		audio_name := hex.EncodeToString(hasher.Sum(nil)) + ".ogg"
		log.Println("Sent -> " + audio_name)
		_, err = cmd.cli.SendVoice(chatID, tg.File{Name: audio_name, Data: grepbytes})
	}
	return err
}
