package asol

import (
	"context"
	"encoding/csv"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type GameProcess struct {
	process   *process.Process
	name      string
	app       string
	region    string
	username  string
	password  string
	pid       string
	port      string
	directory string
	respawn   string
}

func NewGameProcess() *GameProcess {
	process, _ := GetProcessIndefinitely("LeagueClientUx.exe")
	flags := Flags(process)

	return &GameProcess{
		process:   process,
		username:  "riot",
		password:  flags["remoting-auth-token"],
		name:      flags["ux-name"],
		app:       flags["app-port"],
		region:    flags["region"],
		pid:       flags["app-pid"],
		port:      flags["app-port"],
		directory: flags["install-directory"],
		respawn:   flags["respawn-command"],
	}
}

func (game *GameProcess) Username() string {
	return game.username
}

func (game *GameProcess) Password() string {
	return game.password
}

func (game *GameProcess) Port() string {
	return game.port
}

func (game *GameProcess) Path() string {
	return filepath.Join(game.directory, game.respawn)
}

func Flags(process *process.Process) map[string]string {
	cmdline, err := process.Cmdline()

	if err != nil {
		return nil
	}

	reader := csv.NewReader(
		strings.NewReader(cmdline),
	)
	reader.Comma = ' '
	arguments, _ := reader.Read()

	flags := make(map[string]string)

	for _, argument := range arguments {
		if strings.HasPrefix(argument, "--") && strings.Contains(argument, "=") {
			argument := strings.Split(argument[2:], "=")
			flags[argument[0]] = argument[1]
		}
	}

	return flags
}

func GetProcess(application string) (*process.Process, error) {
	processes, _ := process.Processes()

	for _, process := range processes {
		name, _ := process.Name()

		if name == application {
			return process, nil
		}
	}

	return nil, &ProcessNotFoundError{application}
}

func GetProcessIndefinitely(application string) (*process.Process, error) {
	for {
		processes, _ := process.Processes()

		for _, process := range processes {
			name, _ := process.Name()

			if name == application {
				return process, nil
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}

	return nil, &ProcessNotFoundError{application}
}

func IsGameOrClient(channel chan bool, game string, client string) {
	var wg sync.WaitGroup
	wg.Add(3)

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	go func() {
		defer wg.Done()

		for {
			if len(channel) == cap(channel) {
				break
			}

			process, _ := GetProcess(game)

			if process != nil {
				channel <- true
				break
			}

			time.Sleep(1000 * time.Millisecond)
		}

		select {
		case <-ctx.Done():
			return
		}
	}()

	go func() {
		defer wg.Done()

		for {
			if len(channel) == cap(channel) {
				break
			}

			process, err := GetProcess(client)

			if process == nil || err != nil {
				channel <- false
				break
			}

			time.Sleep(1000 * time.Millisecond)
		}

		select {
		case <-ctx.Done():
			return
		}
	}()

	go func() {
		defer wg.Done()
		cancel()
	}()

	wg.Wait()
}
