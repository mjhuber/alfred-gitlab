package main

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

func search(cfg *Options, searchQuery string) {
	showUpdateStatus()

	git, err := gitlab.NewClient(cfg.Token, gitlab.WithBaseURL(fmt.Sprintf("%s/api/v4", cfg.BaseURL)))
	if err != nil {
		wf.Fatalf("Error connecting to gitlab: %s", err)
		return
	}
	opt := &gitlab.ListProjectsOptions{Search: gitlab.String(searchQuery)}
	projects, _, err := git.Projects.ListProjects(opt)
	if err != nil {
		wf.Warn("Error getting projects", err.Error())
		return
	}

	for _, project := range projects {
		wf.NewItem(fmt.Sprintf("%s/%s", project.PathWithNamespace, project.Name)).Valid(true).Arg(project.WebURL).Subtitle(project.Description)
	}
}
