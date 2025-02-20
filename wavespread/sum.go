package wavespread

// func (sa *SmartAlpha) JuniorLiquiditySumIndividual() (decimal.Decimal, error) {
// 	paid, err := sa.AbondPaidSumIndividual()
// 	if err != nil {
// 		return decimal.Zero, err
// 	}
//
// 	juniorProfits, err := sa.JuniorProfitsSumIndividual()
// 	if err != nil {
// 		return decimal.Zero, err
// 	}
//
// 	return sa.TotalBalance.Sub(sa.Abond.Principal).Sub(paid).Add(juniorProfits), nil
// }

// func (sa *SmartAlpha) AbondPaidSumIndividual() (decimal.Decimal, error) {
// 	sum := decimal.Zero
// 	for _, b := range sa.seniors {
// 		if b.Redeemed {
// 			continue
// 		}
//
// 		bp, err := sa.SeniorBondDownside(b)
// 		if err != nil {
// 			return decimal.Zero, err
// 		}
//
// 		sum = sum.Add(bp)
// 	}
//
// 	return sum, nil
// }

// func (sa *SmartAlpha) JuniorProfitsSumIndividual() (decimal.Decimal, error) {
// 	sum := decimal.Zero
// 	for _, b := range sa.seniors {
// 		if b.Redeemed {
// 			continue
// 		}
//
// 		bp, err := sa.SeniorBondUpside(b)
// 		if err != nil {
// 			return decimal.Zero, err
// 		}
//
// 		sum = sum.Add(bp)
// 	}
//
// 	return sum, nil
// }
//
// func (sa *SmartAlpha) JtokenPriceSUM() (decimal.Decimal, error) {
// 	// junior liquidity / junior liquidity - junior token supply
// 	if sa.JTokenSupply.Equal(decimal.NewFromInt(0)) {
// 		return decimal.NewFromInt(1), nil
// 	}
//
// 	juniorLiq, err := sa.JuniorLiquiditySumIndividual()
// 	if err != nil {
// 		return decimal.Decimal{}, err
// 	}
//
// 	return juniorLiq.Div(sa.JTokenSupply), nil
// }
