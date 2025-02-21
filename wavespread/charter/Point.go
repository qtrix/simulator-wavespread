package charter

import (
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/qtrix/simulator-wavespread/wavespread"
)

type Point struct {
	Timestamp           time.Time
	TotalAnchors        decimal.Decimal
	AnchorLiq           decimal.Decimal
	AnchorValue         decimal.Decimal
	AnchorValueRegular  decimal.Decimal
	SurferLiq           decimal.Decimal
	SurferValue         decimal.Decimal
	SurferValueRegular  decimal.Decimal
	TotalProfits        decimal.Decimal
	TotalLoss           decimal.Decimal
	Price               decimal.Decimal
	EntryPrice          decimal.Decimal
	MinPrice            decimal.Decimal
	SurferTokenPrice    decimal.Decimal
	AnchorTokenPrice    decimal.Decimal
	MaxProtectionAmount decimal.Decimal
}

func (c *Charter) getPoint(sa *wavespread.Wavespread) (*Point, error) {
	anchorLiq, err := sa.AnchorLiquidity()
	if err != nil {
		return nil, errors.Wrap(err, "could not get anchor liquidity")
	}

	surferLiq, err := sa.SurferLiquidity()
	if err != nil {
		return nil, errors.Wrap(err, "could not get surfer liquidity")
	}

	totalProfits, err := sa.TotalProfits()
	if err != nil {
		return nil, errors.Wrap(err, "could not get surfer profits")
	}

	totalLoss, err := sa.TotalLoss()
	if err != nil {
		return nil, errors.Wrap(err, "could not get abond paid")
	}

	price, err := sa.GetCurrentPrice()
	if err != nil {
		return nil, errors.Wrap(err, "could not get price")
	}

	surferTokenPrice, err := sa.SurferTokenPrice()
	if err != nil {
		return nil, errors.Wrap(err, "could not get surfer token price")
	}

	anchorTokenPrice, err := sa.AnchorTokenPrice()
	if err != nil {
		return nil, errors.Wrap(err, "could not get anchor Token price")
	}

	return &Point{
		Timestamp:           sa.Clock.Now(),
		TotalAnchors:        sa.TotalAnchors.Round(2),
		AnchorLiq:           anchorLiq.Round(2),
		AnchorValue:         anchorLiq.Mul(price).Round(2),
		AnchorValueRegular:  sa.TotalAnchors.Mul(price).Round(2),
		SurferLiq:           surferLiq.Round(2),
		SurferValue:         surferLiq.Mul(price).Round(2),
		SurferValueRegular:  sa.TotalBalance.Sub(sa.TotalAnchors).Mul(price).Round(2),
		TotalProfits:        totalProfits.Round(2),
		TotalLoss:           totalLoss.Round(2),
		Price:               price.Round(2),
		EntryPrice:          sa.EntryPrice.Round(2),
		MinPrice:            sa.MinPrice().Round(2),
		SurferTokenPrice:    surferTokenPrice.Round(2),
		AnchorTokenPrice:    anchorTokenPrice.Round(2),
		MaxProtectionAmount: sa.MaxProtectionAmount().Round(2),
	}, nil
}
