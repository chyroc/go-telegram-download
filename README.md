# go-telegram-download
Download Telegram File by MTProto(Support Big File).

## Installation

```shell
go get github.com/chyroc/go-telegram-download
```

## Prerequisites

- how to get bot token:
  - create bot: https://core.telegram.org/bots#creating-a-new-bot
- how to get app id and app hash:
  - create app: https://core.telegram.org/api/obtaining_api_id

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/chyroc/go-telegram-download"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	r := telegram_download.New(&telegram_download.Option{
		AppID:        123,
		AppHash:      "xxx",
		BotToken:     "456:789",
		SaveDir:      "/Users/some/Downloads",
	})
	fmt.Println(r.Listen(context.Background()))
}

```