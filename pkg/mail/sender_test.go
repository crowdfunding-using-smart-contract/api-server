package mail

import (
	"fund-o/api-server/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	appConfig, err := config.LoadAppConfig(".")
	require.NoError(t, err)

	sender := NewGmailSender(&GmailSenderOptions{
		Name:              appConfig.EMAIL_SENDER_NAME,
		FromEmailAddress:  appConfig.EMAIL_SENDER_ADDRESS,
		FromEmailPassword: appConfig.EMAIL_SENDER_PASSWORD,
	})

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="https://bizzy.cool">DIZZY</a></p>
	`
	to := []string{"danzkikii@gmail.com"}

	err = sender.SendEmail(subject, content, to, nil, nil)
	require.NoError(t, err)
}
