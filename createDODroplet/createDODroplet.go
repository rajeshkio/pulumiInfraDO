package createDODroplet

import (
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateDODroplet(ctx *pulumi.Context, sshFingerprint, projectId string) ([]*digitalocean.Droplet, error) {
	dropletName := []string{"rajeshRancher-1", "rajeshRancher-2"}
	droplets := make([]*digitalocean.Droplet, 0)
	for _, name := range dropletName {
		d, err := digitalocean.NewDroplet(ctx, name, &digitalocean.DropletArgs{
			Image:  pulumi.String("ubuntu-23-10-x64"),
			Region: pulumi.String("blr1"),
			Size:   pulumi.String("s-4vcpu-8gb"),
			SshKeys: pulumi.StringArray{
				pulumi.String(sshFingerprint),
			},
			Tags: pulumi.StringArray{
				pulumi.String("rajesh-rancher-servers"),
				pulumi.String("project:" + projectId),
			},
		})
		if err != nil {
			return []*digitalocean.Droplet{}, err
		}
		droplets = append(droplets, d)
	}
	return droplets, nil
}
