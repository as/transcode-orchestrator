package job

import (
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/cbsinteractive/pkg/timecode"
	"github.com/cbsinteractive/pkg/video"
	"github.com/gofrs/uuid"
)

// Job is a transcoding job
type Job struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt time.Time
	Labels    []string

	Provider      string `json:"provider"`
	ProviderJobID string

	Input  File
	Output Dir

	Streaming Streaming

	ExecutionFeatures  ExecutionFeatures
	ExecutionEnv       ExecutionEnvironment
	ExecutionCfgReport string

	SidecarAssets map[SidecarAssetKind]string
}

func (j *Job) Asset(sidecar string) *File {
	loc := j.SidecarAssets[sidecar]
	if loc == "" {
		return nil
	}
	return &File{Name: loc}
}

// State is the state of a transcoding job.
type State string

const (
	StateUnknown  = State("unknown")
	StateQueued   = State("queued")
	StateStarted  = State("started")
	StateFinished = State("finished")
	StateFailed   = State("failed")
	StateCanceled = State("canceled")
)

type Provider struct {
	Name   string                 `json:"name,omitempty"`
	JobID  string                 `json:"job_id,omitempty"`
	Status map[string]interface{} `json:"status,omitempty"`
}

// Status is the representation of the status
type Status struct {
	ID     string   `json:"jobID,omitempty"`
	Labels []string `json:"labels,omitempty"`

	State    State   `json:"status,omitempty"`
	Msg      string  `json:"msg,omitempty"`
	Progress float64 `json:"progress"`

	Input  File `json:"input"`
	Output Dir  `json:"output"`

	ProviderName   string                 `json:"providerName,omitempty"`
	ProviderJobID  string                 `json:"providerJobId,omitempty"`
	ProviderStatus map[string]interface{} `json:"providerStatus,omitempty"`
}

// Dir is a named directory of files
type Dir struct {
	Path string `json:"path,omitempty"`
	File []File `json:"files,omitempty"`
}

func (d *Dir) Add(f ...File) {
	d.File = append(d.File, f...)
}

func (d Dir) Location() url.URL {
	u, _ := url.Parse(d.Path)
	if u == nil {
		return url.URL{}
	}
	return *u
}

func (j Job) Location(file string) string {
	u := j.Output.Location()
	u.Path = path.Join(u.Path, j.rootFolder(), file)
	return u.String()
}

func (j Job) rootFolder() string {
	if j.Name != "" {
		if _, err := uuid.FromString(j.Name); err == nil {
			return j.Name
		}
	}
	return j.ID
}

type SidecarAssetKind = string

const SidecarAssetKindDolbyVisionMetadata SidecarAssetKind = "dolbyVisionMetadata"

// ExecutionEnvironment contains configurations for the environment used while transcoding
type ExecutionEnvironment struct {
	Cloud       string
	Region      string
	ComputeTags map[ComputeClass]string
	InputAlias  string
	OutputAlias string
}

// ComputeClass represents a group of resources with similar capability
type ComputeClass = string

// ComputeClassTranscodeDefault runs any default transcodes
// ComputeClassDolbyVisionTranscode runs Dolby Vision transcodes
// ComputeClassDolbyVisionPreprocess runs Dolby Vision pre-processing
// ComputeClassDolbyVisionMezzQC runs QC check on the mezzanine
const (
	ComputeClassTranscodeDefault      ComputeClass = "transcodeDefault"
	ComputeClassDolbyVisionTranscode  ComputeClass = "doViTranscode"
	ComputeClassDolbyVisionPreprocess ComputeClass = "doViPreprocess"
	ComputeClassDolbyVisionMezzQC     ComputeClass = "doViMezzQC"
)

// Streaming configures Adaptive Streaming jobs
type Streaming struct {
	SegmentDuration  uint
	Protocol         string
	PlaylistFileName string
}

// ScanType is a string that represents the scan type of the content.
type ScanType string

// ScanTypeProgressive and other supported types
const (
	ScanTypeProgressive ScanType = "progressive"
	ScanTypeInterlaced  ScanType = "interlaced"
	ScanTypeUnknown     ScanType = "unknown"
)

//ChannelLayout describes layout of an audio channel
type ChannelLayout string

