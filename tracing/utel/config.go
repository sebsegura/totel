package utel

type Config struct {
	ServiceName string
	Owner       string
	Flow        string
}

var _cfg *Config

func SetUtelConfig(cfg *Config) {
	_cfg = cfg
}

func GetUtelConfig() *Config {
	return _cfg
}
