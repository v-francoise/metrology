package main

import (
	"metrology/agent"

	"github.com/go-ini/ini"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {

	var (
		configFilepath = kingpin.Arg("config-file", "Configuration filepath").Default("/etc/metrology/metrology.conf").String()
	)

	kingpin.Version("0.0.1")
	kingpin.Parse()

	cfg, err := ini.Load(*configFilepath)

	if err != nil {
		panic(err)
	}

	metricAgentEndpoint := cfg.Section("agent").Key("metric-agent-endpoint").String()
	metricPublicationEndpoint := cfg.Section("agent").Key("metric-publication-endpoint").String()

	meteringAgent := agent.CreateAgent(metricAgentEndpoint, metricPublicationEndpoint)
	meteringAgent.Run()
}
