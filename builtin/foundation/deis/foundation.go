package deis

import (
  "fmt"
  "os/exec"

  "github.com/hashicorp/otto/directory"
  "github.com/hashicorp/otto/foundation"
  "github.com/hashicorp/otto/helper/bindata"
  "github.com/hashicorp/otto/helper/compile"
  "github.com/hashicorp/otto/helper/terraform"
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
    foundationInfra, err := ctx.Directory.GetInfra(&directory.Infra{Lookup: lookup})
    fmt.Println(foundationInfra.Outputs["ip"])

    cmd := exec.Command("deisctl", "config", "platform", "set", "domain=goings.space")
    err = ottoExec.Run(ctx.Ui, cmd)

    cmd = exec.Command("deisctl", "install", "platform")
    err = ottoExec.Run(ctx.Ui, cmd)

    return err
  }
  return (&terraform.Foundation{}).Infra(ctx)
}
