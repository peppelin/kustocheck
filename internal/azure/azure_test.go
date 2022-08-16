package azure

import (
	"context"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"kustocheck/internal/files"
	"testing"
)

func TestGetYAMLUrls(t *testing.T) {
	type args struct {
		repo      *files.Repo
		gitClient git.Client
		ctx       context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetYAMLUrls(tt.args.repo, tt.args.gitClient, tt.args.ctx)
		})
	}
}
