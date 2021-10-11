package config

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
)

// Config defines the configuration structure.
type Config struct {
	General struct {
		LogLevel  int    `mapstructure:"log_level"`
		LogPath   string `mapstructure:"log_path"`
		Listen    string `mapstructure:"listen"`
	}

	MySQL struct {
		Database    string   `mapstructure:"database"`
		Host        string   `mapstructure:"host"`
		Port        string   `mapstructure:"port"`
		User        string   `mapstructure:"user"`
		Password    string   `mapstructure:"password"`
		Automigrate bool     `mapstructure:"automigrate"`
		MaxIdle     int      `mapstructure:"max_idle"`
		MaxOpen     int      `mapstructure:"max_open"`
		DB          *gorm.DB `mapstructure:"db"`
		Timezone    string   `mapstructure:"timezone"`
	} `mapstructure:"mysql"`

	Redis struct {
		URL       string `mapstructure:"url"`
		MaxIdle   int    `mapstructure:"max_idle"`
		MaxActive int    `mapstructure:"max_active"`
		Pool      *redis.Pool
	} `mapstructure:"redis"`

	Mongo struct {
		Conn     string `mapstructure:"conn"`
		UserName string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
		Frame    string `mapstructure:"frame"`
		Index    bool   `mapstructure:"index"`
		Client   *mongo.Client
	} `mapstructure:"mongo"`

	AssetManager struct {
		Http string `mapstructure:"http"`
		Ak   string `mapstructure:"ak"`
	} `mapstructure:"asset_manager"`

	WeiXinApi struct {
		Http      string `mapstructure:"http"`
		GrantType string `mapstructure:"grant_type"`
	} `mapstructure:"weixin_api"`
}

// C holds the global configuration.
var C Config

