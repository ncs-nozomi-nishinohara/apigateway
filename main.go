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

type Setting struct {
	Url      string `yaml:"url"`
	Prefix   string `yaml:"prefix"`
	Resolver string `yaml:"resolver"`
	Location string `yaml:"location"`
}

func (s *Setting) GetResolver() string {
	if s.Resolver == "-" {
		return ""
	} else {
		return "resolver " + s.Resolver
	}
}

type YamlSetting struct {
	Setting               []*Setting        `yaml:"locations"`
	ServerName            string            `yaml:"serverName"`
	ProxySetHeader        map[string]string `yaml:"proxy_set_header"`
	ProxySetHeaderDefault bool              `yaml:"proxy_set_header_default"`
	ServerSection         map[string]string `yaml:"server_section"`
}

var (
	nginxConf = `server {
  listen 80;
  server_name {{ .ServerName }};
  access_log /dev/stdout;
  error_log /dev/stderr warn;
  server_tokens off;
  {{range $key, $value := .ProxySetHeader -}}
  proxy_set_header {{$key}} {{$value}};
  {{end}}
  {{range $key, $value := .ServerSection -}}
  {{$key}} {{$value}};
  {{end -}}
{{range $index, $setting := .Setting}}
  location {{$setting.Prefix}} {{$setting.Location}} {
    {{$setting.GetResolver}};
    proxy_pass {{$setting.Url}};
  }
{{end}}
}
`
	proxy_set_header_default = map[string]string{
		"Host":              "$host",
		"X-Real-IP":         "$remote_addr",
		"X-Forwarded-Proto": "$http_x_forwarded_proto",
		"X-Forwarded-For":   "$proxy_add_x_forwarded_for",
	}
)

func getenv(key string, default_ string) string {
	value := os.Getenv(key)
	if value == "" {
		return default_
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
	var yamlSetting YamlSetting
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
	if err = yaml.Unmarshal(buf, &yamlSetting); err != nil {
		log.Printf("File to Yaml Error: %s\n", filename)
		return false
	}
	if yamlSetting.ServerName == "" {
		yamlSetting.ServerName = "_"
	}
	for _, setting := range yamlSetting.Setting {
		if setting.Resolver == "" {
			setting.Resolver = "kube-dns.kube-system valid=2s"
		}
	}

	if yamlSetting.ProxySetHeaderDefault {
		yamlSetting.ProxySetHeader = merge(proxy_set_header_default, yamlSetting.ProxySetHeader)
	}
	tpl := template.Must(template.New("nginxconf").Parse(nginxConf))
	if err = tpl.Execute(writer, yamlSetting); err != nil {
		panic(err)
	}

	return true
}
