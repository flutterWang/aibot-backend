package model

type Config struct {
	Host          string `toml:"host"`
	Port          string `toml:"port"`
	TcpPort       string `toml:"tcp_port"`
	AiServer      string `toml:"ai_server"`
	JWTSecret     string `toml:"jwt_secret"`
	JWTExpiration int64  `toml:"jwt_expiration"`
	AppID         string `toml:"app_id"`
	AppSecret     string `toml:"app_secret"`

	Mongo *configMongo `toml:"mongo"`
}

type configMongo struct {
	URI      string `toml:"uri"`
	Database string `toml:"database"`
}
