//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

//go:build linux
// +build linux

package build

import (
	"github.com/energye/energy/v2/cmd/internal/command"
	"github.com/energye/energy/v2/cmd/internal/project"
	"github.com/energye/energy/v2/cmd/internal/term"
	"github.com/energye/energy/v2/cmd/internal/tools"
	toolsCommand "github.com/energye/golcl/tools/command"
	"os"
	"strings"
)

func build(c *command.Config, proj *project.Project) (err error) {
	// go build
	cmd := toolsCommand.NewCMD()
	cmd.Dir = proj.ProjectPath
	cmd.IsPrint = false
	term.Section.Println("Building", proj.OutputFilename)
	var args = []string{"build"}
	if proj.TempDll {
		c.Build.Gtk = strings.ToLower(c.Build.Gtk)
		if c.Build.Gtk != "gtk3" && c.Build.Gtk != "gtk2" {
			term.Logger.Error("Compiling and enabling TempDll. gtk can only be gtk2 or gtk3")
			os.Exit(1)
		}
		args = append(args, "--tags=tempdll "+c.Build.Gtk)
	}
	args = append(args, "-ldflags", "-s -w")
	args = append(args, "-o", proj.OutputFilename)
	cmd.Command("go", args...)
	cmd.Command("strip", proj.OutputFilename)
	// upx
	if c.Build.Upx && tools.CommandExists("upx") {
		term.Section.Println("Upx compression")
		args = []string{"--best", "--no-color", "--no-progress", proj.OutputFilename}
		if c.Build.UpxFlag != "" {
			args = strings.Split(c.Build.UpxFlag, " ")
			args = append(args, proj.OutputFilename)
		}
		cmd.Command("upx", args...)
	}
	cmd.Close()
	if err == nil {
		term.Section.Println("Build Successfully")
	}
	return nil
}
