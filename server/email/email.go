package email

import (
	"context"
	"errors"
	"time"

	cfg "github.com/digisan/go-config"
	lk "github.com/digisan/logkit"
	"github.com/mailgun/mailgun-go/v4"
)

var (
	domain = ""
	apiKey = ""
	sender = ""
	mg     *mailgun.MailgunImpl
)

func init() {

	lk.Log("starting...email")

	if err := cfg.Init("email", false, "email-config.json"); err == nil {
		domain = cfg.Val[string]("domain")
		apiKey = cfg.Val[string]("apikey")
		sender = cfg.Val[string]("sender")
	}

	lk.FailOnErrWhen(len(domain) == 0, "%v", errors.New("email domain is empty"))
	lk.FailOnErrWhen(len(apiKey) == 0, "%v", errors.New("email apiKey is empty"))
	lk.FailOnErrWhen(len(sender) == 0, "%v", errors.New("email sender is empty"))
	mg = mailgun.NewMailgun(domain, apiKey)
}

func SetSender(s string) {
	sender = s
}

func Send(recipient, subject, body string) (string, string, error) {

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	msg, id, err := mg.Send(ctx, message)

	if err != nil {
		lk.Warn("%v", err)
		return "", "", err
	}

	lk.Log("ID: %s Resp: %s\n", id, msg)
	return msg, id, nil
}
