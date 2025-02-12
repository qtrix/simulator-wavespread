package actions

import (
	"time"

	"github.com/qtrix/simulator-wavespread/smartalpha"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type JuniorPosition struct {
	Action  string
	User    string
	Amount  decimal.Decimal
	EntryAt time.Time
}

func (p *JuniorPosition) ExecuteAt() time.Time {
	return p.EntryAt
}

func (p *JuniorPosition) Execute(sa *smartalpha.SmartAlpha) (interface{}, error) {
	logrus.WithFields(logrus.Fields{
		"amount": p.Amount,
	}).Info("-> junior deposit")

	return nil, sa.DepositJunior(p.User, p.Amount)
}
