package wavespread

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"

	"github.com/qtrix/simulator-wavespread/wavespread/interfaces"
)

type Wavespread struct {
	oracle interfaces.IOracle
	Clock  interfaces.IClock

	anchors []*AnchorPosition
	surfers []*SurferPosition

	EntryPrice             decimal.Decimal
	UpsideExposureRate     decimal.Decimal
	DownsideProtectionRate decimal.Decimal

	AnchorTokenSupply decimal.Decimal //anchors
	SurferTokenSupply decimal.Decimal //surfers

	TotalBalance decimal.Decimal
	TotalAnchors decimal.Decimal

	// need to figure out how to do entry and exit without loops
	EntryQueueSurfers []*SurferPosition
	EntryQueueAnchors []*AnchorPosition

	ExitQueueSurfers []*SurferPosition
	ExitQueueAnchors []*AnchorPosition

	SwitchSideQueueAnchorToSurfer []*AnchorPosition
	SwitchSideQueueSurferToAnchor []*SurferPosition
}

func New(oracle interfaces.IOracle, clock interfaces.IClock) *Wavespread {
	return &Wavespread{
		oracle: oracle,
		Clock:  clock,
	}
}

func (sa *Wavespread) StartNextEpoch() error {
	ts := sa.Clock.Now().Unix()

	ep, err := sa.oracle.PriceAtTs(ts)
	if err != nil {
		return err
	}

	// 1. finalize the profits and losses of the previous epoch
	loss, err := sa.TotalLoss()
	if err != nil {
		return errors.Wrap(err, "could not get total loss")
	}

	profit, err := sa.TotalProfits()
	if err != nil {
		return errors.Wrap(err, "could not get total profits")
	}

	// 2. move the funds from one side to the other
	sa.TotalAnchors = sa.TotalAnchors.Add(loss).Sub(profit)

	logrus.WithFields(logrus.Fields{
		"ts":     sa.Clock.Now(),
		"loss":   loss,
		"profit": profit,
	}).Info("starting next epoch")

	// 3. set the new entry price => profits and losses will be reset to 0
	sa.EntryPrice = ep

	// 3. process the entry and exit queues

	surferTokenPrice, err := sa.SurferTokenPrice()
	if err != nil {
		return errors.Wrap(err, "could not get jToken price")
	}

	anchorTokenPrice, err := sa.AnchorTokenPrice()
	if err != nil {
		return errors.Wrap(err, "could not get sToken price")
	}

	logrus.WithField("count", len(sa.EntryQueueSurfers)).Info("processing surfers entry queue")
	for _, j := range sa.EntryQueueSurfers {
		sa.TotalBalance = sa.TotalBalance.Add(j.Amount)

		// calculate the amount to mint for the user
		STokenAmount := j.Amount.Div(surferTokenPrice)

		// mint the jTokens
		sa.SurferTokenSupply = sa.SurferTokenSupply.Add(STokenAmount)

		// update the junior position
		j.EntryPrice = ep
		j.STokenBalance = STokenAmount

		sa.surfers = append(sa.surfers, j)
	}

	logrus.WithField("count", len(sa.ExitQueueSurfers)).Info("processing surfers exit queue")
	for _, j := range sa.ExitQueueSurfers {
		amount := j.STokenBalance.Mul(surferTokenPrice)

		sa.SurferTokenSupply = sa.SurferTokenSupply.Sub(j.STokenBalance)
		sa.TotalBalance = sa.TotalBalance.Sub(amount)
	}

	logrus.WithField("count", len(sa.EntryQueueAnchors)).Info("processing anchors entry queue")
	for _, s := range sa.EntryQueueAnchors {
		sa.TotalAnchors = sa.TotalAnchors.Add(s.Amount)
		sa.TotalBalance = sa.TotalBalance.Add(s.Amount)

		aTokenAmount := s.Amount.Div(anchorTokenPrice)

		sa.AnchorTokenSupply = sa.AnchorTokenSupply.Add(aTokenAmount)

		s.ATokenBalance = aTokenAmount
		s.EntryPrice = ep

		sa.anchors = append(sa.anchors, s)
	}

	logrus.WithField("count", len(sa.ExitQueueAnchors)).Info("processing anchors exit queue")
	for _, s := range sa.ExitQueueAnchors {
		amount := s.ATokenBalance.Mul(anchorTokenPrice)

		sa.AnchorTokenSupply = sa.AnchorTokenSupply.Sub(s.ATokenBalance)
		sa.TotalAnchors = sa.TotalAnchors.Sub(amount)
		sa.TotalBalance = sa.TotalBalance.Sub(amount)
	}

	/// @dev User will join the switch queue and his anchor tokens will be transferred back to the pool.
	/// @dev Their tokens will be burned when the epoch is finalized and the underlying due will be moved to the surfer side of the pool.
	/// @dev Users can increase their queue amount but can't exit the queue
	/// @param amountAnchorTokens The amount of tokens the user wants to move to surfer side
	logrus.WithField("count", len(sa.SwitchSideQueueAnchorToSurfer)).Info("processing anchor to surfer switch queue")
	spew.Dump(sa.SwitchSideQueueAnchorToSurfer)
	spew.Dump(sa.AnchorTokenSupply)
	spew.Dump(sa.TotalAnchors)

	for _, s := range sa.SwitchSideQueueAnchorToSurfer {
		// calculate the amount in underlying tokens for the user
		//amount := s.ATokenBalance.Mul(anchorTokenPrice)
		//calculate the amount of surfer tokens to mint for the user

		//this would be burned
		aTokenAmount := s.Amount.Div(anchorTokenPrice)
		sa.AnchorTokenSupply = sa.AnchorTokenSupply.Sub(aTokenAmount)
		sa.TotalAnchors = sa.TotalAnchors.Sub(s.Amount)

		//now we mint new surfers tokens based on amount
		STokenAmount := s.Amount.Div(surferTokenPrice)
		// mint the sTokens
		sa.SurferTokenSupply = sa.SurferTokenSupply.Add(STokenAmount)

		a := &SurferPosition{
			Owner:         s.Owner,
			Amount:        s.Amount,
			EntryTime:     ts,
			EntryPrice:    ep,
			STokenBalance: STokenAmount,
		}

		sa.surfers = append(sa.surfers, a)

		//pop the anchor from the switch queue
		var newAnchors []*AnchorPosition
		for _, anchor := range sa.anchors {
			if anchor.Owner != s.Owner {
				newAnchors = append(newAnchors, anchor)
			}
		}
		sa.anchors = newAnchors
	}

	logrus.WithField("count", len(sa.SwitchSideQueueSurferToAnchor)).Info("processing surfer to anchor switch queue")
	for _, s := range sa.SwitchSideQueueSurferToAnchor {
		// Calculate the amount to burn for the user
		//this would be burned
		sTokenAmount := s.Amount.Div(surferTokenPrice)
		sa.SurferTokenSupply = sa.SurferTokenSupply.Sub(sTokenAmount)

		//now we mint new surfers tokens based on amount
		ATokenAmount := s.Amount.Div(anchorTokenPrice)
		sa.AnchorTokenSupply = sa.AnchorTokenSupply.Add(ATokenAmount)
		sa.TotalAnchors = sa.TotalAnchors.Add(s.Amount)

		// Create a new anchor position
		a := &AnchorPosition{
			Owner:         s.Owner,
			Amount:        s.Amount,
			EntryTime:     ts,
			EntryPrice:    ep,
			ATokenBalance: ATokenAmount,
		}

		// Add the new anchor position
		sa.anchors = append(sa.anchors, a)

		// Pop the surfer from the switch queue
		var newSurfers []*SurferPosition
		for _, surfer := range sa.surfers {
			if surfer.Owner != s.Owner {
				newSurfers = append(newSurfers, surfer)
			}
		}
		sa.surfers = newSurfers
	}

	// empty the queues
	sa.EntryQueueSurfers = make([]*SurferPosition, 0)
	sa.EntryQueueAnchors = make([]*AnchorPosition, 0)
	sa.ExitQueueSurfers = make([]*SurferPosition, 0)
	sa.ExitQueueAnchors = make([]*AnchorPosition, 0)
	sa.SwitchSideQueueSurferToAnchor = make([]*SurferPosition, 0)
	sa.SwitchSideQueueAnchorToSurfer = make([]*AnchorPosition, 0)

	// update the upside and downside rates for the new pool composition
	sa.UpsideExposureRate, sa.DownsideProtectionRate = sa.CalcAnchorRates()

	logrus.WithField("up", sa.UpsideExposureRate).WithField("down", sa.DownsideProtectionRate).Info("got new rates")

	return nil
}

