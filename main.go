package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"

	"gopkg.in/yaml.v2"
)

type setting struct {
	URL         string            `yaml:"url"`
	Prefix      string            `yaml:"prefix"`
	Resolver    string            `yaml:"resolver"`
	Location    string            `yaml:"location"`
	Auth        int               `yaml:"auth"`
	AuthURL     string            `yaml:"authurl"`
	Free        map[string]string `yaml:"free"`
	AuthService int               `yaml:"auth_service"`
}

func (s *setting) GetResolver() string {
	if s.Resolver == "-" {
		return ""
	}
	return "resolver " + s.Resolver
}

type yamlSetting struct {
	Setting               []*setting        `yaml:"locations"`
	ServerName            string            `yaml:"serverName"`
	ProxySetHeader        map[string]string `yaml:"proxy_set_header"`
	ProxySetHeaderDefault bool              `yaml:"proxy_set_header_default"`
	ServerSection         map[string]string `yaml:"server_section"`
}

var (
	nginxConf             = ""
	proxySetHeaderDefault = map[string]string{
		"Host":              "$host",
		"X-Real-IP":         "$remote_addr",
		"X-Forwarded-Proto": "$http_x_forwarded_proto",
		"X-Forwarded-For":   "$proxy_add_x_forwarded_for",
	}
)

func getenv(key string, defValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defValue
	}
	return value
}

func main() {
	var nginxconf = getenv("NGINX_CONF", "/etc/nginx/conf.d/default.conf")
	_, err := os.Stat(nginxconf)
	if !os.IsNotExist(err) {
		os.Remove(nginxconf)
	}
	if writer, err := os.OpenFile(nginxconf, os.O_RDWR|os.O_CREATE, 0600); err == nil {
		defer writer.Close()
		logic(writer)
		cmd := exec.Command("nginx", "-s", "reload")
		var buf []byte
		var err error
		if buf, err = cmd.Output(); err != nil {
			panic(err)
		}
		log.Println(string(buf))
	} else {
		panic(err)
	}

}

func merge(values ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, mapvalue := range values {
		for key, value := range mapvalue {
			result[key] = value
		}
	}
	return result
}

func logic(writer io.Writer) bool {
	var filename = os.Getenv("SETTING_FILE_NAME")
	var yamlsetting yamlSetting
	var buf []byte
	var err error

	if filename == "" {
		log.Printf("FileName is empty")
		return false
	}
	if buf, err = ioutil.ReadFile(filename); err != nil {
		log.Printf("File Read Error : %s\n", filename)
		return false
	}
	if err = yaml.Unmarshal(buf, &yamlsetting); err != nil {
		log.Printf("File to Yaml Error: %s\n", filename)
		return false
	}
	if yamlsetting.ServerName == "" {
		yamlsetting.ServerName = "_"
	}
	for _, setting := range yamlsetting.Setting {
		if setting.Resolver == "" {
			setting.Resolver = "kube-dns.kube-system valid=2s"
		}
	}

	if yamlsetting.ProxySetHeaderDefault {
		yamlsetting.ProxySetHeader = merge(proxySetHeaderDefault, yamlsetting.ProxySetHeader)
	}
	// tpl := template.Must(template.New("nginxconf").Parse(nginxConf))
	var configTpl = getenv("CONFIG_DIR", "/app/templates/nginxConf.conf")
	tpl := template.Must(template.ParseFiles(configTpl))
	if err = tpl.Execute(writer, yamlsetting); err != nil {
		panic(err)
	}

	return true
}
