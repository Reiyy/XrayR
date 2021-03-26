package panel

import (
	"github.com/XrayR-project/XrayR/api"
	"github.com/XrayR-project/XrayR/service/controller"
)

type Config struct {
	LogConfig     *LogConfig     `mapstructure:"Log"`
	DnsConfigPath string         `mapstructure:"DnsConfigPath"`
	NodesConfig   []*NodesConfig `mapstructure:"Nodes"`
}

type NodesConfig struct {
	PanelType        string             `mapstructure:"PanelType"`
	ApiConfig        *api.Config        `mapstructure:"ApiConfig"`
	ControllerConfig *controller.Config `mapstructure:"ControllerConfig"`
}

type LogConfig struct {
	Level      string `mapstructure:"Level"`
	AccessPath string `mapstructure:"AccessPath"`
	ErrorPath  string `mapstructure:"ErrorPath"`
}
