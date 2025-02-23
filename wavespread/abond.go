package wavespread

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (sa *Wavespread) MinPrice() decimal.Decimal {
	return sa.EntryPrice.Mul(decimal.NewFromInt(1).Sub(sa.DownsideProtectionRate))
}

// anchor loss - surfer profits : total loss
func (sa *Wavespread) TotalLoss() (decimal.Decimal, error) {
	currentPrice, err := sa.GetCurrentPrice()
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "could not get current price")
	}
	if currentPrice.LessThanOrEqual(sa.EntryPrice) {
		return decimal.Zero, nil
	}

	x := currentPrice.Sub(sa.EntryPrice)
	y := OneDec.Sub(sa.DownsideProtectionRate)

	return x.Mul(y).Mul(sa.TotalAnchors).Div(currentPrice).Round(18), nil
}

// anchor profit : total profits
func (sa *Wavespread) TotalProfits() (decimal.Decimal, error) {
	currentPrice, err := sa.GetCurrentPrice()
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "could not get current price")
	}

	if sa.EntryPrice.LessThanOrEqual(currentPrice) {
		return decimal.Zero, nil
	}

	minPrice := (sa.EntryPrice.Mul(OneDec.Sub(sa.DownsideProtectionRate))).Round(18).Add(decimal.NewFromInt(1))

	spew.Dump(minPrice)
	if sa.EntryPrice.LessThanOrEqual(minPrice) {
		return decimal.Zero, nil
	}

	calcPrice := currentPrice
	if calcPrice.LessThan(minPrice) {
		calcPrice = minPrice
	}

	return (sa.TotalAnchors.Mul(sa.EntryPrice)).Div(calcPrice.Sub(sa.TotalAnchors)), nil

	//x := currentPrice.Sub(sa.EntryPrice)
	//y := OneDec.Sub(sa.UpsideExposureRate)
	//
	////x := OneDec.Sub(calcPrice.Div(currentPrice))
	////y := OneDec.Sub(sa.UpsideExposureRate)
	//
	//return x.Mul(y).Mul(sa.TotalAnchors).Div(currentPrice).Round(18), nil
}

func (sa *Wavespread) MaxProtectionAmount() decimal.Decimal {
	if sa.MinPrice().Equal(decimal.Zero) {
		return decimal.Zero
	}

	return sa.TotalAnchors.Mul(sa.DownsideProtectionRate)
}
