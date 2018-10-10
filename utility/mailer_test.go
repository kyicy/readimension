package utility

import (
	"bytes"
	"html/template"
	"os"
	"os/exec"
	"testing"
)

func TestTemplate(t *testing.T) {
	type tplBind struct {
		Username string
		UUID     string
	}

	cmd := exec.Command("wc", "-l")

	var mailTpl = `
		{{define "T"}}
			<main>
				<h2>Just one more Step...</h2>
				<h3>{{.Username}}</h3>
				<p>Click the link below to activate your Readimension account</p>
				<a href="https://www.redimension.com/activate/{{.UUID}}">Activate Account</a>
				<p>The Readimension Team</p>
			</main>
		{{end}}`

	tee, _ := template.New("tml").Parse(mailTpl)

	data := tplBind{
		"kyicy", "godlike",
	}

	buf := bytes.NewBuffer([]byte{})
	tee.ExecuteTemplate(buf, "T", &data)

	cmd.Stdin = buf
	cmd.Stdout = os.Stdout

	cmd.Start()
	cmd.Wait()

}
