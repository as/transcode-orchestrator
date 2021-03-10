package mediaconvert

import (
	"context"
	"fmt"
	"path"
	"strings"

	mc "github.com/aws/aws-sdk-go-v2/service/mediaconvert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/cbsinteractive/pkg/timecode"
	"github.com/cbsinteractive/transcode-orchestrator/config"
	"github.com/cbsinteractive/transcode-orchestrator/job"
	"github.com/cbsinteractive/transcode-orchestrator/provider"
	"github.com/pkg/errors"
)

const (
	Name = "mediaconvert"

	defaultAudioSampleRate     = 48000
	defaultQueueHopTimeoutMins = 1
)

type (
	Status = job.Status
	Job    = job.Job
)

func init() {
	err := provider.Register(Name, mediaconvertFactory)
	if err != nil {
		fmt.Printf("registering mediaconvert factory: %v", err)
	}
}

type mediaconvertClient interface {
	CreateJobRequest(*mc.CreateJobInput) mc.CreateJobRequest
	GetJobRequest(*mc.GetJobInput) mc.GetJobRequest
	ListJobsRequest(*mc.ListJobsInput) mc.ListJobsRequest
	CancelJobRequest(*mc.CancelJobInput) mc.CancelJobRequest
	CreatePresetRequest(*mc.CreatePresetInput) mc.CreatePresetRequest
	GetPresetRequest(*mc.GetPresetInput) mc.GetPresetRequest
	DeletePresetRequest(*mc.DeletePresetInput) mc.DeletePresetRequest
}

type driver struct {
	client mediaconvertClient
	cfg    config.MediaConvert
}

type outputCfg struct {
	output   mc.Output
	filename string
}

func splice2clippings(s timecode.Splice, fps float64) (ic []mc.InputClipping) {
	// NOTE(as): While this could be a helper function in the time/timecode package
	// we probably don't want the uglyness of importing the AWS API in that package
	// and having to recognize mc.InputClippings

	// NOTE(as): We need to take into account embedded timecodes. Maybe it would
	// be better to have this be a method on a timecode object or have it passed in as
	// a reference argument (object could also provide fps info)

	for _, r := range s {
		s, e := r.Timecodes(fps)
		ic = append(ic, mc.InputClipping{
			StartTimecode: &s,
			EndTimecode:   &e,
		})
	}
	return ic
}

func (p *driver) createRequest(ctx context.Context, j *Job) (*mc.CreateJobInput, error) {
	outputGroups, err := p.outputGroupsFrom(j)
	if err != nil {
		return nil, fmt.Errorf("mediaconvert: output group generator: %w", err)
	}

	queue := aws.String(p.cfg.DefaultQueueARN)

	var hopDestinations []mc.HopDestination
	if preferred := p.cfg.PreferredQueueARN; p.canUsePreferredQueue(j.Input) && preferred != "" {
		queue = aws.String(preferred)
		hopDestinations = append(hopDestinations, mc.HopDestination{
			WaitMinutes: aws.Int64(defaultQueueHopTimeoutMins),
		})
	}

	var acceleration *mc.AccelerationSettings
	if p.requiresAcceleration(j.Input) {
		acceleration = &mc.AccelerationSettings{Mode: mc.AccelerationModePreferred}
	}

	audio := mc.AudioSelector{
		DefaultSelection: mc.AudioDefaultSelectionDefault,
	}
	if err = audioSelectorFrom(j.Input.Downmix, &audio); err != nil {
		return nil, fmt.Errorf("mediaconvert: audio selectors generator: %w", err)
	}

	return &mc.CreateJobInput{
		AccelerationSettings: acceleration,
		Queue:                queue,
		HopDestinations:      hopDestinations,
		Role:                 aws.String(p.cfg.Role),
		Settings: &mc.JobSettings{
			Inputs: []mc.Input{
				{
					InputClippings: splice2clippings(j.Input.Splice, 0), // TODO(as): Find FPS in job
					FileInput:      aws.String(j.Input.Name),
					AudioSelectors: map[string]mc.AudioSelector{
						"Audio Selector 1": audio,
					},
					VideoSelector: &mc.VideoSelector{
						ColorSpace: mc.ColorSpaceFollow,
					},
					TimecodeSource: mc.InputTimecodeSourceZerobased,
				},
			},
			OutputGroups:   outputGroups,
			TimecodeConfig: &mc.TimecodeConfig{Source: mc.TimecodeSourceZerobased},
		},
		Tags: p.tagsFrom(j.Labels),
	}, nil
}

