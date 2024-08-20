package llpm

import (
	"fmt"
	"strings"
)

func PrintHelp() {
	fmt.Println(strings.Join([]string{
		"llpm: a basic manager for long-lived processes",
		"usage: llpm [command] [args]",
		"commands:",
		"llpm add [alias] [command]",
		"    adds the specified command to an llpm alias.",
		"    An alias may contain multiple commands,",
		"    and they will all be started/stopped when necessary.",
		"    examples:",
		"        llpm add nightlight gammastep -O 4000K",
		"llpm start [alias] [command?]",
		"    starts the commands specified by the alias.",
		"    Any command specified after the alias is automatically",
		"    added to that alias.",
		"    examples:",
		"        llpm start nightlight",
		"        llpm start nightlight gammastep -O 4000K",
		"llpm stop [alias]",
		"    stops all currently running commands started",
		"    with this alias.",
		"    examples:",
		"        llpm stop nightlight",
		"llpm toggle [alias]",
		"    toggles the commands under the alias; if they are",
		"    running, it will stop them. If they are not, it",
		"    will start them",
		"    examples:",
		"        llpm toggle nightlight",
		"llpm restart [alias]",
		"    stops and then starts the commands under the alias",
		"    examples:",
		"        llpm restart nightlight",
		"llpm setenv [alias] [key] [value]",
		"    sets the value of one environment variable that",
		"    will be specified every time the commands under",
		"    the alias are run",
		"    examples:",
		"        llpm setenv nightlight HOME /home/other-user",
		"llpm list [active|aliases]",
		"    lists either all active commands or the available",
		"    aliases that have been configured (default=active)",
		"    examples:",
		"        llpm list aliases",
		"        llpm list active",
		"llpm help",
		"    prints the help message",
	}, "\n"))
}
