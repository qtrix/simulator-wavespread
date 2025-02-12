package actions

import (
	"time"

	"github.com/qtrix/simulator-wavespread/smartalpha"
)

type Action interface {
	ExecuteAt() time.Time
	Execute(sa *smartalpha.SmartAlpha) (interface{}, error)
}

type Actions []Action

func (a Actions) Len() int {
	return len(a)
}

func (a Actions) Less(i, j int) bool {
	return a[i].ExecuteAt().Before(a[j].ExecuteAt())
}

func (a Actions) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