func (p *driver) Create(ctx context.Context, j *Job) (*Status, error) {
	input, err := p.createRequest(ctx, j)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.CreateJobRequest(input).Send(ctx)
	if err != nil {
		return nil, err
	}
	return &Status{
		Provider:      Name,
		ProviderJobID: aws.StringValue(resp.Job.Id),
		State:         job.StateQueued,
	}, nil
}

func (p *driver) Status(ctx context.Context, job *Job) (*Status, error) {
	jobResp, err := p.client.GetJobRequest(&mc.GetJobInput{
		Id: aws.String(job.ProviderJobID),
	}).Send(ctx)
	if err != nil {
		return &Status{}, errors.Wrap(err, "fetching job info with the mediaconvert API")
	}

	return p.status(job, jobResp.Job), nil
}

func (p *driver) Cancel(ctx context.Context, id string) error {
	_, err := p.client.CancelJobRequest(&mc.CancelJobInput{
		Id: aws.String(id),
	}).Send(ctx)

	return err
}

func (p *driver) Healthcheck() error {
	_, err := p.client.ListJobsRequest(nil).Send(context.Background()) // TODO(as): plump context
	if err != nil {
		return errors.Wrap(err, "listing jobs")
	}
	return nil
}

func (p *driver) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		InputFormats:  []string{"h264", "h265", "hdr10"},
		OutputFormats: []string{"mp4", "hls", "hdr10", "cmaf", "mov"},
		Destinations:  []string{"s3"},
	}
}

func (p *driver) outputGroupsFrom(j *Job) ([]mc.OutputGroup, error) {
	cfg := map[mc.ContainerType][]outputCfg{}
	for _, f := range j.Output.File {
		mc, err := outputFrom(f, j.Input)
		if err != nil {
			return nil, fmt.Errorf("output: %q: %w", f.Name, err)
		}

		cs := mc.ContainerSettings
		if cs == nil {
			return nil, fmt.Errorf("no container was found on outout settings %+v", mc)
		}

		cfg[cs.Container] = append(cfg[cs.Container], outputCfg{
			output:   mc,
			filename: f.Name,
		})
	}

	mcOutputGroups := []mc.OutputGroup{}
	for container, outputs := range cfg {
		mcOutputGroup := mc.OutputGroup{}

		mcOutputs := make([]mc.Output, len(outputs))
		for i, o := range outputs {
			rawExtension := path.Ext(o.filename)
			filename := strings.Replace(path.Base(o.filename), rawExtension, "", 1)
			extension := strings.Replace(rawExtension, ".", "", -1)

			mcOutputs[i] = mc.Output{
				NameModifier:      aws.String(filename),
				Extension:         aws.String(extension),
				ContainerSettings: o.output.ContainerSettings,
				AudioDescriptions: o.output.AudioDescriptions,
				VideoDescription:  o.output.VideoDescription,
			}
		}
		mcOutputGroup.Outputs = mcOutputs

		destination := p.destinationPath(*j, "m")

		switch container {
		case mc.ContainerTypeMp4, mc.ContainerTypeMov, mc.ContainerTypeWebm, mc.ContainerTypeMxf:
			mcOutputGroup.OutputGroupSettings = &mc.OutputGroupSettings{
				Type: mc.OutputGroupTypeFileGroupSettings,
				FileGroupSettings: &mc.FileGroupSettings{
					Destination: aws.String(destination),
				},
			}
		default:
			return nil, fmt.Errorf("container: %w: %q", ErrUnsupported, string(container))
		}

		mcOutputGroups = append(mcOutputGroups, mcOutputGroup)
	}

	return mcOutputGroups, nil
}

func (p *driver) destinationPath(j Job, file string) string {
	if j.Output.Path == "" {
		j.Output.Path = p.cfg.Destination
	}
	return j.Location(file)
}

