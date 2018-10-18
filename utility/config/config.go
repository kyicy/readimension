package config

type ConfigStruct map[string]struct {
	Addr          string   `json:"addr"`
	Port          string   `json:"port"`
	SessionSecret string   `json:"session_secret"`
	Emails        []string `json:"emails"`
}

var Configuratiosn ConfigStruct

var ENV string

func SetENV(env string) {
	ENV = env
}

func HasUser(email string) bool {
	configObj := Configuratiosn[ENV]

	for _, _email := range configObj.Emails {
		if email == _email {
			return true
		}
	}

	return false
}
