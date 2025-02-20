package actions

import (
	"time"

	"github.com/qtrix/simulator-wavespread/wavespread"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type SurferPosition struct {
	Action  string
	User    string
	Amount  decimal.Decimal
	EntryAt time.Time
}

func (p *SurferPosition) ExecuteAt() time.Time {
	return p.EntryAt
}

func (p *SurferPosition) Execute(sa *wavespread.Wavespread) (interface{}, error) {
	logrus.WithFields(logrus.Fields{
		"amount": p.Amount,
	}).Info("-> surfer deposit")

	return nil, sa.DepositSurfer(p.User, p.Amount)
}
