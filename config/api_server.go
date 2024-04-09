package config

type ApiServerConfig struct {
	Host                   string `mapstructure:"APP_HOST"`
	Port                   int    `mapstructure:"APP_PORT"`
	PathPrefix             string `mapstructure:"APP_PATH_PREFIX"`
	RequestIdHeader        string `mapstructure:"APP_REQUEST_ID_HEADER"`
	TrustProxy             string `mapstructure:"APP_TRUST_PROXY"`
	CorsEnabled            bool   `mapstructure:"APP_CORS_ENABLED"`
	CorsAllowedOrigin      string `mapstructure:"APP_CORS_ALLOWED_ORIGIN"`
	CorsAllowedCredentials bool   `mapstructure:"APP_CORS_ALLOWED_CREDENTIALS"`
	CorsMaxAge             int    `mapstructure:"APP_CORS_MAX_AGE"`
	ReadOnly               bool   `mapstructure:"APP_READ_ONLY"`
	LogRequest             bool   `mapstructure:"LOG_REQUEST"`
	JwtSecretKey           string `mapstructure:"JWT_SECRET_KEY"`
	GoogleClientId         string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	RedisAddress           string `mapstructure:"REDIS_ADDRESS"`
	EmailSenderName        string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress     string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword    string `mapstructure:"EMAIL_SENDER_PASSWORD"`
	AwsRegion              string `mapstructure:"AWS_REGION"`
	AwsBucketName          string `mapstructure:"AWS_BUCKET_NAME"`
	AwsAccessKeyID         string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey     string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
}
