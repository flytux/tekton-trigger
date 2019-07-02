package githook

import (
	"fmt"

	"gitlab.com/pongsatt/githook/api/v1alpha1"
	"gitlab.com/pongsatt/githook/pkg/client"
	"gitlab.com/pongsatt/githook/pkg/model"
)

// ProjectHookClient provides webhook inteface
type ProjectHookClient interface {
	Create(options *model.HookOptions) (string, error)
	Update(options *model.HookOptions, hookID string) error
	Delete(options *model.HookOptions) error
}

// GitClient provides git client functionalities
type GitClient interface {
	Validate(options *model.HookOptions) (exists bool, changed bool, err error)
	Create(options *model.HookOptions) (string, error)
	Update(options *model.HookOptions) (string, error)
	Delete(options *model.HookOptions) error
}

// Client provides webhook client
type Client struct {
	GitClient GitClient
}

// New creates new client with dependencies
func New(gitProvider v1alpha1.GitProvider, baseURL, accessToken string) (*Client, error) {
	var gitClient GitClient

	switch gitProvider {
	case v1alpha1.Gogs:
		gitClient = client.NewGogsClient(baseURL, accessToken)
	}

	if gitClient == nil {
		return nil, fmt.Errorf("git provider %s not support", gitProvider)
	}

	return &Client{
		GitClient: gitClient,
	}, nil
}

// Create creates webhook
func (client Client) Create(options *model.HookOptions) (string, error) {
	return client.GitClient.Create(options)
}

// Update updates webhook
func (client Client) Update(options *model.HookOptions) (string, error) {
	return client.GitClient.Update(options)
}

// Validate checks if hook has been changed
func (client Client) Validate(options *model.HookOptions) (exists bool, changed bool, err error) {
	return client.GitClient.Validate(options)
}

// Delete webhook
func (client Client) Delete(options *model.HookOptions) error {
	return client.GitClient.Delete(options)
}
