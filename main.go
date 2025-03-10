package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mjhuber/alfred-gitlab/pkg/cache"
	"github.com/mjhuber/alfred-gitlab/pkg/gitlab"
)

var (
	repo        = "mjhuber/alfred-gitlab"
	maxCacheAge = 1440 // 24 hours
)

/*
ScriptFilterItem contains the structure for an Alfred script filter item
See https://www.alfredapp.com/help/workflows/inputs/script-filter/json/ for more information
*/
type ScriptFilterItem struct {
	Uid          string            `json:"uid,omitzero"`
	Title        string            `json:"title"`
	Subtitle     string            `json:"subtitle,omitzero"`
	Arg          []string          `json:"arg"`
	Icon         Icon              `json:"icon,omitzero"`
	Valid        bool              `json:"valid,default:true"`
	Match        string            `json:"match,omitzero"`
	Autocomplete string            `json:"autocomplete"`
	Type         string            `json:"type,omitzero"`
	Mods         map[string]Mods   `json:"mods,omitzero"`
	Action       Action            `json:"action,omitzero"`
	Text         map[string]string `json:"text,omitzero"`
	Quicklookurl string            `json:"quicklookurl,omitzero"`
}

type Icon struct {
	Type string `json:"type,omitzero"`
	Path string `json:"path"`
}

type Mods struct {
	Valid    bool   `json:"valid"`
	Arg      string `json:"arg"`
	Subtitle string `json:"subtitle"`
}

type Action struct {
	Text string `json:"text"`
	Url  string `json:"url,omitzero"`
	File string `json:"file,omitzero"`
	Auto string `json:"auto,omitzero"`
}

func main() {
	output := []ScriptFilterItem{}
	age, err := cache.FromCache("repos.json", &output)
	if err != nil {
		var glErr error
		output, glErr = getFromGitlab()
		if glErr != nil {
			log.Fatalf("error getting Gitlab projects: %s", glErr)
		}
		go cache.ToCache("repos.json", output)
	} else {
		if age > float64(maxCacheAge) {
			go updateCache()
		}
	}

	jsonData, err := json.MarshalIndent(map[string]interface{}{"items": output}, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %s", err)
	}
	fmt.Println(string(jsonData))
}

func getFromGitlab() ([]ScriptFilterItem, error) {
	output := []ScriptFilterItem{}
	client, err := gitlab.NewGitlabClient()
	if err != nil {
		return nil, fmt.Errorf("Error creating Gitlab client: %s", err)
	}

	projects, err := client.GetProjects()
	if err != nil {
		return nil, fmt.Errorf("Error getting projects: %s", err)
	}

	for _, project := range projects {
		sfObj := ScriptFilterItem{
			Uid:      project.PathWithNamespace,
			Title:    project.Name,
			Subtitle: project.PathWithNamespace,
			Arg:      []string{project.WebURL},
			Icon: Icon{
				Path: "icon.png",
			},
			Valid:        true,
			Match:        project.Name,
			Autocomplete: project.Name,
			Type:         "default",
			Mods: map[string]Mods{
				"alt": Mods{
					Valid:    true,
					Arg:      fmt.Sprintf("%s/-/pipelines", project.WebURL),
					Subtitle: "Open pipelines",
				},
				"cmd": Mods{
					Valid:    true,
					Arg:      fmt.Sprintf("%s/-/merge_requests", project.WebURL),
					Subtitle: "Open merge requests",
				},
			},
			Text: map[string]string{
				"copy":      project.WebURL,
				"largeType": project.PathWithNamespace,
			},
			Quicklookurl: project.WebURL,
		}
		output = append(output, sfObj)
	}
	return output, nil
}

func updateCache() {
	output, err := getFromGitlab()
	if err != nil {
		log.Printf("error getting Gitlab projects: %s", err)
		return
	}
	cache.ToCache("repos.json", output)
}
