package createDODroplet

import (
	"fmt"
	"os"
	"strings"

	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateDODroplet(ctx *pulumi.Context, sshFingerprint, projectId, dropletName, userDataFilePath string, isMaster bool, masterIP, nodeToken pulumi.StringOutput) (*digitalocean.Droplet, error) {
	userData, err := os.ReadFile(userDataFilePath)
	//	var userDataWithToken pulumi.StringOutput
	fmt.Println(nodeToken)
	if err != nil {
		return nil,

			fmt.Errorf("failed to read user data file: %w", err)
	}

	if !isMaster {
		userDataWithMasterIP := masterIP.ApplyT(func(ip interface{}) (string, error) {
			ipStr, ok := ip.(string)
			if !ok {
				return "", fmt.Errorf("IP address is not a string")
			}
			return strings.ReplaceAll(string(userData), "MASTER_IP", ipStr), nil
		}).(pulumi.StringOutput)

		ctx.Export("userdatawithip", userDataWithMasterIP)
	}

	droplet, err := digitalocean.NewDroplet(ctx, dropletName, &digitalocean.DropletArgs{
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
		UserData: pulumi.String(userData),
	})
	if err != nil {
		return nil, err
	}
	return droplet, nil
}
