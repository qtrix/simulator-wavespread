package wavespread

import (
	"github.com/shopspring/decimal"
)

func (sa *Wavespread) SurferLiquidity() (decimal.Decimal, error) {
	// total deposits - total seniors - total loss + total profits
	totalLoss, err := sa.TotalLoss()
	if err != nil {
		return decimal.Decimal{}, err
	}

	totalProfits, err := sa.TotalProfits()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return sa.TotalBalance.Sub(sa.TotalAnchors).Sub(totalLoss).Add(totalProfits), nil
}

func (sa *Wavespread) AnchorLiquidity() (decimal.Decimal, error) {
	totalLoss, err := sa.TotalLoss()
	if err != nil {
		return decimal.Decimal{}, err
	}

	totalProfits, err := sa.TotalProfits()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return sa.TotalAnchors.Sub(totalProfits).Add(totalLoss), nil
}

func (sa *Wavespread) SurferLoanable() (decimal.Decimal, error) {
	return sa.TotalBalance.Sub(sa.TotalAnchors).Sub(sa.MaxProtectionAmount()), nil
}
