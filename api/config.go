package api

import (
	"io/ioutil"

	"github.com/starkandwayne/covalence/db"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	DatabaseType string `yaml:"database_type"`
	DatabaseDSN  string `yaml:"database_dsn"`

	Addr string `yaml:"listen_addr"`

	WebRoot string `yaml:"web_root"`
}

func (a *Api) ReadConfig(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return err
	}

	if config.Addr == "" {
		config.Addr = ":8888"
	}

	if config.WebRoot == "" {
		config.WebRoot = "/usr/share/covalence/webui"
	}

	if a.Database == nil {
		a.Database = &db.DB{}
	}

	a.Database.Driver = config.DatabaseType
	a.Database.DSN = config.DatabaseDSN

	ws := WebServer{
		Database: a.Database.Copy(),
		Addr:     config.Addr,
		WebRoot:  config.WebRoot,
		Api:      a,
	}
	a.Web = &ws
	return nil
}
