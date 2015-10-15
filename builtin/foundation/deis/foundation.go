package deis

import (
  "fmt"
  "os"
  "os/exec"

  "github.com/hashicorp/otto/directory"
  "github.com/hashicorp/otto/foundation"
  "github.com/hashicorp/otto/helper/bindata"
  "github.com/hashicorp/otto/helper/compile"
  // "github.com/hashicorp/otto/helper/terraform"
  ottoExec "github.com/hashicorp/otto/helper/exec"
)

//go:generate go-bindata -pkg=deis -nomemcopy -nometadata ./data/...

// Foundation is an implementation of foundation.Foundation
type Foundation struct{}

func (f *Foundation) Compile(ctx *foundation.Context) (*foundation.CompileResult, error) {
  var opts compile.FoundationOptions
  opts = compile.FoundationOptions{
    Ctx: ctx,
    Bindata: &bindata.Data{
      Asset:    Asset,
      AssetDir: AssetDir,
    },
  }

  return compile.Foundation(&opts)
}

func (f *Foundation) Infra(ctx *foundation.Context) error {
  if ctx.Action == "" {
    appInfra := ctx.Appfile.ActiveInfrastructure()
    lookup := directory.Lookup{Infra: appInfra.Type}
    infra, err := ctx.Directory.GetInfra(&directory.Infra{Lookup: lookup})
    os.Setenv("DEISCTL_TUNNEL", infra.Outputs["ip"])

    fmt.Println("DEISCTL_TUNNEL is " + os.Getenv("DEISCTL_TUNNEL"))

    cmd := exec.Command("deisctl", "config", "platform", "set", "version=v1.11.1")
    cmd.Env = os.Environ()
    err = ottoExec.Run(ctx.Ui, cmd)

    cmd = exec.Command("deisctl", "config", "platform", "set", "sshPrivateKey=~/.ssh/deis-test")
    cmd.Env = os.Environ()
    err = ottoExec.Run(ctx.Ui, cmd)

    cmd = exec.Command("deisctl", "config", "platform", "set", "domain=goings.space")
    cmd.Env = os.Environ()
    err = ottoExec.Run(ctx.Ui, cmd)

    cmd = exec.Command("deisctl", "install", "platform")
    cmd.Env = os.Environ()
    err = ottoExec.Run(ctx.Ui, cmd)

    cmd = exec.Command("deisctl", "start", "platform")
    cmd.Env = os.Environ()
    err = ottoExec.Run(ctx.Ui, cmd)

    return err
  }
  return nil
}
