package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadConfigFromEnv(t *testing.T) {
	os.Clearenv()
	accessLog := "/var/log/transcoding-api-access.log"
	setEnvs(map[string]string{
		"ENV":                                      "some_env",
		"SENTRY_DSN":                               "some_dsn",
		"SENTINEL_ADDRS":                           "10.10.10.10:26379,10.10.10.11:26379,10.10.10.12:26379",
		"SENTINEL_MASTER_NAME":                     "super-master",
		"REDIS_ADDR":                               "localhost:6379",
		"REDIS_PASSWORD":                           "super-secret",
		"REDIS_POOL_SIZE":                          "100",
		"REDIS_POOL_TIMEOUT_SECONDS":               "10",
		"ENCODINGCOM_USER_ID":                      "myuser",
		"ENCODINGCOM_USER_KEY":                     "secret-key",
		"ENCODINGCOM_DESTINATION":                  "https://safe-stuff",
		"ENCODINGCOM_STATUS_ENDPOINT":              "https://safe-status",
		"ENCODINGCOM_REGION":                       "sa-east-1",
		"AWS_ACCESS_KEY_ID":                        "AKIANOTREALLY",
		"AWS_SECRET_ACCESS_KEY":                    "secret-key",
		"AWS_REGION":                               "us-east-1",
		"ELASTICTRANSCODER_PIPELINE_ID":            "mypipeline",
		"ELEMENTALCONDUCTOR_HOST":                  "elemental-server",
		"ELEMENTALCONDUCTOR_USER_LOGIN":            "myuser",
		"ELEMENTALCONDUCTOR_API_KEY":               "secret-key",
		"ELEMENTALCONDUCTOR_AUTH_EXPIRES":          "30",
		"ELEMENTALCONDUCTOR_AWS_ACCESS_KEY_ID":     "AKIANOTREALLY",
		"ELEMENTALCONDUCTOR_AWS_SECRET_ACCESS_KEY": "secret-key",
		"ELEMENTALCONDUCTOR_DESTINATION":           "https://safe-stuff",
		"BITMOVIN_API_KEY":                         "secret-key",
		"BITMOVIN_ENDPOINT":                        "bitmovin",
		"BITMOVIN_TIMEOUT":                         "3",
		"BITMOVIN_AWS_ACCESS_KEY_ID":               "AKIANOTREALLY",
		"BITMOVIN_AWS_SECRET_ACCESS_KEY":           "secret-key",
		"BITMOVIN_DESTINATION":                     "https://safe-stuff",
		"BITMOVIN_AWS_STORAGE_REGION":              "US_WEST_1",
		"BITMOVIN_ENCODING_REGION":                 "GOOGLE_EUROPE_WEST_1",
		"BITMOVIN_ENCODING_VERSION":                "notstable",
		"MEDIACONVERT_AWS_ACCESS_KEY_ID":           "mc-aws-access-key-id",
		"MEDIACONVERT_AWS_SECRET_ACCESS_KEY":       "mc-aws-secret-access-key",
		"MEDIACONVERT_AWS_REGION":                  "mc-aws-region",
		"MEDIACONVERT_ENDPOINT":                    "http://mc-endpoint.tld",
		"MEDIACONVERT_QUEUE_ARN":                   "arn:aws:mediaconvert:us-east-1:some-queue:queues/Default",
		"MEDIACONVERT_PREFERRED_QUEUE_ARN":         "arn:aws:mediaconvert:us-east-1:some-queue:queues/Preferred",
		"MEDIACONVERT_ROLE_ARN":                    "arn:aws:iam::some-account:role/some-role",
		"MEDIACONVERT_DESTINATION":                 "s3://mc-destination/",
		"FLOCK_ENDPOINT":                           "https://flock.domain",
		"FLOCK_CREDENTIAL":                         "secret-token",
		"SWAGGER_MANIFEST_PATH":                    "/opt/video-transcoding-api-swagger.json",
		"HTTP_ACCESS_LOG":                          accessLog,
		"HTTP_PORT":                                "8080",
		"DEFAULT_SEGMENT_DURATION":                 "3",
		"LOGGING_LEVEL":                            "debug",
	})
	cfg := LoadConfig()
	expectedCfg := Config{
		DefaultSegmentDuration: 3,
		Env:                    "some_env",
		SentryDSN:              "some_dsn",
		EncodingCom: &EncodingCom{
			UserID:         "myuser",
			UserKey:        "secret-key",
			Destination:    "https://safe-stuff",
			StatusEndpoint: "https://safe-status",
			Region:         "sa-east-1",
		},
		Hybrik: &Hybrik{
			ComplianceDate: "20170601",
			PresetPath:     "transcoding-api-presets",
		},
		Zencoder: &Zencoder{},
		ElasticTranscoder: &ElasticTranscoder{
			AccessKeyID:     "AKIANOTREALLY",
			SecretAccessKey: "secret-key",
			Region:          "us-east-1",
			PipelineID:      "mypipeline",
		},
		ElementalConductor: &ElementalConductor{
			Host:            "elemental-server",
			UserLogin:       "myuser",
			APIKey:          "secret-key",
			AuthExpires:     30,
			AccessKeyID:     "AKIANOTREALLY",
			SecretAccessKey: "secret-key",
			Destination:     "https://safe-stuff",
		},
		Bitmovin: &Bitmovin{
			APIKey:           "secret-key",
			Endpoint:         "bitmovin",
			Timeout:          3,
			AccessKeyID:      "AKIANOTREALLY",
			SecretAccessKey:  "secret-key",
			AWSStorageRegion: "US_WEST_1",
			Destination:      "https://safe-stuff",
			EncodingRegion:   "GOOGLE_EUROPE_WEST_1",
			EncodingVersion:  "notstable",
		},
		MediaConvert: &MediaConvert{
			AccessKeyID:       "mc-aws-access-key-id",
			SecretAccessKey:   "mc-aws-secret-access-key",
			Region:            "mc-aws-region",
			Endpoint:          "http://mc-endpoint.tld",
			DefaultQueueARN:   "arn:aws:mediaconvert:us-east-1:some-queue:queues/Default",
			PreferredQueueARN: "arn:aws:mediaconvert:us-east-1:some-queue:queues/Preferred",
			Role:              "arn:aws:iam::some-account:role/some-role",
			Destination:       "s3://mc-destination/",
		},
		Flock: &Flock{
			Endpoint:   "https://flock.domain",
			Credential: "secret-token",
		},
	}
	diff := cmp.Diff(*cfg, expectedCfg)
	if diff != "" {
		t.Errorf("LoadConfig(): wrong config\nWant %#v\nGot %#v\nDiff: %v", expectedCfg, *cfg, diff)
	}
}

