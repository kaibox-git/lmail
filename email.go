package lmail

import (
	"encoding/base64"
	"errors"
	"net/mail"
	"strings"
)

type EmailProvider interface {
	Send(data *Data) error
}

type Data struct {
	From        mail.Address
	To          []mail.Address
	Subject     string
	Body        string
	WithLimiter bool
}

func (d *Data) Validate() error {
	if d.Subject == `` {
		return errors.New("empty email subject")
	}
	if d.Body == `` {
		return errors.New("empty email body")
	}
	return nil
}

func (d *Data) MakeCc() string {
	var out = make([]string, 0, len(d.To))
	for i := 1; i < len(d.To); i++ {
		out = append(out, d.To[i].String())
	}
	return strings.Join(out, `, `)
}

func EncodeStr(str string) string {
	if len(str) == 0 {
		return ``
	}
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(str)) + "?="
}
