package telegram_download

type Client struct {
	appID   int
	appHash string
}

func New(appID int, appHash string) *Client {
	return &Client{
		appID:   appID,
		appHash: appHash,
	}
}
