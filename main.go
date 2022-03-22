package main

import (
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
)

var (
	// icons
	updateAvailable = &aw.Icon{Value: "icons/update-available.png"}

	repo  = "mjhuber/alfred-gitlab"
	query string

	// aw.Workflow is the main API
	wf *aw.Workflow
)

// Options contains options for connecting to the gitlab API
type Options struct {
	BaseURL string `env:"GITLAB_URL"`
	Token   string `env:"GITLAB_TOKEN"`
}

func init() {
	wf = aw.New(update.GitHub(repo), aw.HelpURL(repo+"/issues"))
}

func main() {
	wf.Run(run)
}

func run() {
	opts := &Options{}
	cfg := aw.NewConfig()
	if err := cfg.To(opts); err != nil {
		wf.Fatalf("Error loading variables: %v", err)
		return
	}

	query := strings.Join(os.Args[1:], " ")
	search(opts, query)
	wf.SendFeedback()
}
