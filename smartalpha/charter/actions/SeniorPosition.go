package actions

import (
	"time"

	"github.com/qtrix/simulator-wavespread/smartalpha"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type SeniorPosition struct {
	Action  string
	User    string
	Amount  decimal.Decimal
	EntryAt time.Time
}

func (s *SeniorPosition) ExecuteAt() time.Time {
	return s.EntryAt
}

func (s *SeniorPosition) Execute(sa *smartalpha.SmartAlpha) (interface{}, error) {
	entryPrice, _ := sa.GetCurrentPrice()

	logrus.WithFields(logrus.Fields{
		"amount":     s.Amount,
		"entryPrice": entryPrice,
	}).Info("-> senior deposit")

	err := sa.DepositSenior(s.User, s.Amount)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
