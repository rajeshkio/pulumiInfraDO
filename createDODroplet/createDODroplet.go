package createDODroplet

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ModifyUserData(userData pulumi.String, isMaster bool, masterIP, nodeToken pulumi.StringOutput) pulumi.StringOutput {

	if isMaster {
		return userData.ToStringOutput()
	}
	return pulumi.All(masterIP, nodeToken, userData).ApplyT(func(args []interface{}) (string, error) {
		masterIP := args[0].(string)
		nodeToken := args[1].(string)
		userData := args[2].(string)

		userData = strings.ReplaceAll(userData, "MASTER_IP", masterIP)
		userData = strings.ReplaceAll(userData, "NODE_TOKEN", nodeToken)

		return userData, nil
	}).(pulumi.StringOutput)

}

func CreateDODroplet(ctx *pulumi.Context, sshFingerprint, projectId, dropletName, image, region, size string, userData pulumi.StringOutput) (*digitalocean.Droplet, error) {

	droplet, err := digitalocean.NewDroplet(ctx, dropletName, &digitalocean.DropletArgs{
		Name:   pulumi.String(dropletName),
		Image:  pulumi.String(image),
		Region: pulumi.String(region),
		Size:   pulumi.String(size),

		SshKeys: pulumi.StringArray{
			pulumi.String(sshFingerprint),
		},
		Tags: pulumi.StringArray{
			pulumi.String("rajesh-rancher-servers"),
			pulumi.String("project:" + projectId),
		},
		UserData: userData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create droplet: %v", err)
	}

	return droplet, nil
}

func CreateMasterNode(ctx *pulumi.Context, userData pulumi.StringOutput, sshFingerprint, projectId, image, dropletName, region, size string) (*digitalocean.Droplet, error) {
	// You need to specify image, size, and region for master node

	return CreateDODroplet(ctx, sshFingerprint, projectId, dropletName, image, region, size, userData)
}

func CreateWorkerNode(ctx *pulumi.Context, userData pulumi.StringOutput, sshFingerprint, projectId, image, dropletName, region, size string) (*digitalocean.Droplet, error) {
	// You need to specify image, size, and region for worker node

	return CreateDODroplet(ctx, sshFingerprint, projectId, dropletName, image, region, size, userData)
}
