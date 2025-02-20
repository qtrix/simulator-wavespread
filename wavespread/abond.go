package wavespread

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (sa *Wavespread) MinPrice() decimal.Decimal {
	return sa.EntryPrice.Mul(decimal.NewFromInt(1).Sub(sa.DownsideProtectionRate))
}

func (sa *Wavespread) TotalLoss() (decimal.Decimal, error) {
	currentPrice, err := sa.GetCurrentPrice()
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "could not get current price")
	}

	calcPrice := decimal.Max(
		decimal.Min(currentPrice, sa.EntryPrice),
		sa.MinPrice(),
	)

	if calcPrice.Equal(decimal.Zero) {
		return decimal.Zero, nil
	}

	x := sa.EntryPrice.Div(calcPrice).Sub(OneDec)

	return sa.TotalAnchors.Mul(x).Round(18), nil
}

func (sa *Wavespread) TotalProfits() (decimal.Decimal, error) {
	currentPrice, err := sa.GetCurrentPrice()
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "could not get current price")
	}

	calcPrice := decimal.Min(currentPrice, sa.EntryPrice)

	x := OneDec.Sub(calcPrice.Div(currentPrice))
	y := OneDec.Sub(sa.UpsideExposureRate)

	return x.Mul(y).Mul(sa.TotalAnchors).Round(18), nil
}

func (sa *Wavespread) MaxProtectionAmount() decimal.Decimal {
	if sa.MinPrice().Equal(decimal.Zero) {
		return decimal.Zero
	}

	return sa.TotalAnchors.Mul(sa.DownsideProtectionRate)
}
