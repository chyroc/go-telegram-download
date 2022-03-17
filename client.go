package telegram_download

import (
	"os/user"
	"strings"
	"sync"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type Option struct {
	AppID        int
	AppHash      string
	BotToken     string
	AllowChannel map[int64]bool
	SaveDir      string
	Logger       *zap.Logger
}

type Client struct {
	appID        int
	appHash      string
	botToken     string
	allowChannel map[int64]bool
	saveDir      string
	logger       *zap.Logger

	dispatcher tg.UpdateDispatcher
	lock       sync.Mutex
}

func New(opt *Option) *Client {
	if opt.SaveDir == "" {
		home, _ := user.Current()
		opt.SaveDir = home.HomeDir + "/Downloads"
	}
	if !strings.HasSuffix(opt.SaveDir, "/") {
		opt.SaveDir += "/"
	}
	return &Client{
		appID:        opt.AppID,
		appHash:      opt.AppHash,
		botToken:     opt.BotToken,
		allowChannel: opt.AllowChannel,
		saveDir:      opt.SaveDir,
		logger:       opt.Logger,
		// dispatcher:   tg.UpdateDispatcher{},
		// lock:         sync.Mutex{},
	}
}
