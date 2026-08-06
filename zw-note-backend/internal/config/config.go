package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Wechat   WechatConfig   `mapstructure:"wechat"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Redemption RedemptionConfig `mapstructure:"redemption"`
	Notes    NotesConfig    `mapstructure:"notes"`
}

// NotesConfig holds settings for the notes module.
type NotesConfig struct {
	DefaultUserID uint64 `mapstructure:"default_user_id"` // MVP 阶段未接登录，所有请求默认归属该用户
}

// RedemptionConfig holds settings for redemption codes.
type RedemptionConfig struct {
	CooldownDays int `mapstructure:"cooldown_days"` // 用户使用兑换码冷却期（天），默认 30
}

// UploadConfig holds settings for local file uploads.
type UploadConfig struct {
	Dir       string `mapstructure:"dir"`         // local directory, e.g. ./uploads
	MaxSizeMB int    `mapstructure:"max_size_mb"` // max file size in MB
}

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`  // seconds
	WriteTimeout int    `mapstructure:"write_timeout"` // seconds
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	DSN             string `mapstructure:"dsn"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // seconds
}

// JWTConfig holds separate secrets and TTLs for the admin and mini-program audiences.
type JWTConfig struct {
	AdminSecret      string `mapstructure:"admin_secret"`
	MPSecret         string `mapstructure:"mp_secret"`
	AdminExpireHours int    `mapstructure:"admin_expire_hours"`
	MPExpireHours    int    `mapstructure:"mp_expire_hours"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
	AllowMethods []string `mapstructure:"allow_methods"`
	AllowHeaders []string `mapstructure:"allow_headers"`
	MaxAge       int      `mapstructure:"max_age"`
}

// WechatConfig supports multiple mini-programs identified by app_key.
type WechatConfig struct {
	Apps map[string]WechatApp `mapstructure:"apps"`
}

// WechatApp holds credentials for one WeChat mini-program.
type WechatApp struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
	Name      string `mapstructure:"name"`
}

// GetApp returns the WechatApp for the given key, or an error if not found.
func (w *WechatConfig) GetApp(appKey string) (WechatApp, error) {
	app, ok := w.Apps[appKey]
	if !ok {
		return WechatApp{}, fmt.Errorf("unknown wechat app key: %s", appKey)
	}
	return app, nil
}

// Load reads the configuration file and environment variables.
// Environment variables override file values; keys use APP_ prefix with dots → underscores.
// e.g. APP_SERVER_PORT=9090 overrides server.port.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
