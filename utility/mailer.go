package utility

import (
	"bytes"
	"fmt"
	"html/template"
	"net"
	"strconv"

	"gopkg.in/gomail.v2"
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

	t, _ := template.New("tml").Parse(mailTpl)
	data := tplBind{
		username, uuid,
	}

	buf := bytes.NewBuffer([]byte{})
	t.ExecuteTemplate(buf, "T", &data)

	message := gomail.NewMessage()
	message.SetHeader("From", Postman.sender)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "verify your email")
	message.SetBody("text/html", buf.String())

	host, port, _ := net.SplitHostPort(Postman.smtp)

	_port, _ := strconv.Atoi(port)

	d := gomail.NewDialer(host, _port, Postman.sender, Postman.password)

	if err := d.DialAndSend(message); err != nil {
		fmt.Println(err)
	}
}
