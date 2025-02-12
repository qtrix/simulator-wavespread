package charter

import (
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/qtrix/simulator-wavespread/smartalpha"
)

type Point struct {
	Timestamp           time.Time
	TotalSeniors        decimal.Decimal
	SeniorLiq           decimal.Decimal
	SeniorValue         decimal.Decimal
	SeniorValueRegular  decimal.Decimal
	JuniorLiq           decimal.Decimal
	JuniorValue         decimal.Decimal
	JuniorValueRegular  decimal.Decimal
	TotalProfits        decimal.Decimal
	TotalLoss           decimal.Decimal
	Price               decimal.Decimal
	EntryPrice          decimal.Decimal
	MinPrice            decimal.Decimal
	JTokenPrice         decimal.Decimal
	STokenPrice         decimal.Decimal
	MaxProtectionAmount decimal.Decimal
}

func (c *Charter) getPoint(sa *smartalpha.SmartAlpha) (*Point, error) {
	seniorLiq, err := sa.SeniorLiquidity()
	if err != nil {
		return nil, errors.Wrap(err, "could not get senior liquidity")
	}

	juniorLiq, err := sa.JuniorLiquidity()
	if err != nil {
		return nil, errors.Wrap(err, "could not get junior liquidity")
	}

	totalProfits, err := sa.TotalProfits()
	if err != nil {
		return nil, errors.Wrap(err, "could not get junior profits")
	}

	totalLoss, err := sa.TotalLoss()
	if err != nil {
		return nil, errors.Wrap(err, "could not get abond paid")
	}

	price, err := sa.GetCurrentPrice()
	if err != nil {
		return nil, errors.Wrap(err, "could not get price")
	}

	jTokenPrice, err := sa.JTokenPrice()
	if err != nil {
		return nil, errors.Wrap(err, "could not get jtoken price")
	}

	sTokenPrice, err := sa.STokenPrice()
	if err != nil {
		return nil, errors.Wrap(err, "could not get sToken price")
	}

	return &Point{
		Timestamp:           sa.Clock.Now(),
		TotalSeniors:        sa.TotalSeniors.Round(2),
		SeniorLiq:           seniorLiq.Round(2),
		SeniorValue:         seniorLiq.Mul(price).Round(2),
		SeniorValueRegular:  sa.TotalSeniors.Mul(price).Round(2),
		JuniorLiq:           juniorLiq.Round(2),
		JuniorValue:         juniorLiq.Mul(price).Round(2),
		JuniorValueRegular:  sa.TotalBalance.Sub(sa.TotalSeniors).Mul(price).Round(2),
		TotalProfits:        totalProfits.Round(2),
		TotalLoss:           totalLoss.Round(2),
		Price:               price.Round(2),
		EntryPrice:          sa.EntryPrice.Round(2),
		MinPrice:            sa.MinPrice().Round(2),
		JTokenPrice:         jTokenPrice.Round(2),
		STokenPrice:         sTokenPrice.Round(2),
		MaxProtectionAmount: sa.MaxProtectionAmount().Round(2),
	}, nil
}
