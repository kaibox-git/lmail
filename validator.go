package lmail

import (
	"net/mail"
	"regexp"
)

const (
	minEmailLen = 3
	maxEmailLen = 255
)

var emailMask = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// to validate email addresses from config file
func IsValid(mails []mail.Address) string {
	for i := range mails {
		l := len(mails[i].Address)
		if l < minEmailLen || l > maxEmailLen {
			return mails[i].Address
		}
		if !emailMask.MatchString(mails[i].Address) {
			return mails[i].Address
		}
	}
	return ``
}
