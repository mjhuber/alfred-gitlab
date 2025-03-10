package gitlab

import (
	ggl "gitlab.com/gitlab-org/api/client-go"
)

func (c *GitlabClient) GetProjects() ([]*ggl.Project, error) {
	projects := []*ggl.Project{}
	opt := &ggl.ListProjectsOptions{
		ListOptions: ggl.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	keepGoing := true
	for keepGoing {
		currProjects, resp, err := c.api.Projects.ListProjects(opt)
		if err != nil {
			return projects, err
		}
		projects = append(projects, currProjects...)
		keepGoing = resp.NextPage > 0
		opt.Page = resp.NextPage
	}
	return projects, nil
}
