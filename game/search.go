package game

import (
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type (
	Search struct {
		cancel chan struct{}
	}

	SearchCancelled struct {
		string
	}
)

func NewSearch() *Search {
	return &Search{
		cancel: make(chan struct{}, 1),
	}
}

func (error *SearchCancelled) Error() string {
	return "The search was cancelled"
}

func (search *Search) Cancel() {
	search.cancel <- struct{}{}
}

func (search *Search) Close() {
	close(search.cancel)
}

func (search *Search) Start() (*process.Process, error) {
	var application string = "LeagueClientUx.exe"

	timeout := time.After(30 * time.Second)
	ticker := time.Tick(1 * time.Second)

	for {
		select {
		case <-ticker:
			processes, _ := process.Processes()

			for _, process := range processes {
				name, _ := process.Name()

				if name == application {
					return process, nil
				}
			}
		case <-timeout:
			return nil, &ProcessNotFoundError{application}
		case <-search.cancel:
			search.Close()
			return nil, &SearchCancelled{}
		}
	}

	return nil, &ProcessNotFoundError{application}
}