// DepositSurfer
// junior deposit into next epoch
// - we don't know yet the jToken price at which the underlying will be converted
//   - the conversion will happen when the new epoch starts
//
// - junior must be added to a queue
//   - in the real thing, we'll probably issue a NFT which will be redeemable after the epoch starts
func (sa *Wavespread) DepositSurfer(user string, amount decimal.Decimal) error {
	ts := sa.Clock.Now().Unix()

	p := &SurferPosition{
		Owner:     user,
		Amount:    amount,
		EntryTime: ts,
	}

	sa.EntryQueueSurfers = append(sa.EntryQueueSurfers, p)

	return nil
}

func (sa *Wavespread) DepositAnchor(user string, amount decimal.Decimal) error {
	ts := sa.Clock.Now().Unix()

	b := &AnchorPosition{
		Owner:     user,
		Amount:    amount,
		EntryTime: ts,
	}

	sa.EntryQueueAnchors = append(sa.EntryQueueAnchors, b)

	return nil
}

func (sa *Wavespread) ExitAnchor(user string) error {
	var p *AnchorPosition
	for _, v := range sa.anchors {
		if v.Owner == user {
			spew.Dump("----------------------", user, "----------------", v.Owner)
			p = v
			break
		}
	}

	if p == nil {
		return nil
	}

	if p.Exited {
		return errors.New("already exited")
	}

	p.Exited = true
	sa.ExitQueueAnchors = append(sa.ExitQueueAnchors, p)

	return nil
}

