package actions

import (
	"github.com/qtrix/simulator-wavespread/wavespread"
	"github.com/sirupsen/logrus"
	"time"
)

type SurferSwitchSideAnchor struct {
	Action     string
	User       string
	SwitchedAt time.Time
}

func (s *SurferSwitchSideAnchor) ExecuteAt() time.Time {
	return s.SwitchedAt
}

func (s *SurferSwitchSideAnchor) Execute(sa *wavespread.Wavespread) (interface{}, error) {
	logrus.WithFields(logrus.Fields{}).Info("-> surfer switch side anchor")

	err := sa.SwitchSideSurferToAnchor(s.User)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
