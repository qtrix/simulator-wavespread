package smartalpha

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (sa *SmartAlpha) MinPrice() decimal.Decimal {
	return sa.EntryPrice.Mul(decimal.NewFromInt(1).Sub(sa.DownsideProtectionRate))
}

func (sa *SmartAlpha) TotalLoss() (decimal.Decimal, error) {
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

	return sa.TotalSeniors.Mul(x).Round(18), nil
}

func (sa *SmartAlpha) TotalProfits() (decimal.Decimal, error) {
	currentPrice, err := sa.GetCurrentPrice()
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "could not get current price")
	}

	calcPrice := decimal.Min(currentPrice, sa.EntryPrice)

	x := OneDec.Sub(calcPrice.Div(currentPrice))
	y := OneDec.Sub(sa.UpsideExposureRate)

	return x.Mul(y).Mul(sa.TotalSeniors).Round(18), nil
}

func (sa *SmartAlpha) MaxProtectionAmount() decimal.Decimal {
	if sa.MinPrice().Equal(decimal.Zero) {
		return decimal.Zero
	}

	return sa.TotalSeniors.Mul(sa.DownsideProtectionRate)
}
