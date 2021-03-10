package hybrik

import (
	"path"
	"regexp"
	"strings"

	hy "github.com/cbsinteractive/hybrik-sdk-go"
	"github.com/cbsinteractive/transcode-orchestrator/job"
)

type taskWithOutputMatcher struct {
	kind     string
	uidRegex *regexp.Regexp
}

var match = regexp.MustCompile

var outputMatchers = []struct {
	kind     string
	uidRegex *regexp.Regexp
}{
	{"Dolby Vision", match(`post_transcode_stage_[\d]+$`)},
	{"Dolby Vision", match(`dolby_vision_[\d]+$`)},
	{"Transcode", match(`transcode_task_[\d]+$`)},
	{"Package", match(`packager$`)},
	{"Combine Segments", match(`combiner_[\d]+$`)},
}

func hasOutputs(task hy.TaskResult) bool {
	for _, matcher := range outputMatchers {
		if matcher.kind != task.Kind {
			continue
		}
		if matcher.uidRegex.Match([]byte(task.UID)) {
			return true
		}
	}

	return false
}

func filesFrom(task hy.TaskResult) (files []job.File, ok bool, err error) {
	// ensure the task type results in outputs
	if !hasOutputs(task) {
		return nil, false, nil
	}

	for _, d := range task.Documents {
		for _, a := range d.ResultPayload.Payload.AssetVersions {
			dir := job.File{Name: a.Location.Path}
			for _, c := range a.AssetComponents {
				files = append(files, job.File{
					Name:      dir.Join(c.Name).Name,
					Container: containerFrom(c),
					Size:      int64(c.Descriptor.Size),
				})
			}
		}
	}

	return files, len(files) > 0, nil
}

const assetMediaInfoType = "ASSET"

func containerFrom(component hy.AssetComponentResult) string {
	if infos := component.MediaAnalyze.MediaInfo; len(infos) > 0 {
		for _, i := range infos {
			if i.StreamType == assetMediaInfoType && i.ASSET.Format != "" {
				return i.ASSET.Format
			}
		}
	}

	return strings.Replace(path.Ext(component.Name), ".", "", -1)
}
