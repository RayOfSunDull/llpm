package main

import (
	"fmt"
	"llpm/src/lib"
	"os"
	"path/filepath"
)

func printAndExit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func exitIfErr(err error) {
	if err != nil {
		printAndExit(fmt.Sprint(err))
	}
}

func main() {
	homeDir, err := os.UserHomeDir()
	exitIfErr(err)

	configPath := filepath.Join(homeDir, ".config/llpm/config.json")
	activePath := "/tmp/llpm_active_processes.json"

	procMan, err := llpm.NewProcessManager(configPath, activePath)
	defer procMan.Save()

	argParser, err := llpm.NewArgParser(os.Args)
	exitIfErr(err)

	mainCommand := argParser.GetCommand()
	if mainCommand == "start" {
		alias := argParser.GetNext("alias")
		argParser.ExitIfErr()

		commandList := argParser.GetRemaining()
		if len(commandList) > 0 {
			procMan.StartExplicit(alias, commandList, make(map[string]string))
			procMan.AddAliasCommand(alias, commandList)
		} else {
			procMan.StartAuto(alias)
		}
	} else if mainCommand == "stop" {
		alias := argParser.GetNext("alias")
		argParser.ExitIfErr()

		err = procMan.Stop(alias)
		exitIfErr(err)
	} else if mainCommand == "toggle" {
		alias := argParser.GetNext("alias")
		argParser.ExitIfErr()

		err = procMan.Toggle(alias)
		exitIfErr(err)
	} else if mainCommand == "restart" {
		alias := argParser.GetNext("alias")
		argParser.ExitIfErr()

		procMan.Stop(alias)
		procMan.StartAuto(alias)
	} else if mainCommand == "list" {
		subCommand := argParser.GetNextOr("subcommand", "active")

		if subCommand == "active" {
			procMan.UpdateActives()
			for alias, pids := range procMan.Active {
				fmt.Printf("%s (PIDs: %v)\n", alias, pids)
			}
		} else if subCommand == "aliases" {
			for _, alias := range procMan.GetAliases() {
				aliasConfig, err := procMan.GetAliasConfig(alias)
				if err != nil {
					continue
				}
				fmt.Printf("%s:%s\n", alias, aliasConfig.Describe())
			}
		} else {
			fmt.Printf("Subcommand \"%s\" not recognised\n", subCommand)
		}
	} else if mainCommand == "add" {
		alias := argParser.GetNext("alias")
		argParser.ExitIfErr()

		commandList := argParser.GetRemaining()
		procMan.AddAliasCommand(alias, commandList)
	} else if mainCommand == "setenv" {
		alias := argParser.GetNext("alias")
		envKey := argParser.GetNext("key")
		envVal := argParser.GetNext("value")
		argParser.ExitIfErr()

		aliasConfig, _ := procMan.GetAliasConfig(alias)
		aliasConfig.SetEnvVar(envKey, envVal)
	} else if mainCommand == "help" || mainCommand == "--help" || mainCommand == "-h" {
		llpm.PrintHelp()
	} else {
		fmt.Printf(
			"Command \"%s\" not recognised. Use \"llpm help\" for a list of available commands\n",
			mainCommand)
	}
}
