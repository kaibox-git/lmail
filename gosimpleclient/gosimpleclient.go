package gosimpleclient

import (
	"errors"
	"strings"
	"time"

	"github.com/kaibox-git/lmail"
	gosimplemail "github.com/xhit/go-simple-mail/v2"
)

type SMTP struct {
	host        string
	port        int
	connTimeout time.Duration
	limiter     *lmail.Limiter
}

func New(host string, port int, connTimeout time.Duration, limiterMax int, limiterPeriod time.Duration) (*SMTP, error) {
	if host == `` {
		return nil, errors.New(`mail host is empty`)
	}
	if port == 0 {
		return nil, errors.New(`mail port is 0`)
	}
	if limiterMax == 0 {
		limiterMax = 20 // default: 20 emails per 30 minutes
	}
	if limiterPeriod == 0 {
		limiterPeriod = 30 * time.Minute // default: 20 emails per 30 minutes
	}
	return &SMTP{
		host:        host,
		port:        port,
		connTimeout: connTimeout,
		limiter:     lmail.NewLimiter(limiterMax, limiterPeriod),
	}, nil
}

func (s *SMTP) Send(data *lmail.Data) error {
	if len(data.To) == 0 {
		return nil
	}

	if data.WithLimiter {
		if s.limiter.IsOn() { // ignore the email if limiter has switched to "On"
			return nil
		} else {
			s.limiter.Add()
		}
	}

	client := gosimplemail.NewSMTPClient()

	//SMTP Client
	client.Host = s.host
	client.Port = s.port
	client.ConnectTimeout = s.connTimeout
	client.SendTimeout = s.connTimeout
	client.Encryption = gosimplemail.EncryptionNone

	//Connect to client
	smtpClient, err := client.Connect()
	if err != nil {
		return err
	}

	//Create the email message
	email := gosimplemail.NewMSG()

	var toS []string
	for i := range data.To {
		toS = append(toS, data.To[i].String())
	}

	email.SetFrom(data.From.String()).AddTo(toS...).SetSubject(data.Subject)
	if email.Error != nil {
		return err
	}

	if strings.Contains(data.Body, `</`) {
		email.SetBody(gosimplemail.TextHTML, data.Body)
	} else {
		email.SetBody(gosimplemail.TextPlain, data.Body)
	}
	if email.Error != nil {
		return err
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}
