package main

import (
	"fmt"
	"log"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/rajeshkio/pulumiDigitalocean/checkForDOProject"
	"github.com/rajeshkio/pulumiDigitalocean/checkForSSHKeys"
	"github.com/rajeshkio/pulumiDigitalocean/createDODroplet"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		projectName := "Rancher Support"
		sshkeyName := "rajesh-keys"

		project, err := checkForDOProject.CheckForDOProject(ctx, projectName)
		if err != nil {
			log.Fatal("Cannot find the project", err)
		}

		sshKey, err := checkForSSHKeys.CheckForSSHKeys(ctx, sshkeyName)
		if err != nil {
			log.Fatal("Cannot find the sshkeys", err)
		}

		droplets, err := createDODroplet.CreateDODroplet(ctx, sshKey.Fingerprint, project.Id)
		if err != nil {
			log.Fatal("Failed to create droplet ", err)
		}

		for i, d := range droplets {
			ctx.Export(fmt.Sprintf("ssh server %d", i), pulumi.Sprintf("ssh root@%s", d.Ipv4Address))
		}
		return nil
	})
}
