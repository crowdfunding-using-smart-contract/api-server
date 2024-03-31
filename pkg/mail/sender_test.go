package mail_test

import (
	"fund-o/api-server/config"
	"fund-o/api-server/pkg/mail"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EmailSenderSuite struct {
	suite.Suite
	appConfig config.AppConfig
}

func (s *EmailSenderSuite) SetupTest() {
	var err error
	s.appConfig, err = config.LoadAppConfig("../..")
	require.NoError(s.T(), err)
}

func (s *EmailSenderSuite) TestSendEmailWithGmail() {
	gmailOpts := mail.GmailSenderOptions{
		Name:              s.appConfig.EMAIL_SENDER_NAME,
		FromEmailAddress:  s.appConfig.EMAIL_SENDER_ADDRESS,
		FromEmailPassword: s.appConfig.EMAIL_SENDER_PASSWORD,
	}
	sender := mail.NewGmailSender(&gmailOpts)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="https://bizzy.cool">DIZZY</a></p>
	`
	to := []string{"someemail@email.com"}

	err := sender.SendEmail(subject, content, to, nil, nil)
	require.NoError(s.T(), err)
}

func TestEmailSenderSuite(t *testing.T) {
	suite.Run(t, new(EmailSenderSuite))
}
