package smartalpha

import "github.com/shopspring/decimal"

func (sa *SmartAlpha) JTokenPrice() (decimal.Decimal, error) {
	// junior liquidity / junior liquidity - junior token supply
	if sa.JTokenSupply.Equal(decimal.NewFromInt(0)) {
		return decimal.NewFromInt(1), nil
	}

	juniorLiq, err := sa.JuniorLiquidity()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return juniorLiq.Div(sa.JTokenSupply), nil
}

func (sa *SmartAlpha) STokenPrice() (decimal.Decimal, error) {
	// senior liquidity / senior token supply
	if sa.STokenSupply.Equal(decimal.NewFromInt(0)) {
		return decimal.NewFromInt(1), nil
	}

	seniorLiq, err := sa.SeniorLiquidity()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return seniorLiq.Div(sa.STokenSupply), nil
}
