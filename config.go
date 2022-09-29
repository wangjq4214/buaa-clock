package buaaclock

type Config struct {
	LoginURL string
	InfoURL  string
	SaveURL  string

	Retry int

	UserName string
	Password string
}

var defaultConfig = Config{
	LoginURL: "https://app.buaa.edu.cn/uc/wap/login/check",
	InfoURL:  "https://app.buaa.edu.cn/buaaxsncov/wap/default/get-info",
	SaveURL:  "https://app.buaa.edu.cn/buaaxsncov/wap/default/save",
	Retry:    10,
}

func configDefault(configs ...Config) Config {
	if len(configs) == 0 {
		return defaultConfig
	}

	cfg := configs[0]

	if cfg.LoginURL != "" {
		cfg.LoginURL = defaultConfig.LoginURL
	}

	if cfg.InfoURL != "" {
		cfg.InfoURL = defaultConfig.InfoURL
	}

	if cfg.SaveURL != "" {
		cfg.SaveURL = defaultConfig.SaveURL
	}

	if cfg.Retry != 0 {
		cfg.Retry = defaultConfig.Retry
	}

	return cfg
}
