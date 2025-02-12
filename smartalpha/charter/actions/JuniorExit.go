package actions

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/qtrix/simulator-wavespread/smartalpha"
)

type JuniorExit struct {
	Action string
	User   string
	ExitAt time.Time
}

func (s *JuniorExit) ExecuteAt() time.Time {
	return s.ExitAt
}

func (s *JuniorExit) Execute(sa *smartalpha.SmartAlpha) (interface{}, error) {
	logrus.WithFields(logrus.Fields{}).Info("-> junior exit")

	err := sa.ExitJunior(s.User)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
