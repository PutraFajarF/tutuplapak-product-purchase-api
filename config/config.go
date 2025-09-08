package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App        `yaml:"app"`
		HTTPServer `yaml:"httpserver"`
		HTTPClient `yaml:"httpclient"`
		POSTGRESQL `yaml:"postgresql"`
		Log        `yaml:"logger"`
		Cipher     `yaml:"cipher"`
	}

	// App -.
	App struct {
		Name        string `env-required:"true" yaml:"app_name"    env:"APP_NAME"`
		Version     string `env-required:"true" yaml:"app_version" env:"APP_VERSION"`
		Environment string `yaml:"app_environment" env:"ENVIRONMENT"`
		BaseDir     string `yaml:"app_base_dir" env:"APP_BASE_DIR"`
		TimeZone    string `yaml:"app_time_zone"`
	}

	// HTTP -.
	HTTPServer struct {
		Port    string `env-required:"true" yaml:"httpserver_port" env:"HTTP_PORT"`
		UseSSL  bool   `yaml:"httpserver_use_ssl"`
		SSLKey  string `yaml:"httpserver_ssl_key"`
		SSLCert string `yaml:"httpserver_ssl_cert"`
	}

	// HTTP CLIENT -.
	HTTPClient struct {
		MaxIdleConns        int  `env-required:"true" yaml:"httpc_max_idle_conns"`
		MaxIdleConnsPerHost int  `env-required:"true" yaml:"httpc_max_idle_conns_per_host"`
		InsecureSkipVerify  bool `env-required:"true" yaml:"httpc_insecure_skip_verify"`
		SetTimeOut          int  `env-required:"true" yaml:"httpc_set_timeout"`
		UseClientSSL        bool `yaml:"httpc_use_client_ssl"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
		Path  string `yaml:"log_path"`
	}

	// POSTGRESQL struct -.
	POSTGRESQL struct {
		PostgresqllDriverName string `env-required:"true" yaml:"postgresql_driver_name"`
		MaxOpenConns          int    `env-required:"true" yaml:"postgresql_max_open_conns"`
		MaxIdleConns          int    `env-required:"true" yaml:"postgresql_max_idle_conns"`
		MaxLifetimeConns      int    `env-required:"true" yaml:"postgresql_max_lifetime_conns"`
		LifetimeConnsUnit     string `yaml:"postgresql_lifetime_conns_unit"`
		URL                   string `env-required:"true" yaml:"postgresql_url" env:"POSTGRESQL_URL"`
	}

	Cipher struct {
		CipherPostgresql string `env-required:"true" yaml:"cipher_postgresql" env:"CIPHER_POSTGRESQL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	baseDir := ""
	err := cleanenv.ReadConfig(baseDir+"./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
