package llpm

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type AliasConfig struct {
	Env   map[string]string
	Start [][]string
}

type ProcessManager struct {
	configPath string
	activePath string
	Config     map[string]AliasConfig
	Active     map[string]([]int)
}

func NewAliasConfig() AliasConfig {
	return AliasConfig{
		Env:   make(map[string]string),
		Start: make([][]string, 0)}
}

func NewProcessManager(configPath string, activePath string) (ProcessManager, error) {
	config := make(map[string]AliasConfig)
	active := make(map[string]([]int))
	result := ProcessManager{
		configPath: configPath,
		activePath: activePath,
		Config:     config,
		Active:     active}

	err := LoadJson(configPath, &config)
	result.Config = config
	if err != nil {
		return result, err
	}

	err = LoadJson(activePath, &active)
	result.Active = active
	return result, err
}

func (procMan *ProcessManager) Save() (error, error) {
	errConfig := SaveJson(procMan.configPath, procMan.Config)
	errActive := SaveJson(procMan.activePath, procMan.Active)
	return errConfig, errActive
}

func (procMan *ProcessManager) GetActive(alias string) ([]int, error) {
	activePids, ok := procMan.Active[alias]
	if !ok || len(activePids) == 0 {
		return make([]int, 0), errors.New(fmt.Sprintf(
			"Process with alias %s not found", alias))
	}
	return activePids, nil
}

func (procMan *ProcessManager) SetActive(alias string, pids []int) {
	procMan.Active[alias] = pids
}

func (procMan *ProcessManager) AddActive(alias string, pid int) {
	activePids, _ := procMan.GetActive(alias)
	activePids = append(activePids, pid)
	procMan.SetActive(alias, activePids)
}

func (procMan *ProcessManager) DeleteActive(alias string) {
	delete(procMan.Active, alias)
}

func (procMan *ProcessManager) UpdateActives() (error, error) {
	for alias, pids := range procMan.Active {
		for i, pid := range pids {
			proc, err := os.FindProcess(pid)
			if err = proc.Signal(syscall.Signal(0)); err != nil {
				pids[i] = pids[len(pids)-1]
				pids = pids[:(len(pids) - 1)]
			}
		}
		procMan.SetActive(alias, pids)
		if len(pids) == 0 {
			procMan.DeleteActive(alias)
		}
	}
	return procMan.Save()
}

func (procMan *ProcessManager) GetAliases() []string {
	result := make([]string, 0, len(procMan.Config))
	for alias := range procMan.Config {
		result = append(result, alias)
	}
	return result
}

func (procMan *ProcessManager) GetAliasConfig(alias string) (AliasConfig, error) {
	aliasConfig, ok := procMan.Config[alias]
	if !ok {
		return NewAliasConfig(), errors.New(fmt.Sprintf(
			"Alias %s not found", alias))
	}
	return aliasConfig, nil
}

func (procMan *ProcessManager) SetAliasConfig(alias string, aliasConfig AliasConfig) {
	procMan.Config[alias] = aliasConfig
}

func (procMan *ProcessManager) GetAliasCommands(alias string) ([][]string, error) {
	aliasConfig, err := procMan.GetAliasConfig(alias)
	if err != nil {
		return make([][]string, 0), err
	}
	return aliasConfig.Start, nil
}

func (procMan *ProcessManager) SetAliasCommands(alias string, commandLists [][]string) {
	aliasConfig, _ := procMan.GetAliasConfig(alias)
	aliasConfig.Start = commandLists
	procMan.SetAliasConfig(alias, aliasConfig)
}

func (procMan *ProcessManager) AddAliasCommand(alias string, commandList []string) {
	commandLists, _ := procMan.GetAliasCommands(alias)
	commandLists = append(commandLists, commandList)
	procMan.SetAliasCommands(alias, commandLists)
}

func (procMan *ProcessManager) StartExplicit(alias string, commandList []string, environment map[string]string) {
	evalCommandList := make([]string, len(commandList))
	for i, element := range commandList {
		if len(element) > 0 && element[0] == '$' {
			key := element[1:]
			val, ok := environment[key]
			if ok {
				evalCommandList[i] = val
			}
		} else {
			evalCommandList[i] = element
		}
	}

	command, args := evalCommandList[0], evalCommandList[1:]
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	for key, val := range environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, val))
	}
	cmd.Start()
	procMan.AddActive(alias, cmd.Process.Pid)
}

func (procMan *ProcessManager) StartAuto(alias string) error {
	aliasConfig, err := procMan.GetAliasConfig(alias)
	if err != nil {
		return err
	}

	commandLists, env := aliasConfig.Unpack()

	for _, commandList := range commandLists {
		procMan.StartExplicit(alias, commandList, env)
	}
	return nil
}

func (procMan *ProcessManager) StopPids(pids []int) {
	for _, pid := range pids {
		proc, err := os.FindProcess(pid)
		if err != nil {
			continue
		}
		proc.Kill()
	}
}

func (procMan *ProcessManager) Stop(alias string) error {
	pids, err := procMan.GetActive(alias)
	if err != nil {
		return err
	}
	procMan.DeleteActive(alias)
	procMan.StopPids(pids)
	return nil
}

func (procMan *ProcessManager) Toggle(alias string) error {
	pids, err := procMan.GetActive(alias)
	if err != nil { // error means the processes do not exist; in that case, start them
		return procMan.StartAuto(alias)
	} else {
		procMan.DeleteActive(alias)
		procMan.StopPids(pids)
		return nil
	}
}

func (aliasConfig *AliasConfig) GetCommands() [][]string {
	return aliasConfig.Start
}

func (aliasConfig *AliasConfig) SetCommands(commandLists [][]string) {
	aliasConfig.Start = commandLists
}

func (aliasConfig *AliasConfig) Unpack() ([][]string, map[string]string) {
	return aliasConfig.Start, aliasConfig.Env
}

func (aliasConfig *AliasConfig) GetEnvVar(key string) (string, bool) {
	result, ok := aliasConfig.Env[key]
	return result, ok
}

func (aliasConfig *AliasConfig) SetEnvVar(key string, val string) {
	aliasConfig.Env[key] = val
}

func (aliasConfig *AliasConfig) Describe() string {
	commandLists, env := aliasConfig.Unpack()
	commandBlurb := ""
	for _, commandList := range commandLists {
		commandBlurb += "\n    " + strings.Join(commandList, " ")
	}

	envBlurb := ""
	if len(env) != 0 {
		envBlurb = "\n    environment variables:"
		for key, val := range env {
			envBlurb += fmt.Sprintf("\n        %s=%s", key, val)
		}
	}
	return commandBlurb + envBlurb
}
