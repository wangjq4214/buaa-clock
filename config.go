package buaaclock

type Config struct {
	LoginURL string
	InfoURL  string
	SaveURL  string

	Retry int

	UserName string
	Password string

	// Whether to stay at school, "1" means yes, "0" means no
	Boarder string

	// If you are not staying at school, the reasons you need
	// to choose from "1" to "6" are "Temporary absence from school",
	// "Returning home during summer and winter vacations",
	// "Overseas research and study",
	// "Off-campus business trip or internship",
	// "Sick leave, leave of absence or leave of absence from school",
	// and "Other".
	NotBoarderReasen string

	// The reason for leaving school is other reasons that need to be filled in.
	NotBoarderNote string

	// Clocking BoarderAddress
	BoarderAddress  string
	BoarderArea     string
	BoarderCity     string
	BoarderProvince string

	NotBoarderAddress  string
	NotBoarderArea     string
	NotBoarderCity     string
	NotBoarderProvince string
}

var defaultConfig = Config{
	LoginURL: "https://app.buaa.edu.cn/uc/wap/login/check",
	InfoURL:  "https://app.buaa.edu.cn/buaaxsncov/wap/default/get-info",
	SaveURL:  "https://app.buaa.edu.cn/buaaxsncov/wap/default/save",

	Retry: 10,

	Boarder: "1",

	NotBoarderReasen: "2",
	NotBoarderNote:   "寒暑假回家",

	BoarderAddress:  "北京市海淀区花园路街道北京航空航天大学大运村学生公寓5号楼",
	BoarderArea:     "北京市 海淀区",
	BoarderCity:     "北京市",
	BoarderProvince: "北京市",

	NotBoarderAddress:  "",
	NotBoarderArea:     "",
	NotBoarderCity:     "",
	NotBoarderProvince: "",
}

func configDefault(configs ...Config) Config {
	if len(configs) == 0 {
		return defaultConfig
	}

	cfg := configs[0]

	if cfg.LoginURL == "" {
		cfg.LoginURL = defaultConfig.LoginURL
	}

	if cfg.InfoURL == "" {
		cfg.InfoURL = defaultConfig.InfoURL
	}

	if cfg.SaveURL == "" {
		cfg.SaveURL = defaultConfig.SaveURL
	}

	if cfg.Retry == 0 {
		cfg.Retry = defaultConfig.Retry
	}

	if cfg.Boarder == "" {
		cfg.Boarder = defaultConfig.Boarder
		cfg.BoarderAddress = defaultConfig.BoarderAddress
		cfg.BoarderArea = defaultConfig.BoarderArea
		cfg.BoarderCity = defaultConfig.BoarderCity
		cfg.BoarderProvince = defaultConfig.BoarderProvince
	} else if cfg.Boarder == "0" {
		if cfg.NotBoarderReasen == "" {
			cfg.NotBoarderReasen = "2"
		}

		if cfg.NotBoarderReasen == "6" && cfg.NotBoarderNote == "" {
			cfg.NotBoarderNote = defaultConfig.NotBoarderNote
		}

		if cfg.NotBoarderAddress == "" {
			cfg.NotBoarderAddress = defaultConfig.NotBoarderAddress
		}

		if cfg.NotBoarderArea == "" {
			cfg.NotBoarderArea = defaultConfig.NotBoarderArea
		}

		if cfg.NotBoarderCity == "" {
			cfg.NotBoarderCity = defaultConfig.NotBoarderCity
		}

		if cfg.NotBoarderProvince == "" {
			cfg.NotBoarderProvince = defaultConfig.NotBoarderProvince
		}
	} else {
		cfg.Boarder = defaultConfig.Boarder

		if cfg.BoarderAddress == "" {
			cfg.BoarderAddress = defaultConfig.BoarderAddress
		}

		if cfg.BoarderArea == "" {
			cfg.BoarderArea = defaultConfig.BoarderArea
		}

		if cfg.BoarderCity == "" {
			cfg.BoarderCity = defaultConfig.BoarderCity
		}

		if cfg.BoarderProvince == "" {
			cfg.BoarderProvince = defaultConfig.BoarderProvince
		}
	}

	return cfg
}
