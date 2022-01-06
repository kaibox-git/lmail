package smtpclient

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kaibox-git/lmail"
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

	var ContentType string
	if strings.Contains(data.Body, `</`) {
		ContentType = `text/html; charset="utf-8"`
	} else {
		ContentType = `text/plain; charset="utf-8"`
	}

	var (
		msg     bytes.Buffer
		headers = map[string]string{
			`From`:                      data.From.String(),
			`To`:                        data.To[0].String(),
			`Subject`:                   lmail.EncodeStr(data.Subject),
			`Content-Type`:              ContentType,
			`Content-Transfer-Encoding`: `base64`,
			`MIME-Version`:              `1.0`,
		}
	)
	if len(data.To) > 1 {
		headers["Cc"] = data.MakeCc()
	}

	var keys []string
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(&msg, "%s: %s\r\n", key, headers[key])
	}
	msg.WriteString("\r\n")
	msg.WriteString(base64.StdEncoding.EncodeToString([]byte(data.Body)))
	msg.WriteString("\r\n")

	return s.mailProcess(data.From, data.To, msg.Bytes())
}

func (s *SMTP) mailProcess(from mail.Address, to []mail.Address, body []byte) error {
	conn, err := net.DialTimeout("tcp", s.host+`:`+strconv.FormatInt(int64(s.port), 10), s.connTimeout)
	if err != nil {
		return fmt.Errorf("net.DialTimeout() failed: %w", err)
	}
	c, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("net.DialTimeout() failed: %w", err)
	}
	defer func() {
		_ = c.Quit()
	}()
	// Authentication:
	// auth := smtp.PlainAuth("", from, password, smtpHost)
	// c.Auth(auth)
	for i := range to {
		if err = c.Mail(from.Address); err != nil {
			return fmt.Errorf("c.Mail() failed: %w", err)
		}
		if err = c.Rcpt(to[i].Address); err != nil {
			return fmt.Errorf("smtp.Rcpt(%s) failed: %w", to[i], err)
		}
		w, err := c.Data()
		if err != nil {
			return fmt.Errorf("c.Data() failed: %w", err)
		}
		_, err = w.Write(body)
		if err != nil {
			return fmt.Errorf("w.Write() failed: %w", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("w.Close() failed: %w", err)
		}
	}
	return nil
}
