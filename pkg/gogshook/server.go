package gogshook

import (
	"fmt"
	"log"
	"net/http"

	gogs "github.com/gogits/go-gogs-client"
	"gitlab.com/pongsatt/githook/pkg/tekton"
)

const (
	gogsHeaderEvent = "Gogs-Event"
)

// GogsReceiveAdapter converts incoming Gogs webhook events to
// CloudEvents and then sends them to the specified Sink
type GogsReceiveAdapter struct {
	Client       *http.Client
	TektonClient *tekton.Client

	Namespace   string
	Name        string
	RunSpecJSON string
}

// HandleEvent is invoked whenever an event comes in from Gogs
func (ra *GogsReceiveAdapter) HandleEvent(payload interface{}, header http.Header) {
	err := ra.handleEvent(payload, header)
	if err != nil {
		log.Printf("unexpected error handling git event: %s", err)
	}
}

func (ra *GogsReceiveAdapter) handleEvent(payload interface{}, header http.Header) error {
	gogsEventType := header.Get("X-" + gogsHeaderEvent)

	log.Printf("Handling %s", gogsEventType)

	if gogsEventType == "" {
		return fmt.Errorf("invalid event: %s", gogsEventType)
	}

	options := buildOptionFromPayload(payload)
	options.Namespace = ra.Namespace
	options.Prefix = ra.Name
	options.RunSpecJSON = ra.RunSpecJSON

	pipelineRun, err := ra.TektonClient.CreatePipelineRun(options)

	if err != nil {
		return err
	}

	log.Printf("create pipeline run successfully %s", pipelineRun.Name)

	return nil
}

func buildOptionFromPayload(payload interface{}) tekton.PipelineOptions {
	switch payload.(type) {
	case gogs.CreatePayload:
		p := payload.(gogs.CreatePayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repo.HTMLURL,
			GitRevision: p.Ref,
		}
	case gogs.ReleasePayload:
		p := payload.(gogs.ReleasePayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Release.TargetCommitish,
		}
	case gogs.PushPayload:
		p := payload.(gogs.PushPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repo.HTMLURL,
			GitRevision: p.Ref,
		}
	case gogs.DeletePayload:
		p := payload.(gogs.DeletePayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repo.HTMLURL,
			GitRevision: p.Ref,
		}
	case gogs.ForkPayload:
		p := payload.(gogs.ForkPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repo.HTMLURL,
			GitRevision: p.Repo.DefaultBranch,
		}
	case gogs.IssuesPayload:
		p := payload.(gogs.IssuesPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Repository.DefaultBranch,
		}
	case gogs.IssueCommentPayload:
		p := payload.(gogs.IssueCommentPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Repository.DefaultBranch,
		}
	case gogs.PullRequestPayload:
		p := payload.(gogs.PullRequestPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.PullRequest.HeadBranch,
		}
	}
	return tekton.PipelineOptions{}
}
