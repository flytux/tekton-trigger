package server

import (
	"net/http"

	"gitlab.com/pongsatt/githook/pkg/tekton"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

const (
	gitlabHeaderEvent = "Gitlab-Event"
)

// GitlabServer provides gitlab git functionalities
type GitlabServer struct {
	hook *gitlab.Webhook
}

// NewGitlabServer creates new gitlab provider
func NewGitlabServer(secretToken string) (*GitlabServer, error) {
	hook, err := gitlab.New(gitlab.Options.Secret(secretToken))

	if err != nil {
		return nil, err
	}

	return &GitlabServer{hook}, nil
}

// GetEventHeader returns gitlab event header
func (git *GitlabServer) GetEventHeader() string {
	return gitlabHeaderEvent
}

// Parse returns gitlab payload
func (git *GitlabServer) Parse(r *http.Request) (interface{}, error) {
	return git.hook.Parse(r,
		gitlab.PushEvents,
		gitlab.IssuesEvents,
		gitlab.CommentEvents,
		gitlab.MergeRequestEvents)
}

// BuildOptionFromPayload builds pipeline option from payload information
func (git *GitlabServer) BuildOptionFromPayload(payload interface{}) tekton.PipelineOptions {
	switch payload.(type) {
	case gitlab.PushEventPayload:
		p := payload.(gitlab.PushEventPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Project.HTTPURL,
			GitRevision: p.Ref,
		}
	case gitlab.IssueEventPayload:
		p := payload.(gitlab.IssueEventPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Project.HTTPURL,
			GitRevision: p.Project.DefaultBranch,
		}
	case gitlab.CommentEventPayload:
		p := payload.(gitlab.CommentEventPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Project.HTTPURL,
			GitRevision: p.Project.DefaultBranch,
		}
	case gitlab.MergeRequestEventPayload:
		p := payload.(gitlab.MergeRequestEventPayload)
		return tekton.PipelineOptions{
			GitURL:      p.Project.HTTPURL,
			GitRevision: p.ObjectAttributes.SourceBranch,
		}
	}
	return tekton.PipelineOptions{}
}
