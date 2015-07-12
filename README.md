# tgbot-ng

## Description

Telegram bot

## Usage

```
$ tgbot-ng
usage: tgbot-ng config
```

## Config format

The following snippet shows a typical config file. A
complete example can be found at doc/config.json.

```json
{
	"name": "bot_name",
	"token": "api_token",
	"update_interval": 1,
	"allowed_ids": [],
	"commands": {
		"ano": {
			"enabled": true,
			"search_limit": 10
		},
		...
	}
}
```

## Installation

`go get github.com/jroimartin/tgbot-ng`