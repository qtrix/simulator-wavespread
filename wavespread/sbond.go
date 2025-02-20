package wavespread

// func (sa *SmartAlpha) AbondMaxProtectionSum() decimal.Decimal {
// 	sum := decimal.Zero
// 	for _, b := range sa.seniors {
// 		if b.Redeemed {
// 			continue
// 		}
//
// 		bp := b.Amount.Mul(b.DownsideRate)
//
// 		sum = sum.Add(bp)
// 	}
//
// 	return sum
// }

// func (sa *SmartAlpha) SeniorBondDownside(b *SeniorPosition) (decimal.Decimal, error) {
// 	// calc_price = MAX(cur_price, min_price)
// 	// sBOND.paid() = ((sBOND.principal * sBOND.ep) / calc_price) - principal
// 	currentPrice, err := sa.GetCurrentPrice()
// 	if err != nil {
// 		return decimal.Zero, err
// 	}
//
// 	if currentPrice.GreaterThanOrEqual(b.EntryPrice) {
// 		return decimal.Zero, nil
// 	}
//
// 	calcPrice := decimal.Max(currentPrice, sa.SeniorMinPrice(b))
// 	// calcPrice := currentPrice
//
// 	return b.Amount.Mul(b.EntryPrice).Div(calcPrice).Sub(b.Amount), nil
// }
//
// func (sa *SmartAlpha) SeniorBondUpside(b *SeniorPosition) (decimal.Decimal, error) {
// 	// jProfits = (cur_price - entry_price) * (1 - urate) * principal / cur_price
// 	currentPrice, err := sa.GetCurrentPrice()
// 	if err != nil {
// 		return decimal.Zero, err
// 	}
//
// 	if currentPrice.LessThanOrEqual(b.EntryPrice) {
// 		return decimal.Zero, nil
// 	}
//
// 	return currentPrice.
// 		Sub(b.EntryPrice).
// 		Mul(
// 			decimal.NewFromInt(1).
// 				Sub(b.UpsideRate),
// 		).Mul(b.Amount).
// 		Div(currentPrice), nil
// }
//
// func (sa *SmartAlpha) SeniorBondById(id int64) *SeniorPosition {
// 	for _, v := range sa.seniors {
// 		if v.Id == id {
// 			return v
// 		}
// 	}
//
// 	return nil
// }
