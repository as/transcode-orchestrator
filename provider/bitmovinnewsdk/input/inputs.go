package input

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/bitmovin/bitmovin-api-sdk-go/model"

	"github.com/bitmovin/bitmovin-api-sdk-go"
	"github.com/NYTimes/video-transcoding-api/config"
	"github.com/pkg/errors"
)

var s3Pattern = regexp.MustCompile(`^s3://`)
var httpPattern = regexp.MustCompile(`^http://`)
var httpsPattern = regexp.MustCompile(`^https://`)

// New creates an input and returns an inputID and the media path or an error
func New(srcMediaLoc string, api *bitmovin.BitmovinApi, cfg *config.Bitmovin) (inputID string, path string, err error) {
	if s3Pattern.MatchString(srcMediaLoc) {
		return s3(srcMediaLoc, api, cfg)
	}
	if httpPattern.MatchString(srcMediaLoc) {
		return http(srcMediaLoc, api, cfg)
	}
	if httpsPattern.MatchString(srcMediaLoc) {
		return https(srcMediaLoc, api, cfg)
	}

	return "", "", errors.New("only s3, http, and https urls are supported")
}

func s3(srcMediaLoc string, api *bitmovin.BitmovinApi, cfg *config.Bitmovin) (inputID string, path string, err error) {
	bucket, mediaPath, err := parseS3URL(srcMediaLoc)
	if err != nil {
		return "", "", err
	}

	input, err := api.Encoding.Inputs.S3.Create(model.S3Input{
		CloudRegion: model.AwsCloudRegion(cfg.AWSStorageRegion),
		BucketName:  bucket,
		AccessKey:   cfg.AccessKeyID,
		SecretKey:   cfg.SecretAccessKey,
	})
	if err != nil {
		return "", "", errors.Wrap(err, "creating s3 input")
	}

	return input.Id, mediaPath, nil
}

func parseS3URL(s3URL string) (bucketName string, objectKey string, err error) {
	u, err := url.Parse(s3URL)
	if err != nil || u.Scheme != "s3" {
		return "", "", errors.Wrap(err, "parsing s3 url")
	}
	return u.Host, strings.TrimLeft(u.Path, "/"), nil
}

func http(srcMediaLoc string, api *bitmovin.BitmovinApi, cfg *config.Bitmovin) (inputID string, path string, err error) {
	u, err := url.Parse(srcMediaLoc)
	if err != nil {
		return "", "", errors.Wrap(err, "parsing src media url")
	}

	input, err := api.Encoding.Inputs.Http.Create(model.HttpInput{
		Host: u.Host,
	})
	if err != nil {
		return "", "", errors.Wrap(err, "creating http input")
	}

	return input.Id, u.Path, nil
}

func https(srcMediaLoc string, api *bitmovin.BitmovinApi, cfg *config.Bitmovin) (inputID string, path string, err error) {
	u, err := url.Parse(srcMediaLoc)
	if err != nil {
		return "", "", errors.Wrap(err, "parsing src media url")
	}

	input, err := api.Encoding.Inputs.Https.Create(model.HttpsInput{
		Host: u.Host,
	})
	if err != nil {
		return "", "", errors.Wrap(err, "creating https input")
	}

	return input.Id, u.Path, nil
}