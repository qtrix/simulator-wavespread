package smartalpha

import (
	"github.com/shopspring/decimal"
)

func (sa *SmartAlpha) JuniorLiquidity() (decimal.Decimal, error) {
	// total deposits - total seniors - total loss + total profits
	totalLoss, err := sa.TotalLoss()
	if err != nil {
		return decimal.Decimal{}, err
	}

	totalProfits, err := sa.TotalProfits()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return sa.TotalBalance.Sub(sa.TotalSeniors).Sub(totalLoss).Add(totalProfits), nil
}

func (sa *SmartAlpha) SeniorLiquidity() (decimal.Decimal, error) {
	totalLoss, err := sa.TotalLoss()
	if err != nil {
		return decimal.Decimal{}, err
	}

	totalProfits, err := sa.TotalProfits()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return sa.TotalSeniors.Sub(totalProfits).Add(totalLoss), nil
}

func (sa *SmartAlpha) JuniorLoanable() (decimal.Decimal, error) {
	return sa.TotalBalance.Sub(sa.TotalSeniors).Sub(sa.MaxProtectionAmount()), nil
}
