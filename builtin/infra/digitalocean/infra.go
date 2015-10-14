package digitalocean

import (
  // "os/exec"

  "github.com/hashicorp/otto/helper/bindata"
  "github.com/hashicorp/otto/helper/terraform"
  "github.com/hashicorp/otto/infrastructure"
  "github.com/hashicorp/otto/ui"
  // ottoExec "github.com/hashicorp/otto/helper/exec"
)

//go:generate go-bindata -pkg=digitalocean -nomemcopy -nometadata ./data/...

// Infra returns the infrastructure.Infrastructure implementation.
// This function is a infrastructure.Factory.
func Infra() (infrastructure.Infrastructure, error) {
  // add in make discovery-url here

  return &terraform.Infrastructure{
    CredsFunc:       creds,
    VerifyCredsFunc: verifyCreds,
    Bindata: &bindata.Data{
      Asset:    Asset,
      AssetDir: AssetDir,
    },
  }, nil
}

func verifyCreds(ctx *infrastructure.Context) error {
  return nil
}

func creds(ctx *infrastructure.Context) (map[string]string, error) {
  fields := []*ui.InputOpts{
    &ui.InputOpts{
      Id:          "do_token",
      Query:       "DigitalOcean token",
      Description: "DigitalOcean token used for API calls.",
      EnvVars:     []string{"DO_TOKEN"},
    },
    &ui.InputOpts{
      Id:          "ssh_keys",
      Query:       "SSH key fingerprint",
      Description: "SSH public key fingerprint that will be granted access to DigitalOcean instances",
      EnvVars:     []string{"DO_SSH_FINGERPRINT"},
    },
  }

  result := make(map[string]string, len(fields))
  for _, f := range fields {
    value, err := ctx.Ui.Input(f)
    if err != nil {
      return nil, err
    }

    result[f.Id] = value
  }

  return result, nil
}
