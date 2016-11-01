package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/drael/GOnetstat"
	"github.com/starkandwayne/covalence/db"
	"github.com/voxelbrain/goptions"
)

type Config struct {
	ApiURL         string `yaml:"api_url"`
	UpdateInterval int    `yaml:"update_interval"`
	JobName        string `yaml:"job_name"`
	InstanceID     int    `yaml:"instance_id"`
	DeploymentName string `yaml:"deployment_name"`
}

func readConfig(path string) (Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return Config{}, err
	}

	if config.UpdateInterval == 0 {
		config.UpdateInterval = 30
	}

	return config, nil
}

func sendConnections(apiURL string, connections []db.Connection) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(connections)
	_, err := http.Post(apiURL+"/connections", "application/json; charset=utf-8", b)
	return err
}

type CovalanceAgentOpts struct {
	Help       bool   `goptions:"-h, --help, description='Show the help screen'"`
	ConfigFile string `goptions:"-c, --config, description='Path to the agent configuration file'"`
	Log        string `goptions:"-l, --log-level, description='Set logging level to debug, info, notice, warn, error, crit, alert, or emerg'"`
	Version    bool   `goptions:"-v, --version, description='Display the agent version'"`
}

var Version = "(development)"

func main() {
	var opts CovalanceAgentOpts
	opts.Log = "Info"
	if err := goptions.Parse(&opts); err != nil {
		fmt.Printf("%s\n", err)
		goptions.PrintHelp()
		return
	}

	if opts.Help {
		goptions.PrintHelp()
		os.Exit(0)
	}
	if opts.Version {
		if Version == "" {
			fmt.Printf("covalence-agent (development)\n")
		} else {
			fmt.Printf("covalence-agent v%s\n", Version)
		}
		os.Exit(0)
	}

	if opts.ConfigFile == "" {
		fmt.Fprintf(os.Stderr, "No config specified. Please try again using the -c/--config argument\n")
		os.Exit(1)
	}

	config, err := readConfig(opts.ConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read config file\n")
		os.Exit(1)
	}
	for {
		connections := GOnetstat.Tcp()
		filterConnections := []db.Connection{}
		for _, connection := range connections {
			if (connection.ForeignIp != "127.0.0.1") &&
				(connection.Ip != "127.0.0.1") && (connection.State == "ESTABLISHED") {
				conn := db.Connection{
					Source: db.Source{
						IP:          connection.Ip,
						Port:        strconv.Itoa(int(connection.Port)),
						Deployment:  config.DeploymentName,
						Job:         config.JobName,
						Index:       config.InstanceID,
						User:        connection.User,
						Group:       "",
						Pid:         connection.Pid,
						ProcessName: connection.Name,
						Age:         0,
					},
					Destination: db.Destination{
						IP:   connection.ForeignIp,
						Port: strconv.Itoa(int(connection.ForeignPort)),
					},
				}
				filterConnections = append(filterConnections, conn)
			}
		}
		err := sendConnections(config.ApiURL, filterConnections)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to send connections\n")
			os.Exit(1)
		}
		time.Sleep(time.Duration(config.UpdateInterval) * time.Second)
	}
}
