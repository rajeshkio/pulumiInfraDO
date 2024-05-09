package checkRKE2Server

import (
	"os"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func readPrivateKey(privateKeyPath string) (string, error) {

	privateKeyByte, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}
	privateKey := string(privateKeyByte)
	return privateKey, nil
}
func createServerConnectionArgs(ip pulumi.StringOutput) (*remote.ConnectionArgs, error) {

	privateKeyPath := "/Users/rajeshkumar/.ssh/id_ed25519"
	privateKey, err := readPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}
	connectionArgs := &remote.ConnectionArgs{
		Host:           ip,
		User:           pulumi.String("root"),
		PrivateKey:     pulumi.String(privateKey),
		DialErrorLimit: pulumi.Int(300),
	}
	return connectionArgs, nil
}
func CheckRKE2Server(ctx *pulumi.Context, ip pulumi.StringOutput) (pulumi.StringOutput, error) {

	connectionArgs, err := createServerConnectionArgs(ip)
	if err != nil {
		return pulumi.StringOutput{}, err
	}
	cmd, err := remote.NewCommand(ctx, "waitForRKE2Start", &remote.CommandArgs{
		Create: pulumi.String(`#!/bin/bash
		PULUMI_COMMAND_STDOUT="1"; export PULUMI_COMMAND_STDOUT;
		until systemctl is-active --quiet rke2-server; do
          echo "Waiting for rke2-server service to be active..."
          sleep 10
        done
		echo "rke2-server service is active, userdata setup completed."
		`),
		Connection: connectionArgs,
	})
	if err != nil {
		return pulumi.StringOutput{}, err
	}
	RKE2ServerStatus := cmd.Stdout.ApplyT(func(stdout interface{}) (string, error) {
		return "Active", nil
	}).(pulumi.StringOutput)
	return RKE2ServerStatus, nil
}

func GetRKE2NodeToken(ctx *pulumi.Context, ip pulumi.StringOutput) (pulumi.StringOutput, error) {
	connectionArgs, err := createServerConnectionArgs(ip)
	if err != nil {
		return pulumi.StringOutput{}, err
	}
	nodeToken, err := remote.NewCommand(ctx, "getRKE2ServerToken", &remote.CommandArgs{
		Create:     pulumi.String("PULUMI_COMMAND_STDOUT=1 cat /var/lib/rancher/rke2/server/node-token"),
		Connection: connectionArgs,
	})
	if err != nil {
		return pulumi.StringOutput{}, err
	}

	return nodeToken.Stdout, nil

}

func GetKubeConfig(ctx *pulumi.Context, ip pulumi.StringOutput) (pulumi.StringOutput, error) {
	connectionArgs, err := createServerConnectionArgs(ip)
	if err != nil {
		return pulumi.StringOutput{}, err
	}
	privateKeyPath := "/Users/rajeshkumar/.ssh/id_ed25519"

	return pulumi.All(ip).ApplyT(func(args []interface{}) (string, error) {
		ip := args[0]
		remoteKubeconfigPath := "/etc/rancher/rke2/rke2.yaml"
		localKubeconfigPath := "./rke2.yaml"

		_, err = local.NewCommand(ctx, "copyKubeconfig", &local.CommandArgs{
			Create: pulumi.Sprintf("scp -o StrictHostKeyChecking=accept-new -i %s %s@%s:%s %s", privateKeyPath, connectionArgs.User, ip, remoteKubeconfigPath, localKubeconfigPath),
		})
		if err != nil {
			return "", err
		}
		return localKubeconfigPath, nil
	}).(pulumi.StringOutput), nil
}
