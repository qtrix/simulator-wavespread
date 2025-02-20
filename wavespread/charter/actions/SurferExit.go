package actions

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/qtrix/simulator-wavespread/wavespread"
)

type SurferExit struct {
	Action string
	User   string
	ExitAt time.Time
}

func (s *SurferExit) ExecuteAt() time.Time {
	return s.ExitAt
}

func (s *SurferExit) Execute(sa *wavespread.Wavespread) (interface{}, error) {
	logrus.WithFields(logrus.Fields{}).Info("-> surfer exit")

	err := sa.ExitSurfer(s.User)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
