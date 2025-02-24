package wavespread

import (
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

	// If current price <= entry price, there are no surfer profits
	if currentPrice.LessThanOrEqual(sa.EntryPrice) {
		return decimal.Zero, nil
	}

	// Constants
	scaleFactor := decimal.NewFromInt(1e18) // Equivalent to SCALE_FACTOR in Solidity

	// (currentPrice - entryPrice)
	x := currentPrice.Sub(sa.EntryPrice)

	// (1 - upsideExposureRate)
	y := decimal.NewFromInt(1).Sub(sa.UpsideExposureRate)

	// (current price - entry price) * (1 - upsideExposureRate) * total anchors / current price / SCALE_FACTOR
	profits := x.Mul(y).Mul(sa.TotalAnchors).Div(currentPrice.Mul(scaleFactor))

	return profits.Round(18), nil
}

// anchor profit : total profits
func (sa *Wavespread) TotalProfits() (decimal.Decimal, error) {
	currentPrice, err := sa.GetCurrentPrice()
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "could not get current price")
	}

	// If price went up, there are no losses for surfers (i.e., no profits for anchors)
	if sa.EntryPrice.LessThanOrEqual(currentPrice) {
		return decimal.Zero, nil
	}

	// minPrice = (entryPrice * (1 - downsideProtectionRate)) + 1
	scaleFactor := decimal.NewFromInt(1e18) // Equivalent to SCALE_FACTOR in Solidity
	minPrice := sa.EntryPrice.Mul(scaleFactor.Sub(sa.DownsideProtectionRate)).Div(scaleFactor).Add(decimal.NewFromInt(1))

	// Ensure `calcPrice` is not lower than `minPrice`
	calcPrice := decimal.Max(currentPrice, minPrice)

	// (totalAnchors * entryPrice) / calcPrice - totalAnchors
	anchorProfits := sa.TotalAnchors.Mul(sa.EntryPrice).Div(calcPrice).Sub(sa.TotalAnchors)

	return anchorProfits.Round(18), nil
}

func (sa *Wavespread) MaxProtectionAmount() decimal.Decimal {
	if sa.MinPrice().Equal(decimal.Zero) {
		return decimal.Zero
	}

	return sa.TotalAnchors.Mul(sa.DownsideProtectionRate)
}
