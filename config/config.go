package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/zsiec/pkg/tracing"
)

// Config is a struct to contain all the needed configuration for the
// Transcoding API.
type Config struct {
	DefaultSegmentDuration uint   `envconfig:"DEFAULT_SEGMENT_DURATION" default:"5"`
	SentryDSN              string `envconfig:"SENTRY_DSN"`
	Env                    string `envconfig:"ENV" default:"dev"`
	EnableXray             bool   `envconfig:"ENABLE_XRAY"`
	EnableXrayAWSPlugins   bool   `envconfig:"ENABLE_XRAYAWSPLUGINS"`
	EncodingCom            *EncodingCom
	ElasticTranscoder      *ElasticTranscoder
	ElementalConductor     *ElementalConductor
	Hybrik                 *Hybrik
	Zencoder               *Zencoder
	Bitmovin               *Bitmovin
	MediaConvert           *MediaConvert
	Flock                  *Flock
	Tracer                 tracing.Tracer `ignored:"true"`
}

// EncodingCom represents the set of configurations for the Encoding.com
// provider.
type EncodingCom struct {
	UserID         string `envconfig:"ENCODINGCOM_USER_ID"`
	UserKey        string `envconfig:"ENCODINGCOM_USER_KEY"`
	Destination    string `envconfig:"ENCODINGCOM_DESTINATION"`
	Region         string `envconfig:"ENCODINGCOM_REGION"`
	StatusEndpoint string `envconfig:"ENCODINGCOM_STATUS_ENDPOINT" default:"http://status.encoding.com"`
}

// Zencoder represents the set of configurations for the Zencoder
// provider.
type Zencoder struct {
	APIKey      string `envconfig:"ZENCODER_API_KEY"`
	Destination string `envconfig:"ZENCODER_DESTINATION"`
}

// ElasticTranscoder represents the set of configurations for the Elastic
// Transcoder provider.
type ElasticTranscoder struct {
	AccessKeyID     string `envconfig:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY"`
	Region          string `envconfig:"AWS_REGION"`
	PipelineID      string `envconfig:"ELASTICTRANSCODER_PIPELINE_ID"`
}

// ElementalConductor represents the set of configurations for the Elemental
// Conductor provider.
type ElementalConductor struct {
	Host            string `envconfig:"ELEMENTALCONDUCTOR_HOST"`
	UserLogin       string `envconfig:"ELEMENTALCONDUCTOR_USER_LOGIN"`
	APIKey          string `envconfig:"ELEMENTALCONDUCTOR_API_KEY"`
	AuthExpires     int    `envconfig:"ELEMENTALCONDUCTOR_AUTH_EXPIRES"`
	AccessKeyID     string `envconfig:"ELEMENTALCONDUCTOR_AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"ELEMENTALCONDUCTOR_AWS_SECRET_ACCESS_KEY"`
	Destination     string `envconfig:"ELEMENTALCONDUCTOR_DESTINATION"`
}

// Bitmovin represents the set of configurations for the Bitmovin
// provider.
type Bitmovin struct {
	APIKey             string `envconfig:"BITMOVIN_API_KEY"`
	Endpoint           string `envconfig:"BITMOVIN_ENDPOINT" default:"https://api.bitmovin.com/v1/"`
	Timeout            uint   `envconfig:"BITMOVIN_TIMEOUT" default:"5"`
	AccessKeyID        string `envconfig:"BITMOVIN_AWS_ACCESS_KEY_ID"`
	SecretAccessKey    string `envconfig:"BITMOVIN_AWS_SECRET_ACCESS_KEY"`
	AWSStorageRegion   string `envconfig:"BITMOVIN_AWS_STORAGE_REGION" default:"US_EAST_1"`
	GCSAccessKeyID     string `envconfig:"BITMOVIN_GCS_ACCESS_KEY_ID"`
	GCSSecretAccessKey string `envconfig:"BITMOVIN_GCS_SECRET_ACCESS_KEY"`
	GCSStorageRegion   string `envconfig:"BITMOVIN_GCS_STORAGE_REGION"`
	Destination        string `envconfig:"BITMOVIN_DESTINATION"`
	EncodingRegion     string `envconfig:"BITMOVIN_ENCODING_REGION" default:"AWS_US_EAST_1"`
	EncodingVersion    string `envconfig:"BITMOVIN_ENCODING_VERSION" default:"STABLE"`
}

// Hybrik represents the set of configurations for the Hybrik
// provider.
type Hybrik struct {
	URL               string `envconfig:"HYBRIK_URL"`
	ComplianceDate    string `envconfig:"HYBRIK_COMPLIANCE_DATE" default:"20170601"`
	OAPIKey           string `envconfig:"HYBRIK_OAPI_KEY"`
	OAPISecret        string `envconfig:"HYBRIK_OAPI_SECRET"`
	AuthKey           string `envconfig:"HYBRIK_AUTH_KEY"`
	AuthSecret        string `envconfig:"HYBRIK_AUTH_SECRET"`
	Destination       string `envconfig:"HYBRIK_DESTINATION"`
	GCPCredentialsKey string `envconfig:"HYBRIK_GCP_CREDENTIALS_KEY"`
	PresetPath        string `envconfig:"HYBRIK_PRESET_PATH" default:"transcoding-api-presets"`
}

// MediaConvert represents the set of configurations for the MediaConvert
// provider.
type MediaConvert struct {
	AccessKeyID       string `envconfig:"MEDIACONVERT_AWS_ACCESS_KEY_ID"`
	SecretAccessKey   string `envconfig:"MEDIACONVERT_AWS_SECRET_ACCESS_KEY"`
	Region            string `envconfig:"MEDIACONVERT_AWS_REGION"`
	Endpoint          string `envconfig:"MEDIACONVERT_ENDPOINT"`
	DefaultQueueARN   string `envconfig:"MEDIACONVERT_QUEUE_ARN"`
	PreferredQueueARN string `envconfig:"MEDIACONVERT_PREFERRED_QUEUE_ARN"`
	Role              string `envconfig:"MEDIACONVERT_ROLE_ARN"`
	Destination       string `envconfig:"MEDIACONVERT_DESTINATION"`
}

// Flock represents the set of configurations for the Flock
// provider.
type Flock struct {
	Endpoint   string `envconfig:"FLOCK_ENDPOINT"`
	Credential string `envconfig:"FLOCK_CREDENTIAL"`
}

// LoadConfig loads the configuration of the API using environment variables.
func LoadConfig() *Config {
	var cfg Config
	envconfig.Process("", &cfg)
	return &cfg
}
