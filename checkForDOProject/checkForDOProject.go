package checkForDOProject

import (
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CheckForDOProject(ctx *pulumi.Context, projectName string) (*digitalocean.LookupProjectResult, error) {

	project, err := digitalocean.LookupProject(ctx, &digitalocean.LookupProjectArgs{
		Name: &projectName,
	})
	if err != nil {
		return &digitalocean.LookupProjectResult{}, err
	}
	return project, nil
}
