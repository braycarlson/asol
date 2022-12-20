package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/braycarlson/asol/authorization"
	"github.com/shirou/gopsutil/v3/process"
)

type (
	Game struct {
		process *process.Process
	}

	ProcessNotFoundError struct {
		Process string
	}
)

func NewGame(process *process.Process) *Game {
	return &Game{
		process: process,
	}
}

func (game *Game) Authorization() *authorization.Authorization {
	flag := game.Flag()

	return &authorization.Authorization{
		Username: flag["username"],
		Password: flag["remoting-auth-token"],
		Name:     flag["ux-name"],
		App:      flag["app-port"],
		Region:   flag["region"],
		PID:      flag["app-pid"],
		Port:     flag["app-port"],
		Respawn:  flag["respawn-command"],
	}
}

func (error *ProcessNotFoundError) Error() string {
	return fmt.Sprintf("%s could not be found", error.Process)
}

func (game *Game) Process() *process.Process {
	return game.process
}

func (game *Game) Flag() map[string]string {
	arguments, err := game.process.CmdlineSlice()

	if err != nil {
		return nil
	}

	flags := make(map[string]string)

	for _, argument := range arguments {
		argument, err = strconv.Unquote(argument)
		flag := strings.Split(argument, "=")

		if len(flag) == 1 {
			continue
		}

		key := flag[0][2:]
		value := flag[1]

		flags[key] = value
	}

	flags["username"] = "riot"
	return flags
}
