package config

import (
	_ "embed"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

//go:embed default.env
var defaultConfig []byte

// Config ...
type (
	Config struct {
		ENV         string         `json:"env" mapstructure:"env" validate:"required"`
		LogLevel    string         `json:"log_level" mapstructure:"log_level"`
		ProjectId   string         `json:"project_id" mapstructure:"project_id"`
		GRPCAddress string         `json:"grpc_address" mapstructure:"grpc_address" validate:"required"`
		Port        string         `json:"port" mapstructure:"port"`
		Tracer      Tracer         `json:"tracer" mapstructure:"tracer"`
		MySQL       DatabaseConfig `json:"mysql" mapstructure:"mysql"`
		GCSBulk     GCS            `json:"gcs_bulk_disbursement" mapstructure:"gcs_bulk_disbursement"`
		Pubsub      PubSubConfig   `json:"pubsub" mapstructure:"pubsub"`
	}
	GCS struct {
		BucketUri                     string `json:"bucket_uri" mapstructure:"bucket_uri"`
		Bucket                        string `json:"bucket" mapstructure:"bucket" validate:"required"`
		GoogleAccessId                string `json:"google_access_id" mapstructure:"google_access_id"`
		SignedUrlExpiredTimeInMinutes int    `json:"signed_url_expired_time_in_minutes" mapstructure:"signed_url_expired_time_in_minutes"`
	}
	Tracer struct {
		TracerName string `json:"name" mapstructure:"name"`
		Enable     bool   `json:"enable" mapstructure:"enable"`
	}
	DatabaseConfig struct {
		Username           string `json:"username" mapstructure:"username"`
		Password           string `json:"password" mapstructure:"password"`
		Protocol           string `json:"protocol" mapstructure:"protocol"`
		Address            string `json:"address" mapstructure:"address"`
		Database           string `json:"database" mapstructure:"database"`
		MaxIdleConnections int    `json:"max_idle_connections" mapstructure:"max_idle_connections"`
		MaxOpenConnections int    `json:"max_open_connections" mapstructure:"max_open_connections"`
	}
	PubSubConfig struct {
		Type          string `json:"type" mapstructure:"type"`
		ClientEmail   string `json:"client_email" mapstructure:"client_email"`
		InquiryTopic  string `json:"inquiry_topic" mapstructure:"inquiry_topic"`
		TransferTopic string `json:"transfer_topic" mapstructure:"transfer_topic"`
	}
)

// Validate ...
func (o Config) Validate() error {
	validate := validator.New()
	if err := validate.Struct(o); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return err
		}
	}

	return nil
}

var Cfg *Config

func init() {
	Cfg = LoadConfig()
}

// LoadConfig ...
// LoadConfig loads the configuration from the default.env file and environment variables
func LoadConfig() *Config {
	cfg := &Config{}

	// Initialize viper with options
	v := viper.NewWithOptions(viper.KeyDelimiter("__"))

	// Set the config type and file
	v.SetConfigType("env")
	v.SetConfigFile(".env") // Remove this line since you're using default.env, not .env
	v.AddConfigPath("/app/config")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")
	v.SetConfigName("default") // Set the config name to default since your file is named default.env

	// Read the config file into viper
	if err := v.ReadInConfig(); err != nil {
		log.Println("failed to read config from file", err)
	}

	// Unmarshal the config into the cfg struct
	err := v.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	// Validate the config
	if err = cfg.Validate(); err != nil {
		panic(err)
	}

	return cfg
}
