package wavespread

import "github.com/shopspring/decimal"

func (sa *Wavespread) SurferTokenPrice() (decimal.Decimal, error) {
	// surfer liquidity / surfer token supply
	if sa.SurferTokenSupply.Equal(decimal.NewFromInt(0)) {
		return decimal.NewFromInt(1), nil
	}

	surferLiq, err := sa.SurferLiquidity()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return surferLiq.Div(sa.SurferTokenSupply), nil
}

func (sa *Wavespread) AnchorTokenPrice() (decimal.Decimal, error) {
	// anchor liquidity / anchor token supply
	if sa.AnchorTokenSupply.Equal(decimal.NewFromInt(0)) {
		return decimal.NewFromInt(1), nil
	}

	anchorLiq, err := sa.AnchorLiquidity()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return anchorLiq.Div(sa.AnchorTokenSupply), nil
}
