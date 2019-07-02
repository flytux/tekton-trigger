package server

import (
	"net/http"

	gogsclient "github.com/gogits/go-gogs-client"
	"gitlab.com/pongsatt/githook/pkg/tekton"
	"gopkg.in/go-playground/webhooks.v5/gogs"
)

const (
	gogsHeaderEvent = "Gogs-Event"
)

// GogsServer provides gogs git functionalities
type GogsServer struct {
	hook *gogs.Webhook
}

// NewGogs creates new gogs provider
func NewGogsServer(secretToken string) (*GogsServer, error) {
	hook, err := gogs.New(gogs.Options.Secret(secretToken))

	if err != nil {
		return nil, err
	}

	return &GogsServer{hook}, nil
}

// GetEventHeader returns gogs event header
func (git *GogsServer) GetEventHeader() string {
	return gogsHeaderEvent
}

// Parse returns gogs payload
func (git *GogsServer) Parse(r *http.Request) (interface{}, error) {
	return git.hook.Parse(r,
		gogs.CreateEvent,
		gogs.DeleteEvent,
		gogs.ForkEvent,
		gogs.PushEvent,
		gogs.IssuesEvent,
		gogs.IssueCommentEvent,
		gogs.PullRequestEvent,
		gogs.ReleaseEvent)
}

// BuildOptionFromPayload builds pipeline option from payload information
func (git *GogsServer) BuildOptionFromPayload(payload interface{}) tekton.PipelineOptions {
	switch payload.(type) {
	case gogsclient.CreatePayload:
		p := payload.(gogsclient.CreatePayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repo.HTMLURL,
			GitRevision: p.Ref,
		}
	case gogsclient.ReleasePayload:
		p := payload.(gogsclient.ReleasePayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Release.TargetCommitish,
		}
	case gogsclient.PushPayload:
		p := payload.(gogsclient.PushPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repo.HTMLURL,
			GitRevision: p.Ref,
		}
	case gogsclient.DeletePayload:
		p := payload.(gogsclient.DeletePayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repo.HTMLURL,
			GitRevision: p.Ref,
		}
	case gogsclient.ForkPayload:
		p := payload.(gogsclient.ForkPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repo.HTMLURL,
			GitRevision: p.Repo.DefaultBranch,
		}
	case gogsclient.IssuesPayload:
		p := payload.(gogsclient.IssuesPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Repository.DefaultBranch,
		}
	case gogsclient.IssueCommentPayload:
		p := payload.(gogsclient.IssueCommentPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Repository.DefaultBranch,
		}
	case gogsclient.PullRequestPayload:
		p := payload.(gogsclient.PullRequestPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.PullRequest.HeadBranch,
		}
	}
	return tekton.PipelineOptions{}
}
