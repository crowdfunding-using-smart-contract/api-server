package mail

import (
	"fund-o/api-server/config"
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
	gmailOpts := GmailSenderOptions{
		Name:              s.appConfig.EmailSenderName,
		FromEmailAddress:  s.appConfig.EmailSenderAddress,
		FromEmailPassword: s.appConfig.EmailSenderPassword,
	}

	sender := NewGmailSender(&gmailOpts)

	subject := "A test email"
	content := NewVerifyEmailTemplate("http://localhost:3000/verify")
	to := []string{"test-2669dc@test.mailgenius.com"}

	err := sender.SendEmail(subject, content, to, nil, nil)
	require.NoError(s.T(), err)
}

func TestEmailSenderSuite(t *testing.T) {
	suite.Run(t, new(EmailSenderSuite))
}
