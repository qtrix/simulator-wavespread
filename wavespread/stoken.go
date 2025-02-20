package wavespread

import "github.com/shopspring/decimal"

func (sa *Wavespread) STokenPrice() (decimal.Decimal, error) {
	// junior liquidity / junior liquidity - junior token supply
	if sa.STokenSupply.Equal(decimal.NewFromInt(0)) {
		return decimal.NewFromInt(1), nil
	}

	juniorLiq, err := sa.SurferLiquidity()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return juniorLiq.Div(sa.STokenSupply), nil
}

func (sa *Wavespread) ATokenPrice() (decimal.Decimal, error) {
	// senior liquidity / senior token supply
	if sa.ATokenSupply.Equal(decimal.NewFromInt(0)) {
		return decimal.NewFromInt(1), nil
	}

	anchorLiq, err := sa.AnchorLiquidity()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return anchorLiq.Div(sa.ATokenSupply), nil
}
