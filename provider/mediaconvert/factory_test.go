package mediaconvert

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	mc "github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/cbsinteractive/transcode-orchestrator/config"
	"github.com/google/go-cmp/cmp"
)

var cfgWithoutCredsAndRegion = config.Config{
	MediaConvert: &config.MediaConvert{
		Endpoint:        "http://some/endpoint",
		DefaultQueueARN: "arn:some:queue",
		Role:            "arn:some:role",
	},
}

var cfgWithCredsAndRegion = config.Config{
	MediaConvert: &config.MediaConvert{
		AccessKeyID:       "cfg_access_key_id",
		SecretAccessKey:   "cfg_secret_access_key",
		Endpoint:          "http://some/endpoint",
		DefaultQueueARN:   "arn:some:queue",
		PreferredQueueARN: "arn:some:preferred:queue",
		Role:              "arn:some:role",
		Region:            "us-cfg-region-1",
	},
}

func TestFactory(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		cfg        config.Config
		wantCreds  aws.Credentials
		wantRegion string
		wantErrMsg string
	}{
		{
			name: "AWSCredentialsConfig",
			envVars: map[string]string{
				"AWS_ACCESS_KEY_ID":     "env_access_key_id",
				"AWS_SECRET_ACCESS_KEY": "env_secret_access_key",
				"AWS_DEFAULT_REGION":    "us-north-1",
			},
			cfg: cfgWithCredsAndRegion,
			wantCreds: aws.Credentials{
				AccessKeyID:     "cfg_access_key_id",
				SecretAccessKey: "cfg_secret_access_key",
				Source:          aws.StaticCredentialsProviderName,
			},
			wantRegion: "us-cfg-region-1",
		},
		{
			// NODE(as); we should get rid of this in favor of explicit configuration
			name: "AWSCredentialsEnvironment",
			envVars: map[string]string{
				"AWS_ACCESS_KEY_ID":     "env_access_key_id",
				"AWS_SECRET_ACCESS_KEY": "env_secret_access_key",
				"AWS_DEFAULT_REGION":    "us-north-1",
				"AWS_PROFILE":           "",
			},
			cfg: cfgWithoutCredsAndRegion,
			wantCreds: aws.Credentials{
				AccessKeyID:     "env_access_key_id",
				SecretAccessKey: "env_secret_access_key",
				Source:          external.CredentialsSourceName,
			},
			wantRegion: "us-north-1",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				resetFunc, err := setenvReset(k, v)
				if err != nil {
					t.Errorf("running os env reset: %v", err)
				}
				defer resetFunc()
			}

			provider, err := mediaconvertFactory(&tt.cfg)
			if err != nil {
				if tt.wantErrMsg != err.Error() {
					t.Fatalf("mcProvider.CreatePreset() error = %v, wantErr %q", err, tt.wantErrMsg)
				}
			}

			p, ok := provider.(*driver)
			if !ok {
				t.Fatalf("factory didn't return a mediaconvert provider")
			}

			client, ok := p.client.(*mc.Client)
			if !ok {
				t.Fatalf("factory returned a mediaconvert provider with a non-aws client implementation")
			}

			creds, err := client.Credentials.Retrieve(context.Background())
			if err != nil {
				t.Fatalf("error retrieving aws credentials: %v", err)
			}

			if g, e := creds, tt.wantCreds; !reflect.DeepEqual(g, e) {
				t.Fatalf("unexpected credentials\nWant %+v\nGot %+v\nDiff %s",
					e, g, cmp.Diff(e, g))
			}

			if g, e := client.Config.Region, tt.wantRegion; g != e {
				t.Fatalf("expected region %q, got %q", e, g)
			}
		})
	}
}

func setenvReset(name, val string) (resetEnv func(), rerr error) {
	cached := os.Getenv(name)
	err := os.Setenv(name, val)
	if err != nil {
		return nil, err
	}
	return func() {
		os.Setenv(name, cached)
	}, nil
}
