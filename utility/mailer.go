package utility

import (
	"bytes"
	"html/template"
	"os/exec"
)

type mailer struct {
	env string
}

var Postman *mailer

func SetUpMailer(env string) {
	Postman = new(mailer)
	Postman.env = env
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

func (m *mailer) Send(username, email, uuid string) {
	if m.env != "production" {
		return
	}
	t, _ := template.New("tml").Parse(mailTpl)
	cmd := exec.Command("mail", "-a", "Content-type: text/html", "-s", "verify your email", email)

	data := tplBind{
		username, uuid,
	}

	buf := bytes.NewBuffer([]byte{})
	t.ExecuteTemplate(buf, "T", &data)

	cmd.Stdin = buf
	cmd.Start()
	cmd.Wait()
}
