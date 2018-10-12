package utility

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
)

type mailer struct {
	env      string
	sender   string
	password string
	smtp     string
}

var Postman *mailer

func SetUpMailer(env, sender, password, smtp string) {
	Postman = new(mailer)
	Postman.env = env
	Postman.sender = sender
	Postman.password = password
	Postman.smtp = smtp
}

var mailTpl = `
{{define "T"}}
    <main>
        <h2>Just one more Step...</h2>
        <h3>{{.Username}}</h3>
        <p>Click the link below to activate your Readimension account</p>
        <a href="https://www.readimension.com/activate/{{.UUID}}">Activate Account</a>
        <p>The Readimension Team</p>
    </main>
{{end}}
`

type tplBind struct {
	Username string
	UUID     string
}

func (m *mailer) SendVerification(username, email, uuid string) {
	if m.env != "production" {
		return
	}
	t, _ := template.New("tml").Parse(mailTpl)

	data := tplBind{
		username, uuid,
	}

	buf := bytes.NewBuffer([]byte{})
	t.ExecuteTemplate(buf, "T", &data)

	auth := smtp.PlainAuth("", Postman.sender, Postman.password, Postman.smtp)

	to := []string{email}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"

	msg := []byte(
		"Subject: verify your email\r\n" +
			mime +
			buf.String() + "\r\n")

	err := smtp.SendMail(Postman.smtp+":25", auth, Postman.sender, to, msg)

	if err != nil {
		log.Fatal(err)
	}
}
