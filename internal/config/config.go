package config

import (
	"github.com/spf13/viper"
)

// Config holds the top-level application configuration.
type Config struct {
	Docker DockerConfig `yaml:"docker" mapstructure:"docker"`
	SSH    SSHConfig    `yaml:"ssh" mapstructure:"ssh"`
	Nginx  NginxConfig  `yaml:"nginx" mapstructure:"nginx"`
	Env    EnvConfig    `yaml:"env" mapstructure:"env"`
}

// DockerConfig holds Docker-related settings.
type DockerConfig struct {
	Host string `yaml:"host" mapstructure:"host"`
}

// SSHConfig holds SSH manager settings.
type SSHConfig struct {
	ConfigPath string `yaml:"config_path" mapstructure:"config_path"`
}

// NginxConfig holds nginx generator settings.
type NginxConfig struct {
	ConfigDir    string `yaml:"config_dir" mapstructure:"config_dir"`
	TemplatesDir string `yaml:"templates_dir" mapstructure:"templates_dir"`
}

// EnvConfig holds env sync settings.
type EnvConfig struct {
	DefaultEnvFile string `yaml:"default_env_file" mapstructure:"default_env_file"`
}

// Load reads the config from Viper into a Config struct.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
