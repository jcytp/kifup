package env

import (
	"log"
	"os"
)

const (
	envProduction  = "production"
	envStaging     = "staging"
	envDevelopment = "development"
)

type Config struct {
	Environment    string
	SecretKey      []byte
	AllowedOrigins []string
	DatabasePath   string
	AwsRegion      string
	S3BucketName   string
	EmailSender    string
	SwaggerEnable  bool
}

var conf *Config

func Initialize() {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY environment variable is required")
	}
	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		log.Fatal("FRONTEND_ORIGIN environment variable is required")
	}

	env := os.Getenv("ENV")
	switch env {
	case envProduction:
		conf = &Config{
			Environment:    envProduction,
			SecretKey:      []byte(secretKey),
			AllowedOrigins: []string{frontendOrigin},
			DatabasePath:   "kifup.db",
			AwsRegion:      "ap-northeast-1",
			S3BucketName:   "s3-kifup",
			EmailSender:    "system.kifup@jcytp.net",
			SwaggerEnable:  false,
		}
	case envStaging:
		conf = &Config{
			Environment:    envStaging,
			SecretKey:      []byte(secretKey),
			AllowedOrigins: []string{frontendOrigin},
			DatabasePath:   "kifup.db",
			AwsRegion:      "ap-northeast-1",
			S3BucketName:   "s3-kifup-stg",
			EmailSender:    "system.kifup@jcytp.net",
			SwaggerEnable:  true,
		}
	case envDevelopment:
		conf = &Config{
			Environment:    envDevelopment,
			SecretKey:      []byte(secretKey),
			AllowedOrigins: []string{"http://localhost", frontendOrigin},
			DatabasePath:   "tmp/kifup.db",
			AwsRegion:      "",
			S3BucketName:   "",
			EmailSender:    "",
			SwaggerEnable:  true,
		}
	default:
		log.Fatalf("invalid environment: %s", env)
	}
}

func IsProductionn() bool {
	return conf.Environment == envProduction
}

func IsStaging() bool {
	return conf.Environment == envStaging
}

func IsDevelopment() bool {
	return conf.Environment == envDevelopment
}

func Environment() string {
	return conf.Environment
}

func SecretKey() []byte {
	return conf.SecretKey
}

func AllowedOrigins() []string {
	return conf.AllowedOrigins
}

func DatabasePath() string {
	return conf.DatabasePath
}

func AwsRegion() string {
	return conf.AwsRegion
}

func S3BucketName() string {
	return conf.S3BucketName
}

func EmailSender() string {
	return conf.EmailSender
}

func SwaggerEnable() bool {
	return conf.SwaggerEnable
}
