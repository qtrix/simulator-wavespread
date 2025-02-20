package actions

import (
	"time"

	"github.com/qtrix/simulator-wavespread/wavespread"
	"github.com/sirupsen/logrus"
)

type AnchorExit struct {
	Action string
	User   string
	ExitAt time.Time
}

func (s *AnchorExit) ExecuteAt() time.Time {
	return s.ExitAt
}

func (s *AnchorExit) Execute(sa *wavespread.Wavespread) (interface{}, error) {
	logrus.WithFields(logrus.Fields{}).Info("-> anchor exit")

	err := sa.ExitAnchor(s.User)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