func (p *driver) status(j *Job, mcJob *mc.Job) *Status {
	status := &Status{
		ProviderJobID: j.ProviderJobID,
		Provider:      Name,
		State:         state(mcJob.Status),
		Msg:           message(mcJob),
		Output: job.Dir{
			Path: p.destinationPath(*j, ""),
		},
	}

	if status.State == job.StateFinished {
		status.Progress = 100
	} else if p := mcJob.JobPercentComplete; p != nil {
		status.Progress = float64(*p)
	}

	var files []job.File
	if settings := mcJob.Settings; settings != nil {
		for _, group := range settings.OutputGroups {
			groupDestination, err := outputGroupDestinationFrom(group)
			if err != nil {
				continue
			}
			for _, output := range group.Outputs {
				file := job.File{}

				if modifier := output.NameModifier; modifier != nil {
					if extension, err := fileExtensionFromContainer(output.ContainerSettings); err == nil {
						file.Name = groupDestination + *modifier + extension
					} else {
						continue
					}
				} else {
					continue
				}

				if video := output.VideoDescription; video != nil {
					if h := video.Height; h != nil {
						file.Video.Height = int(*h)
					}
					if w := video.Width; w != nil {
						file.Video.Width = int(*w)
					}
				}

				if container, err := containerIdentifierFrom(output.ContainerSettings); err == nil {
					file.Container = container
				}

				files = append(files, file)
			}
		}
	}

	status.Output.File = files

	return status
}

func outputGroupDestinationFrom(group mc.OutputGroup) (string, error) {
	if group.OutputGroupSettings == nil {
		return "", errors.New("output group contained no settings")
	}

	switch group.OutputGroupSettings.Type {
	case mc.OutputGroupTypeFileGroupSettings:
		fsSettings := group.OutputGroupSettings.FileGroupSettings
		if fsSettings == nil {
			return "", errors.New("file group settings were nil")
		}

		if fsSettings.Destination == nil {
			return "", errors.New("file group destination was nil")
		}

		return *fsSettings.Destination, nil
	default:
		return "", fmt.Errorf("output enumeration not supported for output group %q",
			group.OutputGroupSettings.Type)
	}
}

// NOTE(as): it seems like the code in this repository wants to decouple
// file extensions from file types for some reason, but the implementation
// in some places resists that
func fileExtensionFromContainer(settings *mc.ContainerSettings) (string, error) {
	if settings == nil {
		return "", errors.New("container settings were nil")
	}
	switch settings.Container {
	case mc.ContainerTypeMp4:
		return ".mp4", nil
	case mc.ContainerTypeMov:
		return ".mov", nil
	case mc.ContainerTypeWebm:
		return ".webm", nil
	default:
		return "", fmt.Errorf("could not determine extension from output container %q", settings.Container)
	}
}

// TODO(as): the only difference between fileExtensionFromContainer
// and containerIdentifierFrom is a period... and the error message...
func containerIdentifierFrom(settings *mc.ContainerSettings) (string, error) {
	if settings == nil {
		return "", errors.New("container settings were nil")
	}
	switch settings.Container {
	case mc.ContainerTypeMp4:
		return "mp4", nil
	case mc.ContainerTypeMov:
		return "mov", nil
	case mc.ContainerTypeWebm:
		return "webm", nil
	default:
		return "", fmt.Errorf("could not determine container identifier from output container %q", settings.Container)
	}
}

func message(job *mc.Job) string {
	if job.ErrorMessage != nil {
		return *job.ErrorMessage
	}
	return string(job.CurrentPhase)
}

func (p *driver) tagsFrom(labels []string) map[string]string {
	tags := make(map[string]string)

	for _, label := range labels {
		tags[label] = "true"
	}

	return tags
}

func mediaconvertFactory(cfg *config.Config) (provider.Provider, error) {
	if cfg.MediaConvert.Endpoint == "" || cfg.MediaConvert.DefaultQueueARN == "" || cfg.MediaConvert.Role == "" {
		return nil, errors.New("incomplete MediaConvert config")
	}

	mcCfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "loading default aws config")
	}

	if cfg.MediaConvert.AccessKeyID+cfg.MediaConvert.SecretAccessKey != "" {
		mcCfg.Credentials = &aws.StaticCredentialsProvider{Value: aws.Credentials{
			AccessKeyID:     cfg.MediaConvert.AccessKeyID,
			SecretAccessKey: cfg.MediaConvert.SecretAccessKey,
		}}
	}

	if cfg.MediaConvert.Region != "" {
		mcCfg.Region = cfg.MediaConvert.Region
	}

	mcCfg.EndpointResolver = &aws.ResolveWithEndpoint{
		URL: cfg.MediaConvert.Endpoint,
	}

	return &driver{
		client: mc.New(mcCfg),
		cfg:    *cfg.MediaConvert,
	}, nil
}
