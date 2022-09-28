package buaaclock

type Config struct {
	LoginURL string
	InfoURL  string
	SaveURL  string

	Retry int

	UserName string
	Password string
}

func overwriteConfig(configs ...Config) Config {
	defaultConfig := Config{
		LoginURL: "https://app.buaa.edu.cn/uc/wap/login/check",
		InfoURL:  "https://app.buaa.edu.cn/buaaxsncov/wap/default/get-info",
		SaveURL:  "https://app.buaa.edu.cn/buaaxsncov/wap/default/save",
		Retry:    0,
	}

	if len(configs) == 0 {
		return defaultConfig
	}

	cfg := configs[0]

	if cfg.LoginURL != "" {
		defaultConfig.LoginURL = cfg.LoginURL
	}

	if cfg.InfoURL != "" {
		defaultConfig.InfoURL = cfg.InfoURL
	}

	if cfg.SaveURL != "" {
		defaultConfig.SaveURL = cfg.SaveURL
	}

	if cfg.Retry != 0 {
		defaultConfig.Retry = cfg.Retry
	}

	if cfg.UserName != "" {
		defaultConfig.UserName = cfg.UserName
	}

	if cfg.Password != "" {
		defaultConfig.Password = cfg.Password
	}

	return defaultConfig
}
