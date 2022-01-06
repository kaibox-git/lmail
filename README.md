# lmail

Emailing with optional limiter. It will send no more then limiterMax emails during limiterPeriod if you set WithLimiter = true. An emails with `WithLimiter = false` (or just skip this parameter) have no limitations.

## Install

```
go get github.com/kaibox-git/lmail
```

## Usage

You can use any email client with wrapper function to implement the lmail.EmailProvider interface:

```go
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
```

The example below uses default smtpclient.

```
go get github.com/kaibox-git/lmail
go get github.com/kaibox-git/lmail/smtpclient
```


```go
import (
    "github.com/kaibox-git/lmail"
    "github.com/kaibox-git/lmail/smtpclient"
)

...

host := `localhost`
port := 25
connTimeout := time.Second // for local smtp server
/*
Only 20 emails per 30 minutes. The rest is ignored.
This is useful for notifications of errors, but has a limitation if emailing is too often.
In this case keep logging info to file.
*/
emailNumber := 20 
period := 30 * time.Minute
emailSender, err := smtpclient.New(host, port, connTimeout, emailNumber, period)
if err != nil {
    println(err.Error())
    os.Exit(1)
}

// plain text with limiter
emailSender.Send(&lmail.Data{
    From: mail.Address{
        Name: `Robot`,
        Address: `robot@domain.com`,
    },
    To: []mail.Address{
            {
                Name: `Test address`,
                Address: `test@domain.com`,
            },
        },
    Subject: `test subject`,
    Body:    `test message`,
    WithLimiter: true,
})

// html body for 2 addresses with no limiter
body := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Emailing HTML format</title>
</head>
<body>
    <h1>Test HTML format</h1>
    <p>This is a test body.</p>
</body>
</html>`

emailSender.Send(&lmail.Data{
    From: mail.Address{
        Name: `Robot`,
        Address: `robot@domain.com`,
    },
    To: []mail.Address{
            {
                Name: `Test address`,
                Address: `test@domain.com`,
            },
            {
                Name: `Test address 2`,
                Address: `test2@domain.com`,
            },
        },
    Subject: `test subject`,
    Body:    body,
})
```

### Using [go-simple-mail](https://github.com/xhit/go-simple-mail) client

```
go get github.com/kaibox-git/lmail
go get github.com/kaibox-git/lmail/gosimpleclient
```

```go
import (
    "github.com/kaibox-git/lmail"
    "github.com/kaibox-git/lmail/gosimpleclient"
)

...

host := `localhost`
port := 25
connTimeout := time.Second // for local smtp server
/*
Only 20 emails per 30 minutes. The rest is ignored.
This is useful for notifications of errors, but has a limitation if emailing is too often.
In this case keep logging info to file.
*/
emailNumber := 20 
period := 30 * time.Minute
emailSender, err := gosimpleclient.New(host, port, connTimeout, emailNumber, period)
if err != nil {
    println(err.Error())
    os.Exit(1)
}

// plain text with limiter
emailSender.Send(&lmail.Data{
    From: mail.Address{
        Name: `Robot`,
        Address: `robot@domain.com`,
    },
    To: []mail.Address{
            {
                Name: `Test address`,
                Address: `test@domain.com`,
            },
        },
    Subject: `test subject`,
    Body:    `test message`,
    WithLimiter: true,
})

// html body for 2 addresses with no limiter
body := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Emailing HTML format</title>
</head>
<body>
    <h1>Test HTML format</h1>
    <p>This is a test body.</p>
</body>
</html>`

emailSender.Send(&lmail.Data{
    From: mail.Address{
        Name: `Robot`,
        Address: `robot@domain.com`,
    },
    To: []mail.Address{
            {
                Name: `Test address`,
                Address: `test@domain.com`,
            },
            {
                Name: `Test address 2`,
                Address: `test2@domain.com`,
            },
        },
    Subject: `test subject`,
    Body:    body,
})
```

### Wrapper function example:

```go
type SomeHandler struct {
    ...
    email lmail.EmailProvider
}

// wrapper function:
func (h *SomeHandler) ErrorMail(message string) error {
    return h.email.Send(&lmail.Data{
        From: mail.Address{
            Name: `Robot`,
            Address: `robot@domain.com`,
        },
        To: []mail.Address{
                {
                    Name: `Test address`,
                    Address: `test@domain.com`,
                },
            },
        Subject: `Error occured`,
        Body:    message,
        WithLimiter: true,
    })
}
```
Now use:

```go
if err := someHandler.ErrorMail(`This is a test body`); err != nil {
    // log it
}
```