package gitlab

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	ggl "gitlab.com/gitlab-org/api/client-go"
)

type GitlabClient struct {
	AuthToken string `envconfig:"GITLAB_TOKEN"`
	BaseURL   string `envconfig:"GITLAB_URL"`

	api *ggl.Client
}

func NewGitlabClient() (*GitlabClient, error) {
	client := GitlabClient{}
	err := client.init()
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *GitlabClient) init() error {
	err := envconfig.Process("", c)
	if err != nil {
		return err
	}
	if c.AuthToken == "" || c.BaseURL == "" {
		return fmt.Errorf("AuthToken and BaseURL are required: %v", c)
	}

	client, err := ggl.NewClient(c.AuthToken, ggl.WithBaseURL(fmt.Sprintf("%s/api/v4", c.BaseURL)))
	if err != nil {
		return err
	}
	c.api = client
	return nil
}
