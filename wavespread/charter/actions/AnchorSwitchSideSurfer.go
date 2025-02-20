package actions

import (
	"github.com/qtrix/simulator-wavespread/wavespread"
	"github.com/sirupsen/logrus"
	"time"
)

type AnchorSwitchSideSurfer struct {
	Action     string
	User       string
	SwitchedAt time.Time
}

func (s *AnchorSwitchSideSurfer) ExecuteAt() time.Time {
	return s.SwitchedAt
}

func (s *AnchorSwitchSideSurfer) Execute(sa *wavespread.Wavespread) (interface{}, error) {
	logrus.WithFields(logrus.Fields{}).Info("-> anchor switch side surfer")

	err := sa.SwitchSideAnchorToSurfer(s.User)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
