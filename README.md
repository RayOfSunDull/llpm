# Long-Lived Process Manager

`llpm` is basically a bootleg version of `systemctl`. I would use `systemctl` but they refuse to implement a toggle command, which would be really convenient for me as a user. It also makes it cumbersome to run multiple long-lived commands in the same unit, which I also want to do. If this is not the case for you, you should consider using `systemctl` instead.

`llpm`'s main function is starting processes, keeping track of their PIDs and killing them on command.

## Installation (Build from source)

As it stands, `llpm` only works on Linux, but it should be able to be ported with minimal effort. To build from source, you need the (go)[https://go.dev/] compiler and GNU make:

```sh
git clone https://github.com/RayOfSunDull/llpm
cd llpm
make
make install
```

## Usage

`llpm`'s smallest organizational unit is the *alias*; Each alias is comprised of one or more commands which can be started and stopped by invoking the alias' name. The environment variables for all commands may also be set. You can:

### Add a command to an alias, creating it if it does not exist:
```sh
llpm add [alias] [commands...]
```

### Start the commands in the alias:
```sh
llpm start [alias]
```

### Stop the processes in the alias:
```sh
llpm stop [alias]
```

### Toggle the processes in the alias (if they are running, they will be stopped; if they are not, they will be started):
```sh
llpm toggle [alias]
```

### Set the value of an environment variable:
```sh
llpm setenv [alias] [key] [value]
```
Note that environment variables are replaced in the commands themselves. For example, if a command is `fswatch $MYDIRECTORY`, `llpm` will replace `$MYDIRECTORY` with the value of the `MYDIRECTORY` environment variable, if it exists.

Keeping track of logs is not supported, so if your processes crash, you may have a hard time figuring out why.

Use `llpm help` for more details.
