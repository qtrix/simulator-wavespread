package actions

import (
	"time"

	"github.com/qtrix/simulator-wavespread/smartalpha"
	"github.com/sirupsen/logrus"
)

type SeniorExit struct {
	Action string
	User   string
	ExitAt time.Time
}

func (s *SeniorExit) ExecuteAt() time.Time {
	return s.ExitAt
}

func (s *SeniorExit) Execute(sa *smartalpha.SmartAlpha) (interface{}, error) {
	logrus.WithFields(logrus.Fields{}).Info("-> senior exit")

	err := sa.ExitSenior(s.User)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