const (
	ChannelLayoutCenter        ChannelLayout = "C"
	ChannelLayoutLeft          ChannelLayout = "L"
	ChannelLayoutRight         ChannelLayout = "R"
	ChannelLayoutLeftSurround  ChannelLayout = "Ls"
	ChannelLayoutRightSurround ChannelLayout = "Rs"
	ChannelLayoutLeftBack      ChannelLayout = "Lb"
	ChannelLayoutRightBack     ChannelLayout = "Rb"
	ChannelLayoutLeftTotal     ChannelLayout = "Lt"
	ChannelLayoutRightTotal    ChannelLayout = "Rt"
	ChannelLayoutLFE           ChannelLayout = "LFE"
)

// AudioChannel describes the position and attributes of a
// single channel of audio inside a container
type AudioChannel struct {
	TrackIdx, ChannelIdx int
	Layout               string
}

//AudioDownmix holds source and output channels for providers
//to handle downmixing
type Downmix struct {
	Src []AudioChannel
	Dst []AudioChannel
}

// ExecutionFeatures is a map whose key is a custom feature name and value is a json string
// representing the corresponding custom feature definition
type ExecutionFeatures map[string]interface{}

// File
type File struct {
	// Description     string `json:"description,omitempty"`
	// SourceContainer string `json:"sourceContainer,omitempty"`
	// TwoPass         bool   `json:"twoPass"`

	Size      int64         `json:"size,omitempty"`
	Duration  time.Duration `json:"dur,omitempty"`
	Name      string        `json:"name,omitempty"`
	Container string        `json:"container,omitempty"`
	Video     Video         `json:"video,omitempty"`
	Audio     Audio         `json:"audio,omitempty"`

	Splice                  timecode.Splice `json:"splice,omitempty"`
	Downmix                 *Downmix
	ExplicitKeyframeOffsets []float64
}

func (f File) URL() url.URL {
	u, _ := url.Parse(f.Name)
	if u == nil {
		return url.URL{}
	}
	return *u
}
func (f File) Provider() string {
	return f.URL().Scheme
}
func (f File) Type() string {
	return strings.TrimPrefix(path.Ext(f.URL().Path), ".")
}

// Video transcoding parameters
type Video struct {
	Codec   string `json:"codec,omitempty"`
	Profile string `json:"profile,omitempty"`
	Level   string `json:"level,omitempty"`

	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	Scantype string `json:"scantype,omitempty"`

	FPS     float64 `json:"fps,omitempty"`
	Bitrate Bitrate `json:"bitrate"`
	Gop     Gop     `json:"gop"`

	HDR10       HDR10       `json:"hdr10"`
	DolbyVision DolbyVision `json:"dolbyVision"`
	Overlays    *Overlays   `json:"overlays,omitempty"`
	Crop        video.Crop  `json:"crop"`
}

type Bitrate struct {
	BPS     int    `json:"bps"`
	Control string `json:"control"`
	TwoPass bool   `json:"twopass"`
}

func (b Bitrate) Kbps() int {
	return b.BPS / 1000
}

type Gop struct {
	Unit string  `json:"unit,omitempty"`
	Size float64 `json:"size,omitempty"`
	Mode string  `json:"mode,omitempty"`
}

func (g Gop) Seconds() bool {
	return g.Unit == "seconds"
}

// GopUnit defines the unit used to measure gops
type GopUnit = string

const (
	GopUnitFrames  GopUnit = "frames"
	GopUnitSeconds GopUnit = "seconds"
)

//Overlays defines all the overlay settings for a Video preset
type Overlays struct {
	Images         []Image   `json:"images,omitempty"`
	TimecodeBurnin *Timecode `json:"timecodeBurnin,omitempty"`
}

//Image defines the image overlay settings
type Image struct {
	URL string `json:"url"`
}

// Timecode settings
type Timecode struct {
	FontSize int    `json:"fontSize,omitempty"`
	Position int    `json:"position,omitempty"`
	Prefix   string `json:"prefix,omitempty"`
}

// HDR10 configurations and metadata
type HDR10 struct {
	Enabled       bool   `json:"enabled"`
	MaxCLL        int    `json:"maxCLL,omitempty"`
	MaxFALL       int    `json:"maxFALL,omitempty"`
	MasterDisplay string `json:"masterDisplay,omitempty"`
}

// DolbyVision settings
type DolbyVision struct {
	Enabled bool `json:"enabled"`
}

// Audio defines audio transcoding parameters
type Audio struct {
	Codec     string `json:"codec,omitempty"`
	Bitrate   int    `json:"bitrate,omitempty"`
	Normalize bool   `json:"normalize,omitempty"`
	Discrete  bool   `json:"discrete,omitempty"`
}
