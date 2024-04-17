package checkForSSHKeys

import (
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CheckForSSHKeys(ctx *pulumi.Context, sshKeyName string) (digitalocean.LookupSshKeyResult, error) {
	sshKey, err := digitalocean.LookupSshKey(ctx, &digitalocean.LookupSshKeyArgs{
		Name: sshKeyName,
	})
	if err != nil {
		return digitalocean.LookupSshKeyResult{}, err
	}
	return *sshKey, nil
}
