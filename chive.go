// Copyright 2015 The tgbot-ng Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/tgbot-ng/tg"
)

type cmdChive struct {
	re      *regexp.Regexp
	Api_key string
	config  CmdConfigChive
	cli     *tg.Client
}

type CmdConfigChive struct {
	Enabled bool `json:"enabled"`
}

func NewCmdChive(config CmdConfigChive, cli *tg.Client) Command {
	return &cmdChive{
		re:      regexp.MustCompile(`^/chive(?:@[^ ]+?)?(?:$| +(.+)$)`),
		config:  config,
		Api_key: "API_KEY",
		cli:     cli,
	}
}

func (cmd *cmdChive) Match(text string) bool {
	return cmd.re.MatchString(text)
}

func (cmd *cmdChive) Run(chatID, replyID int, text string, from tg.User, reply_ID *tg.Message) error {
	var (
		img  tg.File
		err  error
		tags string
	)

	m := cmd.re.FindStringSubmatch(text)
	if len(m) == 2 {
		tags = m[1]
	}
	if tags == "" {
		img, err = cmd.randomPic()
	} else {
		img, err = cmd.searchTag(strings.Split(tags, ",")[0])
	}
	if err != nil {
		cmd.cli.SendText(chatID, "No photos found!")
		return err
	}

	if _, err := cmd.cli.SendText(chatID, "Keep calm and chive on..."); err != nil {
		return err
	}

	if path.Ext(img.Name) == ".gif" {
		_, err = cmd.cli.SendDocument(chatID, img)
	} else {
		_, err = cmd.cli.SendPhoto(chatID, img)
	}
	return err
}

func (cmd *cmdChive) randomPic() (img tg.File, err error) {
	var category struct {
		Post_Count struct {
			Total_Posts int
		}
		Posts []struct {
			Guid int
		}
	}

	var post struct {
		Posts []struct {
			Items []struct {
				Identity struct {
					URL string
				}
			}
		}
	}

	//Coger numero total de posts  de la categoria 60 'Girls'
	resp, err := http.Get("http://api.thechive.com/api/category/60?key=" + cmd.Api_key)
	if err != nil {
		return tg.File{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return tg.File{}, fmt.Errorf("HTTP error: %v (%v)", resp.Status, resp.StatusCode)
	}
	repBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tg.File{}, err
	}
	err = json.Unmarshal(repBody, &category)
	if err != nil {
		return tg.File{}, err
	}

	//Coger una página aleatoria (cada página tiene 40 posts)
	rand.Seed(time.Now().Unix())
	randPage := rand.Intn(65) // Tomamos como total de posts 2500 (hay 3000 y pico), y se divide por 39 para obtener las paginas
	resp, err = http.Get("http://api.thechive.com/api/category/60?key=" + cmd.Api_key + "&page=" + strconv.Itoa(randPage))
	if err != nil {
		return tg.File{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return tg.File{}, fmt.Errorf("HTTP error: %v (%v)", resp.Status, resp.StatusCode)
	}
	repBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return tg.File{}, err
	}
	err = json.Unmarshal(repBody, &category)
	if err != nil {
		return tg.File{}, err
	}
	fmt.Println("-------\n" + string(repBody) + "\n-------\n")
	//Coger un post aleatorio de la página aleatoria seleccionada
	rand.Seed(time.Now().Unix())
	randPost := rand.Intn(39)
	//fmt.Println("LONGITUD: " + strconv.Itoa(len(category.Posts)))
	if len(category.Posts) < randPost+1 {
		return tg.File{}, errors.New("Posts argument empty!")
	}
	postNum := category.Posts[randPost].Guid

	resp, err = http.Get("http://api.thechive.com/api/post/" + strconv.Itoa(postNum) + "?key=" + cmd.Api_key)
	if err != nil {
		return tg.File{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return tg.File{}, fmt.Errorf("HTTP error: %v (%v)", resp.Status, resp.StatusCode)
	}
	repBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return tg.File{}, err
	}
	err = json.Unmarshal(repBody, &post)
	if err != nil {
		return tg.File{}, err
	}

	//Coger una foto aleatoria del post seleccionado
	rand.Seed(time.Now().Unix())
	randFoto := rand.Intn(len(post.Posts[0].Items))
	picUrl := post.Posts[0].Items[randFoto].Identity.URL

	resp, err = http.Get(picUrl)
	if err != nil {
		return tg.File{}, err
	}
	defer resp.Body.Close()
	imgData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tg.File{}, err
	}

	return tg.File{Name: picUrl, Data: imgData}, nil
}

func (cmd *cmdChive) searchTag(tag string) (img tg.File, err error) {
	var category struct {
		Post_Count struct {
			Total_Posts int
		}
		Posts []struct {
			Guid int
		}
	}

	var post struct {
		Posts []struct {
			Items []struct {
				Identity struct {
					URL string
				}
			}
		}
	}

	//fmt.Println(cmd.Api_key)

	//Coger numero total de posts  de la busqueda'
	resp, err := http.Get("http://api.thechive.com/api/search/" + tag + "?key=" + cmd.Api_key)
	if err != nil {
		return tg.File{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return tg.File{}, fmt.Errorf("HTTP error: %v (%v)", resp.Status, resp.StatusCode)
	}
	repBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tg.File{}, err
	}
	err = json.Unmarshal(repBody, &category)
	if err != nil {
		return tg.File{}, err
	}

	//Coger un post aleatorio de la página aleatoria seleccionada
	rand.Seed(time.Now().Unix())
	if category.Post_Count.Total_Posts == 0 {
		err = errors.New("No photos!")
		return tg.File{}, err
	}
	randPost := rand.Intn(category.Post_Count.Total_Posts)
	postNum := category.Posts[randPost].Guid
	resp, err = http.Get("http://api.thechive.com/api/post/" + strconv.Itoa(postNum) + "?key=" + cmd.Api_key)
	if err != nil {
		return tg.File{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return tg.File{}, fmt.Errorf("HTTP error: %v (%v)", resp.Status, resp.StatusCode)
	}
	repBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return tg.File{}, err
	}
	err = json.Unmarshal(repBody, &post)
	if err != nil {
		return tg.File{}, err
	}

	//Coger una foto aleatoria del post seleccionado
	rand.Seed(time.Now().Unix())
	if len(post.Posts[0].Items) == 0 {
		err = errors.New("No photos!")
		return tg.File{}, err
	}
	randFoto := rand.Intn(len(post.Posts[0].Items))
	picUrl := post.Posts[0].Items[randFoto].Identity.URL

	resp, err = http.Get(picUrl)
	if err != nil {
		return tg.File{}, err
	}
	defer resp.Body.Close()
	imgData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tg.File{}, err
	}

	return tg.File{Name: picUrl, Data: imgData}, nil
}
