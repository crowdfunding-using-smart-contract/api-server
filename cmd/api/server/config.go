package server

type ApiServerConfig struct {
	APP_HOST                     string `mapstructure:"APP_HOST"`
	APP_PORT                     int    `mapstructure:"APP_PORT"`
	APP_PATH_PREFIX              string `mapstructure:"APP_PATH_PREFIX"`
	APP_REQUEST_ID_HEADER        string `mapstructure:"APP_REQUEST_ID_HEADER"`
	APP_TRUST_PROXY              string `mapstructure:"APP_TRUST_PROXY"`
	APP_CORS_ENABLED             bool   `mapstructure:"APP_CORS_ENABLED"`
	APP_CORS_ALLOWED_ORIGIN      string `mapstructure:"APP_CORS_ALLOWED_ORIGIN"`
	APP_CORS_ALLOWED_CREDENTIALS bool   `mapstructure:"APP_CORS_ALLOWED_CREDENTIALS"`
	APP_CORS_MAX_AGE             int    `mapstructure:"APP_CORS_MAX_AGE"`
	APP_READ_ONLY                bool   `mapstructure:"APP_READ_ONLY"`
	LOG_REQUEST                  bool   `mapstructure:"LOG_REQUEST"`
	JWT_SECRET_KEY               string `mapstructure:"JWT_SECRET_KEY"`
	GOOGLE_CLIENT_ID             string `mapstructure:"GOOGLE_CLIENT_ID"`
	GOOGLE_CLIENT_SECRET         string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	REDIS_ADDRESS                string `mapstructure:"REDIS_ADDRESS"`
	EMAIL_SENDER_NAME            string `mapstructure:"EMAIL_SENDER_NAME"`
	EMAIL_SENDER_ADDRESS         string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EMAIL_SENDER_PASSWORD        string `mapstructure:"EMAIL_SENDER_PASSWORD"`
	AwsRegion                    string `mapstructure:"AWS_REGION"`
	AwsBucketName                string `mapstructure:"AWS_BUCKET_NAME"`
	AwsAccessKeyID               string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey           string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
}
