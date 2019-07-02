package gogshook

import (
	"fmt"
	"strconv"

	gogs "github.com/gogits/go-gogs-client"
)

// HookOptions keeps webhook options
type HookOptions struct {
	AccessToken string
	SecretToken string
	Project     string
	ID          string
	BaseURL     string
	URL         string
	Owner       string
	Events      []string
}

// ProjectHookClient provides webhook inteface
type ProjectHookClient interface {
	Create(options *HookOptions) (string, error)
	Delete(options *HookOptions, hookID string) error
}

// Client provides webhook client
type Client struct {
	GogsClient *gogs.Client
}

// New creates new client with dependencies
func New(baseURL, accessToken string) *Client {
	gogsClient := gogs.NewClient(baseURL, accessToken)

	return &Client{
		GogsClient: gogsClient,
	}
}

// Create creates webhook
func (client Client) Create(options *HookOptions) (string, error) {
	gogsClient := client.GogsClient

	hookID, err := createHook(gogsClient, options)

	if err != nil {
		return "", fmt.Errorf("fail to create new hook: %s", err)
	}

	return strconv.Itoa(int(hookID)), nil
}

// Update updates webhook
func (client Client) Update(options *HookOptions) (string, error) {
	gogsClient := client.GogsClient

	if options.ID == "" {
		return "", fmt.Errorf("webhook id is required to be updated")
	}

	hookID, err := updateHook(gogsClient, options)

	return strconv.Itoa(hookID), err
}

func createHook(gogsClient *gogs.Client, options *HookOptions) (int64, error) {
	hookOptions := gogs.CreateHookOption{
		Active: true,
		Config: map[string]string{
			"content_type": "json",
			"url":          options.URL,
			"secret":       options.SecretToken,
		},
		Events: options.Events,
		Type:   "gogs",
	}

	hook, err := gogsClient.CreateRepoHook(options.Owner, options.Project, hookOptions)
	if err != nil {
		return -1, fmt.Errorf("Failed to add webhook to the Project:" + options.Project + " due to " + err.Error())
	}

	return hook.ID, nil
}

// Validate checks if hook has been changed
func (client Client) Validate(options *HookOptions) (exists bool, changed bool, err error) {
	if options.ID == "" {
		return false, false, nil
	}

	hook, err := getHook(client.GogsClient, options)

	if err != nil {
		return false, false, err
	}

	if hook == nil {
		return false, false, nil
	}

	if hook.Config["url"] != options.URL {
		return true, true, nil
	}

	if len(hook.Events) != len(options.Events) {
		return true, true, nil
	}

	eventSet := make(map[string]bool)

	for _, event := range hook.Events {
		eventSet[event] = true
	}

	for _, event := range options.Events {
		if eventSet[event] == false {
			return true, true, nil
		}
	}

	return true, false, nil
}

func getHook(gogsClient *gogs.Client, options *HookOptions) (*gogs.Hook, error) {
	hooks, err := gogsClient.ListRepoHooks(options.Owner, options.Project)

	if err != nil {
		return nil, fmt.Errorf("Failed to list webhook to the Project:" + options.Project + " due to " + err.Error())
	}

	for _, hook := range hooks {
		if strconv.Itoa(int(hook.ID)) == options.ID {
			return hook, nil
		}
	}

	return nil, nil
}

func updateHook(gogsClient *gogs.Client, options *HookOptions) (int, error) {
	active := true

	hookOptions := gogs.EditHookOption{
		Active: &active,
		Config: map[string]string{
			"content_type": "json",
			"url":          options.URL,
			"secret":       options.SecretToken,
		},
		Events: options.Events,
	}

	hookID, err := strconv.Atoi(options.ID)

	if err != nil {
		return -1, fmt.Errorf("cannot convert hook ID %v", hookID)
	}

	err = gogsClient.EditRepoHook(options.Owner, options.Project, int64(hookID), hookOptions)

	if err != nil {
		return -1, fmt.Errorf("Failed to update webhook to the Project:" + options.Project + " due to " + err.Error())
	}

	return hookID, nil
}

// Delete webhook
func (client Client) Delete(options *HookOptions) error {
	if options.ID != "" {
		hookID, err := strconv.Atoi(options.ID)
		if err != nil {
			return fmt.Errorf("failed to convert hook id to int: " + err.Error())
		}

		gogsClient := gogs.NewClient(options.BaseURL, options.AccessToken)

		err = gogsClient.DeleteRepoHook(options.Owner, options.Project, int64(hookID))
		if err != nil {
			return fmt.Errorf("failed to delete hook owner '%s' project '%s' : %s", options.Owner, options.Project, err)
		}
	}

	return nil
}