func (sa *Wavespread) ExitSurfer(user string) error {
	var p *SurferPosition
	for _, v := range sa.surfers {
		if v.Owner == user {
			spew.Dump("----------------------", user, "----------------", v.Owner)
			p = v
			break
		}
	}

	if p == nil {
		return nil
	}

	if p.Exited {
		return errors.New("already exited")
	}

	p.Exited = true
	sa.ExitQueueSurfers = append(sa.ExitQueueSurfers, p)

	return nil
}

func (sa *Wavespread) SwitchSideAnchorToSurfer(user string) error {
	ts := sa.Clock.Now().Unix()

	//first we check for the user in the anchor queue
	var p *AnchorPosition

	for _, v := range sa.anchors {
		if v.Owner == user {
			p = v
			break
		}
	}

	if p == nil {
		return errors.New("user not found for switch side anchor to surfer")
	}

	//check if the user didnt already exited
	if p.Exited {
		return errors.New("already exited")
	}

	//add the user to the switch queue
	sa.SwitchSideQueueAnchorToSurfer = append(sa.SwitchSideQueueAnchorToSurfer, &AnchorPosition{
		Owner:     user,
		Amount:    p.Amount,
		EntryTime: ts,
	})

	return nil
}

func (sa *Wavespread) SwitchSideSurferToAnchor(user string) error {
	ts := sa.Clock.Now().Unix()

	//first we check for the user in the surfer queue
	var p *SurferPosition

	for _, v := range sa.surfers {
		if v.Owner == user {
			p = v
			break
		}
	}

	if p == nil {
		return errors.New("user not found for switch side surfer to anchor")
	}

	//check if the user didnt already exited
	if p.Exited {
		return errors.New("already exited")
	}

	//add the user to the switch queue
	sa.SwitchSideQueueSurferToAnchor = append(sa.SwitchSideQueueSurferToAnchor, &SurferPosition{
		Owner:     user,
		Amount:    p.Amount,
		EntryTime: ts,
	})

	return nil
}

func (sa *Wavespread) CalcAnchorRates() (decimal.Decimal, decimal.Decimal) {
	surferLiq, _ := sa.SurferLiquidity()
	anchorLiq, _ := sa.AnchorLiquidity()

	if surferLiq.Add(anchorLiq).Equal(decimal.Zero) {
		return OneDec, decimal.Zero
	}

	surferDominance := surferLiq.Div(surferLiq.Add(anchorLiq))

	protection := decimal.NewFromFloat(0.8).Mul(surferDominance)
	if protection.GreaterThan(decimal.NewFromFloat(0.35)) {
		protection = decimal.NewFromFloat(0.35)
	}

	var sum decimal.Decimal

	if surferDominance.LessThan(decimal.NewFromFloat(0.05)) {
		sum = decimal.NewFromInt(-18).Mul(surferDominance).Add(OneDec)
	} else {
		sum = decimal.NewFromInt(18).Div(decimal.NewFromInt(19)).Mul(surferDominance).Add(OneDec.Div(decimal.NewFromInt(19)))
	}

	return sum.Sub(protection), protection
}

func (sa *Wavespread) CalcAnchorRate() (decimal.Decimal, error) {
	// (junior_loanable_liquidity / total_pool_liquidity)
	loanable, err := sa.SurferLiquidity()
	if err != nil {
		return decimal.Decimal{}, err
	}

	total := sa.TotalBalance
	if total.Equal(decimal.Zero) {
		return decimal.Zero, nil
	}

	rate := loanable.Div(total)

	return decimal.Min(rate, MaxRate), nil
}
