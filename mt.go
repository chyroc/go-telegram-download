package telegram_download

import (
	"context"
	"fmt"
	"mime"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func (r *Client) Listen(ctx context.Context) error {
	r.dispatcher = tg.NewUpdateDispatcher()
	client := telegram.NewClient(r.appID, r.appHash, telegram.Options{
		UpdateHandler: r.dispatcher,
		Logger:        r.logger,
	})

	if err := r.setup(ctx, client); err != nil {
		return fmt.Errorf("setup: %w", err)
	}

	fmt.Println("start to run")
	return client.Run(ctx, func(ctx context.Context) error {
		status, err := client.Auth().Status(ctx)
		if err != nil {
			return fmt.Errorf("auth status: %w", err)
		}

		if !status.Authorized {
			if _, err := client.Auth().Bot(ctx, r.botToken); err != nil {
				return fmt.Errorf("login: %w", err)
			}
		}

		return telegram.RunUntilCanceled(ctx, client)
	})
}

func (r *Client) setup(ctx context.Context, client *telegram.Client) error {
	api := tg.NewClient(client)
	sender := message.NewSender(api)

	r.dispatcher.OnNewChannelMessage(func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewChannelMessage) error {
		fmt.Println("got message")
		document, ok := getDocument(u.Message)
		if !ok {
			return nil
		}
		filename := getFilename(document)
		fmt.Println("got file: ", filename)
		builder := downloader.NewDownloader().Download(api, &tg.InputDocumentFileLocation{
			ID:            document.GetID(),
			AccessHash:    document.GetAccessHash(),
			FileReference: document.GetFileReference(),
		})
		_, err := builder.ToPath(context.Background(), r.saveDir+filename)
		if err != nil {
			_, _ = sender.Reply(entities, u).Text(ctx, fmt.Sprintf("[error] save %q to %q: %s", filename, r.saveDir, err))
		} else {
			_, _ = sender.Reply(entities, u).Text(ctx, fmt.Sprintf("[success] save %q to %q", filename, r.saveDir))
		}
		return err
	})
	return nil
}

func (r *Client) isAllowChannel(channelID int64) bool {
	if len(r.allowChannel) == 0 {
		return true
	}
	return r.allowChannel[channelID]
}

func getDocument(msg tg.MessageClass) (*tg.Document, bool) {
	m, ok := msg.(*tg.Message)
	if !ok || m.Out {
		return nil, false
	}
	mm, ok := m.Media.(*tg.MessageMediaDocument)
	if !ok {
		return nil, false
	}
	mmd, ok := mm.Document.(*tg.Document)
	if !ok {
		return nil, false
	}
	return mmd, true
}

func getFilename(document *tg.Document) string {
	name := ""
	for _, attr := range document.Attributes {
		switch attr := attr.(type) {
		case *tg.DocumentAttributeFilename:
			name = attr.FileName
			break
		case *tg.DocumentAttributeAudio:
			name = fmt.Sprintf("%s.%s", attr.Title, mime2ext(document.GetMimeType()))
			break
		}
	}
	if name == "" {
		name = fmt.Sprintf("%d%s", document.GetID(), mime2ext(document.GetMimeType()))
	}
	return name
}

func mime2ext(s string) string {
	res, err := mime.ExtensionsByType(s)
	if err != nil {
		return ""
	}
	if len(res) == 0 {
		return ""
	}
	return res[0]
}
