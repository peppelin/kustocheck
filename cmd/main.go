package main

import (
	"context"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"kustocheck/files"
	"log"
	"os"
)

var configFile = "config/repos.yaml"

const (
	organizationUrl = "https://dev.azure.com/schwarzit"
)

func main() {
	repos, err := files.ReadConfig(configFile)

	if err != nil {
		log.Fatal(err)
	}
	pat := os.Getenv("ADO_TOKEN")
	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(organizationUrl, pat)

	ctx := context.Background()

	// Create a new client to interact with git area
	gitClient, err := git.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}
	repos[0].GetInfo(gitClient, ctx)

}
