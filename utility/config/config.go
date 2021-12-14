package config

type envRecord struct {
	Addr            string   `json:"addr"`
	Port            string   `json:"port"`
	SessionSecret   string   `json:"session_secret"`
	Emails          []string `json:"emails"`
	GoogleAnalytics string   `json:"google_analytics"`
	ServeStatic     bool     `json:"serve_static"`
}

type ConfigStruct map[string]envRecord

var Configuration ConfigStruct

var ENV string

func SetENV(env string) {
	ENV = env
}

func HasUser(email string) bool {
	configObj := Configuration[ENV]

	for _, _email := range configObj.Emails {
		if email == _email {
			return true
		}
	}

	return false
}

func Get() envRecord {
	return Configuration[ENV]
}
