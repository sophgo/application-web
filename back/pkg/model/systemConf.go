package model

type SystemConf struct {
	ServerInfo     Server         `mapstructure:"server"`
	Log            LogConfig      `mapstructure:"log"`
	Db             DbConfig       `mapstructure:"db"`
	MediaHost      string         `mapstructure:"mediahost"`
	Picture        PictureConfig  `mapstructure:"picture"`
	Face           FaceConfig     `mapstructure:"face"`
	AlgorithmHost  string         `mapstructure:"algorithmhost"`
	FaceAlgoHost   string         `mapstructure:"faceAlgohost"`
	UploadHost     string         `mapstructure:"uploadhost"`
	FaceUploadHost string         `mapstructure:"faceUploadHost"`
	LocalHost      string         `mapstructure:"localhost"`
	Abilities      map[int]string `mapstructure:"abilities" yaml:"abilities"`
}

type Server struct {
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Timeout  string `mapstructure:"timeout"`
}
type LogConfig struct {
	Path  string `mapstructure:"path"`
	Level string `mapstructure:"level"`
}

type DbConfig struct {
	Path     string `mapstructure:"path"`
	SaveDays int    `mapstructure:"savedays"`
}

type PictureConfig struct {
	Dir     string `mapstructure:"dir"`
	MaxSize int64  `mapstructure:"maxsize"`
	Quality int64  `mapstructure:"quality"`
}

type FaceConfig struct {
	Dir     string `mapstructure:"dir"`
	MaxSize int64  `mapstructure:"maxsize"`
	Quality int64  `mapstructure:"quality"`
}