func TestLoadConfigFromEnvWithDefaults(t *testing.T) {
	os.Clearenv()
	accessLog := "/var/log/transcoding-api-access.log"
	setEnvs(map[string]string{
		"SENTINEL_ADDRS":                           "10.10.10.10:26379,10.10.10.11:26379,10.10.10.12:26379",
		"SENTINEL_MASTER_NAME":                     "super-master",
		"REDIS_PASSWORD":                           "super-secret",
		"REDIS_POOL_SIZE":                          "100",
		"REDIS_POOL_TIMEOUT_SECONDS":               "10",
		"REDIS_IDLE_TIMEOUT_SECONDS":               "30",
		"REDIS_IDLE_CHECK_FREQUENCY_SECONDS":       "20",
		"ENCODINGCOM_USER_ID":                      "myuser",
		"ENCODINGCOM_USER_KEY":                     "secret-key",
		"ENCODINGCOM_DESTINATION":                  "https://safe-stuff",
		"AWS_ACCESS_KEY_ID":                        "AKIANOTREALLY",
		"AWS_SECRET_ACCESS_KEY":                    "secret-key",
		"AWS_REGION":                               "us-east-1",
		"ELASTICTRANSCODER_PIPELINE_ID":            "mypipeline",
		"ELEMENTALCONDUCTOR_HOST":                  "elemental-server",
		"ELEMENTALCONDUCTOR_USER_LOGIN":            "myuser",
		"ELEMENTALCONDUCTOR_API_KEY":               "secret-key",
		"ELEMENTALCONDUCTOR_AUTH_EXPIRES":          "30",
		"ELEMENTALCONDUCTOR_AWS_ACCESS_KEY_ID":     "AKIANOTREALLY",
		"ELEMENTALCONDUCTOR_AWS_SECRET_ACCESS_KEY": "secret-key",
		"ELEMENTALCONDUCTOR_DESTINATION":           "https://safe-stuff",
		"BITMOVIN_API_KEY":                         "secret-key",
		"BITMOVIN_AWS_ACCESS_KEY_ID":               "AKIANOTREALLY",
		"BITMOVIN_AWS_SECRET_ACCESS_KEY":           "secret-key",
		"BITMOVIN_DESTINATION":                     "https://safe-stuff",
		"SWAGGER_MANIFEST_PATH":                    "/opt/video-transcoding-api-swagger.json",
		"HTTP_ACCESS_LOG":                          accessLog,
		"HTTP_PORT":                                "8080",
	})
	cfg := LoadConfig()
	expectedCfg := Config{
		Env:                    "dev",
		DefaultSegmentDuration: 5,
		EncodingCom: &EncodingCom{
			UserID:         "myuser",
			UserKey:        "secret-key",
			Destination:    "https://safe-stuff",
			StatusEndpoint: "http://status.encoding.com",
		},
		ElasticTranscoder: &ElasticTranscoder{
			AccessKeyID:     "AKIANOTREALLY",
			SecretAccessKey: "secret-key",
			Region:          "us-east-1",
			PipelineID:      "mypipeline",
		},
		ElementalConductor: &ElementalConductor{
			Host:            "elemental-server",
			UserLogin:       "myuser",
			APIKey:          "secret-key",
			AuthExpires:     30,
			AccessKeyID:     "AKIANOTREALLY",
			SecretAccessKey: "secret-key",
			Destination:     "https://safe-stuff",
		},
		Hybrik: &Hybrik{
			ComplianceDate: "20170601",
			PresetPath:     "transcoding-api-presets",
		},
		Zencoder: &Zencoder{},
		Bitmovin: &Bitmovin{
			APIKey:           "secret-key",
			Endpoint:         "https://api.bitmovin.com/v1/",
			Timeout:          5,
			AccessKeyID:      "AKIANOTREALLY",
			SecretAccessKey:  "secret-key",
			Destination:      "https://safe-stuff",
			AWSStorageRegion: "US_EAST_1",
			EncodingRegion:   "AWS_US_EAST_1",
			EncodingVersion:  "STABLE",
		},
		MediaConvert: &MediaConvert{},
		Flock:        &Flock{},
	}
	diff := cmp.Diff(*cfg, expectedCfg)
	if diff != "" {
		t.Errorf("LoadConfig(): wrong config\nWant %#v\nGot %#v\nDiff: %v", expectedCfg, *cfg, diff)
	}
}

func setEnvs(envs map[string]string) {
	for k, v := range envs {
		os.Setenv(k, v)
	}
}
