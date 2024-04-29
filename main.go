package main

import (
	"log"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/rajeshkio/pulumiDigitalocean/checkForDOProject"
	"github.com/rajeshkio/pulumiDigitalocean/checkForSSHKeys"
	"github.com/rajeshkio/pulumiDigitalocean/checkRKE2Server"
	"github.com/rajeshkio/pulumiDigitalocean/createDODroplet"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		projectName := "Rancher Support"
		sshkeyName := "rajesh-keys"
		RKE2MasterDroplet := "rajeshRancherMaster"
		userMasterDataFilePath := "createDODroplet/installRKE2.sh"
		RKE2AgentDroplet := "rajeshRancherAgent"
		userAgentDataFilePath := "createDODroplet/installRKE2Agent.sh"

		project, err := checkForDOProject.CheckForDOProject(ctx, projectName)
		if err != nil {
			log.Printf("Error: Cannot find the project: %v\n", err)
			return err
		}

		sshKey, err := checkForSSHKeys.CheckForSSHKeys(ctx, sshkeyName)
		if err != nil {
			log.Printf("Error: Cannot find the ssh keys: %v\n", err)
			return err
		}

		masterDroplet, err := createDODroplet.CreateDODroplet(ctx, sshKey.Fingerprint, project.Id, RKE2MasterDroplet, userMasterDataFilePath, true, pulumi.StringOutput{}, pulumi.StringOutput{})
		if err != nil {
			log.Printf("Error: Failed to create master droplet: %v\n", err)
			return err
		}
		ctx.Export("ssh server master", pulumi.Sprintf("ssh root@%s", masterDroplet.Ipv4Address))
		var nodeToken pulumi.StringOutput
		RKE2ServerStatus, err := checkRKE2Server.CheckRKE2Server(ctx, masterDroplet.Ipv4Address)
		if err != nil {
			log.Printf("Error: Failed to start RKE2 server: %v\n", err)
			return err
		}
		nodeToken = RKE2ServerStatus.ApplyT(func(status interface{}) (string, error) {
			if status == "RKE2 server is active" {
				token, err := checkRKE2Server.GetRKE2NodeToken(ctx, masterDroplet.Ipv4Address)
				if err != nil {
					return "", err
				}
				return token.ToStringOutput().ElementType().String(), nil
			}
			return "", nil
		}).(pulumi.StringOutput)
		ctx.Export("nodeToken", nodeToken)

		ctx.Export("nodeToken", nodeToken)
		masterIP := masterDroplet.Ipv4Address

		agentDroplet, err := createDODroplet.CreateDODroplet(ctx, sshKey.Fingerprint, project.Id, RKE2AgentDroplet, userAgentDataFilePath, false, masterIP, nodeToken)
		if err != nil {
			log.Printf("Error: Failed to create agent droplet: %v\n", err)
			return err
		} else {
			ctx.Export("ssh server agent", pulumi.Sprintf("ssh root@%s", agentDroplet.Ipv4Address))
		}

		return nil
	})
}
