//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package main

import (
	"fmt"
	"github.com/energye/energy/v2/cmd/internal"
	"github.com/energye/energy/v2/cmd/internal/command"
	"github.com/jessevdk/go-flags"
	"os"
)

var commands = []*command.Command{
	nil,
	internal.CmdInstall,
	internal.CmdPackage,
	internal.CmdVersion,
	internal.CmdSetenv,
	internal.CmdEnv,
	internal.CmdInit,
	internal.CmdBuild,
}

func main() {
	wd, _ := os.Getwd()
	cc := &command.Config{Wd: wd}
	parser := flags.NewParser(cc, flags.HelpFlag|flags.PassDoubleDash)
	if len(os.Args) < 2 {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}
	if extraArgs, err := parser.ParseArgs(os.Args[1:]); err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	} else {
		switch parser.Active.Name {
		case "install":
			cc.Index = 1
		case "package":
			cc.Index = 2
		case "version":
			cc.Index = 3
		case "setenv":
			cc.Index = 4
		case "env":
			cc.Index = 5
		case "init":
			cc.Index = 6
		case "build":
			cc.Index = 7
		}
		command := commands[cc.Index]
		if len(extraArgs) < 1 || extraArgs[len(extraArgs)-1] != "." {
			fmt.Fprintf(os.Stderr, "%s\n%s", command.UsageLine, command.Long)
			os.Exit(1)
		}
		fmt.Println("Energy executing:", command.Short)
		if err := command.Run(cc); err != nil {
			fmt.Fprint(os.Stderr, err.Error()+"\n")
			os.Exit(1)
		}
	}
}
