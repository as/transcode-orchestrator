package mediaconvert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/cbsinteractive/video-transcoding-api/db"
	"github.com/pkg/errors"
)

func h264CodecSettingsFrom(preset db.Preset) (*mediaconvert.VideoCodecSettings, error) {
	bitrate, err := strconv.ParseInt(preset.Video.Bitrate, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing video bitrate %q to int64", preset.Video.Bitrate)
	}

	gopSize, err := strconv.ParseFloat(preset.Video.GopSize, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing gop size %q to int64", preset.Video.GopSize)
	}

	gopUnit, err := h264GopUnitFrom(preset.Video.GopUnit)
	if err != nil {
		return nil, err
	}

	rateControl, err := h264RateControlModeFrom(preset.RateControl)
	if err != nil {
		return nil, err
	}

	profile, err := h264CodecProfileFrom(preset.Video.Profile)
	if err != nil {
		return nil, err
	}

	interlaceMode, err := h264InterlaceModeFrom(preset.Video.InterlaceMode)
	if err != nil {
		return nil, err
	}

	tuning := mediaconvert.H264QualityTuningLevelSinglePassHq
	if preset.TwoPass {
		tuning = mediaconvert.H264QualityTuningLevelMultiPassHq
	}

	return &mediaconvert.VideoCodecSettings{
		Codec: mediaconvert.VideoCodecH264,
		H264Settings: &mediaconvert.H264Settings{
			Bitrate:            aws.Int64(bitrate),
			GopSize:            aws.Float64(gopSize),
			GopSizeUnits:       gopUnit,
			RateControlMode:    rateControl,
			CodecProfile:       profile,
			CodecLevel:         mediaconvert.H264CodecLevelAuto,
			InterlaceMode:      interlaceMode,
			QualityTuningLevel: tuning,
		},
	}, nil
}

func h264GopUnitFrom(gopUnit string) (mediaconvert.H264GopSizeUnits, error) {
	gopUnit = strings.ToLower(gopUnit)
	switch gopUnit {
	case "", "frames":
		return mediaconvert.H264GopSizeUnitsFrames, nil
	case "seconds":
		return mediaconvert.H264GopSizeUnitsSeconds, nil
	default:
		return "", fmt.Errorf("gop unit %q is not supported with mediaconvert", gopUnit)
	}
}

func h264RateControlModeFrom(rateControl string) (mediaconvert.H264RateControlMode, error) {
	rateControl = strings.ToLower(rateControl)
	switch rateControl {
	case "vbr":
		return mediaconvert.H264RateControlModeVbr, nil
	case "", "cbr":
		return mediaconvert.H264RateControlModeCbr, nil
	case "qvbr":
		return mediaconvert.H264RateControlModeQvbr, nil
	default:
		return "", fmt.Errorf("rate control mode %q is not supported with mediaconvert", rateControl)
	}
}

func h264CodecProfileFrom(profile string) (mediaconvert.H264CodecProfile, error) {
	profile = strings.ToLower(profile)
	switch profile {
	case "baseline":
		return mediaconvert.H264CodecProfileBaseline, nil
	case "main":
		return mediaconvert.H264CodecProfileMain, nil
	case "", "high":
		return mediaconvert.H264CodecProfileHigh, nil
	default:
		return "", fmt.Errorf("h264 profile %q is not supported with mediaconvert", profile)
	}
}

func h264InterlaceModeFrom(mode string) (mediaconvert.H264InterlaceMode, error) {
	mode = strings.ToLower(mode)
	switch mode {
	case "", "progressive":
		return mediaconvert.H264InterlaceModeProgressive, nil
	default:
		return "", fmt.Errorf("h264 interlace mode %q is not supported with mediaconvert", mode)
	}
}
