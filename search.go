package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	aw "github.com/deanishe/awgo"
	gitlab "github.com/xanzy/go-gitlab"
)

func find(query string) {
	projects := []*gitlab.Project{}
	if wf.Cache.Exists(cacheName) {
		log.Printf("cache exists")
		if err := wf.Cache.LoadJSON(cacheName, &projects); err != nil {
			wf.FatalError(err)
		}
	}

	// If the cache has expired, set Rerun (which tells Alfred to re-run the
	// workflow), and start the background update process if it isn't already
	// running.
	if wf.Cache.Expired(cacheName, maxCacheAge) {
		log.Printf("cache has expired")
		wf.Rerun(0.3)
		if !wf.IsRunning("download") {
			cmd := exec.Command(os.Args[0], "download")
			log.Printf("running download job")
			if err := wf.RunInBackground("download", cmd); err != nil {
				wf.FatalError(err)
			}
		} else {
			log.Printf("Download job already running")
		}
		// Cache is also "expired" if it doesn't exist. So if there are no
		// cached data, show a corresponding message and exit.
		if len(projects) == 0 {
			log.Printf("zero projects")
			wf.NewItem("Downloading projects...").Icon(aw.IconInfo)
			wf.SendFeedback()
			return
		}
	}

	for _, project := range projects {
		projID := fmt.Sprintf("%s/%s", project.PathWithNamespace, project.Name)
		wf.NewItem(project.Name).
			Valid(true).
			Arg(project.WebURL).
			Subtitle(project.PathWithNamespace).
			UID(projID)
	}

	res := wf.Filter(query)
	log.Printf("%d/%d projects match %q", len(res), len(projects), query)

	// Convenience method that shows a warning if there are no results to show.
	// Alfred's default behaviour if no results are returned is to show its
	// fallback searches, which is also what it does if a workflow errors out.
	//
	// As such, it's a good idea to display a message in this situation,
	// otherwise the user can't tell if the workflow failed or simply found
	// no matching results.
	wf.WarnEmpty("No repos found", "Try a different query?")
	wf.SendFeedback()
}

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

func cacheRepos(cfg *Options) {
	git, err := gitlab.NewClient(cfg.Token, gitlab.WithBaseURL(fmt.Sprintf("%s/api/v4", cfg.BaseURL)))
	if err != nil {
		wf.Fatalf("Error connecting to gitlab: %s", err)
		return
	}

	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	projects := []*gitlab.Project{}
	keepGoing := true
	for keepGoing {
		ps, resp, err := git.Projects.ListProjects(opt)
		if err != nil {
			wf.Fatalf("error getting projects: %v", err)
		}

		for _, p := range ps {
			projects = append(projects, p)
		}
		keepGoing = resp.NextPage > 0
		opt.Page = resp.NextPage
	}

	if err := wf.Cache.StoreJSON(cacheName, projects); err != nil {
		wf.FatalError(err)
	}
}
