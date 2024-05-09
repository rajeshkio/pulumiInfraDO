package main

import (
	"fmt"
	"log"
	"os"

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
		//	RKE2AgentDroplet := "rajeshRancherAgent"
		userAgentDataFilePath := "createDODroplet/installRKE2Agent.sh"
		image := "ubuntu-23-10-x64"
		region := "blr1"
		size := "s-4vcpu-8gb"

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

		masterUserDataBytes, err := os.ReadFile(userMasterDataFilePath)
		if err != nil {
			return fmt.Errorf("failed to read masterUserData file: %v", err)
		}
		masterUserDataStr := pulumi.String(string(masterUserDataBytes))
		masterUserData := createDODroplet.ModifyUserData(masterUserDataStr, true, pulumi.StringOutput{}, pulumi.StringOutput{})
		masterDroplet, err := createDODroplet.CreateDODroplet(ctx, sshKey.Fingerprint, project.Id, RKE2MasterDroplet, image, region, size, masterUserData)
		if err != nil {
			return err
		}

		masterIP := masterDroplet.Ipv4Address
		ctx.Export("MasterIP", masterIP)

		RKE2ServerStatus, err := checkRKE2Server.CheckRKE2Server(ctx, masterIP)
		if err != nil {
			log.Printf("Error: Failed to start RKE2 server: %v\n", err)
			return err
		}
		_ = pulumi.All(RKE2ServerStatus, masterIP).ApplyT(func(args []interface{}) error {
			status := args[0].(string)
			if status == "Active" {
				token, err := checkRKE2Server.GetRKE2NodeToken(ctx, masterIP)
				if err != nil {
					return err
				}
				workerUserDataBytes, err := os.ReadFile(userAgentDataFilePath)
				if err != nil {
					return nil
				}
				workerUserDataStr := pulumi.String(string(workerUserDataBytes))
				workerCount := 1
				workerUserData := createDODroplet.ModifyUserData(workerUserDataStr, false, masterIP, token)
				for i := 0; i < workerCount; i++ {
					workerDroplet, err := createDODroplet.CreateDODroplet(ctx, sshKey.Fingerprint, project.Id, fmt.Sprintf("rajeshRancherAgent%d", i), image, region, size, workerUserData)
					if err != nil {
						return err
					}
					ctx.Export("MasterIP", workerDroplet.Ipv4Address)
				}
				return nil
			}
			return nil
		})

		kubeconfigOutput := RKE2ServerStatus.ApplyT(func(serverStatus interface{}) (pulumi.StringOutput, error) {
			kubeconfig, err := checkRKE2Server.GetKubeConfig(ctx, masterIP)
			if err != nil {
				log.Printf("Error: Cannot get kubeconfig file: %v\n", err)
				return pulumi.StringOutput{}, err
			}
			ctx.Export("kubeconfig: ", kubeconfig)
			return kubeconfig, nil
		})
		kubeconfigOutput.ApplyT(func(kubeconfig interface{}) error {
			kubeconfigBytes, ok := kubeconfig.([]byte)
			if !ok {
				return fmt.Errorf("kubeconfig is not of type []byte")
			}

			// Convert `[]byte` to `string` for display purposes and further processing
			kubeconfigStr := string(kubeconfigBytes)
			fmt.Println("Kubeconfig content:", kubeconfigStr)
			return nil
		})
		return nil
	})
}
