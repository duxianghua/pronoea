package config

// Conf is to store our service configuration

var (
	Cfg Confing
)

type Confing struct {
	Database DBConfig    `yaml:"database"`
	Email    EmailConfig `yaml:"email"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DATABASE_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Database string `yaml:"database" env:"DB_DATABASE"`
}

type EmailConfig struct {
	From     string   `yaml:"from" env:"EMAIL_FROM"`
	Host     string   `yaml:"host" env:"EMAIL_SMARTHOST"`
	Port     int      `yaml:"port" env:"EMAIL_PORT"`
	Username string   `yaml:"username" env:"EMAIL_USERNAME"`
	Password string   `yaml:"password" env:"EMAIL_PASSWORD"`
	Html     string   `yaml:"html" env:"EMAIL_HTML"`
	Subject  string   `yaml:"subject" env:"EMAIL_SUBJECT"`
	Bcc      string   `yaml:"bcc" env:"EMAIL_BCC"`
	Cc       []string `yaml:"cc" env:"EMAIL_CC"`
}

// // NewConf return new Conf instance from env
// func NewConf() (*Confing, error) {
// 	var conf Confing

// 	err := cleanenv.ReadEnv(&conf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &conf, nil
// }

// func InitConfig(path string) (*Confing, error) {
// 	var cfg Confing
// 	var err error
// 	if len(path) > 0 {
// 		err = cleanenv.ReadConfig(path, &cfg)
// 	} else {
// 		err = cleanenv.ReadEnv(&cfg)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	CFG = &cfg
// 	return &cfg, nil
// }
