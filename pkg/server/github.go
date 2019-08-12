package server

import (
	"net/http"

	"gitlab.com/pongsatt/githook/pkg/tekton"
	"gopkg.in/go-playground/webhooks.v5/github"
)

const (
	githubHeaderEvent = "GitHub-Event"
)

// GithubServer provides github git functionalities
type GithubServer struct {
	hook *github.Webhook
}

// NewGithubServer creates new github provider
func NewGithubServer(secretToken string) (*GithubServer, error) {
	hook, err := github.New(github.Options.Secret(secretToken))

	if err != nil {
		return nil, err
	}

	return &GithubServer{hook}, nil
}

// GetEventHeader returns github event header
func (git *GithubServer) GetEventHeader() string {
	return githubHeaderEvent
}

// Parse returns github payload
func (git *GithubServer) Parse(r *http.Request) (interface{}, error) {
	return git.hook.Parse(r,
		github.CreateEvent,
		github.DeleteEvent,
		github.ForkEvent,
		github.PushEvent,
		github.IssuesEvent,
		github.IssueCommentEvent,
		github.PullRequestEvent,
		github.ReleaseEvent)
}

// BuildOptionFromPayload builds pipeline option from payload information
func (git *GithubServer) BuildOptionFromPayload(payload interface{}) tekton.PipelineOptions {
	switch payload.(type) {
	case github.CreatePayload:
		p := payload.(github.CreatePayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Ref,
		}
	case github.ReleasePayload:
		p := payload.(github.ReleasePayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Release.TargetCommitish,
		}
	case github.PushPayload:
		p := payload.(github.PushPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Ref,
			GitCommit:   p.After,
		}
	case github.DeletePayload:
		p := payload.(github.DeletePayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Ref,
		}
	case github.ForkPayload:
		p := payload.(github.ForkPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Repository.DefaultBranch,
		}
	case github.IssuesPayload:
		p := payload.(github.IssuesPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Repository.DefaultBranch,
		}
	case github.IssueCommentPayload:
		p := payload.(github.IssueCommentPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.Repository.DefaultBranch,
		}
	case github.PullRequestPayload:
		p := payload.(github.PullRequestPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Repository.HTMLURL,
			GitRevision: p.PullRequest.Head.Ref,
		}
	}
	return tekton.PipelineOptions{}
}
