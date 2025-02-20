package actions

import (
	"time"

	"github.com/qtrix/simulator-wavespread/wavespread"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type AnchorPosition struct {
	Action  string
	User    string
	Amount  decimal.Decimal
	EntryAt time.Time
}

func (s *AnchorPosition) ExecuteAt() time.Time {
	return s.EntryAt
}

func (s *AnchorPosition) Execute(sa *wavespread.Wavespread) (interface{}, error) {
	entryPrice, _ := sa.GetCurrentPrice()

	logrus.WithFields(logrus.Fields{
		"amount":     s.Amount,
		"entryPrice": entryPrice,
	}).Info("-> anchor deposit")

	err := sa.DepositAnchor(s.User, s.Amount)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
