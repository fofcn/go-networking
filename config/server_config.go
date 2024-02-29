package config

type Config struct {
	// Http server config
	HttpServerConfig struct {
		Port           string `env:"SERVER_PORT"`
		Host           string `env:"SERVER_HOST"`
		GinMode        string `env:"GIN_MODE"`
		ReadTimeout    int    `env:"SERVER_READTIMEOUT, default=10"`
		WriteTimeout   int    `env:"SERVER_WRITETIMEOUT, default=10"`
		MaxHeaderBytes int    `env:"SERVER_MAXHEADERBYTES, default=1048576000"`
	} `yaml:"server"`

	// Tcp server config
	TcpServerConfig struct {
		Port string `env:"TCP_SERVER_PORT"`
		Host string `env:"TCP_SERVER_HOST"`
	}

	// application config
	AppConfig struct {
		StorePath          string `env:"FILE_STORE_PATH"`
		CredentialFilePath string `env:"CREDENTIAL_FILE_PATH"`
	}
}

var ApplicationConfig Config

func GetHttpServerConfig() *struct {
	Port           string `env:"SERVER_PORT"`
	Host           string `env:"SERVER_HOST"`
	GinMode        string `env:"GIN_MODE"`
	ReadTimeout    int    `env:"SERVER_READTIMEOUT, default=10"`
	WriteTimeout   int    `env:"SERVER_WRITETIMEOUT, default=10"`
	MaxHeaderBytes int    `env:"SERVER_MAXHEADERBYTES, default=1048576000"`
} {
	return &ApplicationConfig.HttpServerConfig
}

func GetAppStorePath() string {
	return ApplicationConfig.AppConfig.StorePath
}
