package smartalpha

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"

	"github.com/qtrix/simulator-wavespread/smartalpha/interfaces"
)

type SmartAlpha struct {
	oracle interfaces.IOracle
	Clock  interfaces.IClock

	seniors []*SeniorPosition
	juniors []*JuniorPosition

	EntryPrice             decimal.Decimal
	UpsideExposureRate     decimal.Decimal
	DownsideProtectionRate decimal.Decimal

	JTokenSupply decimal.Decimal
	STokenSupply decimal.Decimal

	TotalBalance decimal.Decimal
	TotalSeniors decimal.Decimal

	// need to figure out how to do entry and exit without loops
	EntryQueueJuniors []*JuniorPosition
	EntryQueueSeniors []*SeniorPosition

	ExitQueueJuniors []*JuniorPosition
	ExitQueueSeniors []*SeniorPosition
}

func New(oracle interfaces.IOracle, clock interfaces.IClock) *SmartAlpha {
	return &SmartAlpha{
		oracle: oracle,
		Clock:  clock,
	}
}

func (sa *SmartAlpha) StartNextEpoch() error {
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
	sa.TotalSeniors = sa.TotalSeniors.Add(loss).Sub(profit)

	logrus.WithFields(logrus.Fields{
		"ts":     sa.Clock.Now(),
		"loss":   loss,
		"profit": profit,
	}).Info("starting next epoch")

	// 3. set the new entry price => profits and losses will be reset to 0
	sa.EntryPrice = ep

	// 3. process the entry and exit queues

	jTokenPrice, err := sa.JTokenPrice()
	if err != nil {
		return errors.Wrap(err, "could not get jToken price")
	}

	sTokenPrice, err := sa.STokenPrice()
	if err != nil {
		return errors.Wrap(err, "could not get sToken price")
	}

	// logrus.WithField("count", len(sa.EntryQueueJuniors)).Info("processing junior entry queue")
	for _, j := range sa.EntryQueueJuniors {
		sa.TotalBalance = sa.TotalBalance.Add(j.Amount)

		// calculate the amount to mint for the user
		jTokenAmount := j.Amount.Div(jTokenPrice)

		// mint the jTokens
		sa.JTokenSupply = sa.JTokenSupply.Add(jTokenAmount)

		// update the junior position
		j.EntryPrice = ep
		j.JTokenBalance = jTokenAmount

		sa.juniors = append(sa.juniors, j)
	}

	// logrus.WithField("count", len(sa.ExitQueueJuniors)).Info("processing junior exit queue")
	for _, j := range sa.ExitQueueJuniors {
		amount := j.JTokenBalance.Mul(jTokenPrice)

		sa.JTokenSupply = sa.JTokenSupply.Sub(j.JTokenBalance)
		sa.TotalBalance = sa.TotalBalance.Sub(amount)
	}

	// logrus.WithField("count", len(sa.EntryQueueSeniors)).Info("processing senior entry queue")
	for _, s := range sa.EntryQueueSeniors {
		sa.TotalSeniors = sa.TotalSeniors.Add(s.Amount)
		sa.TotalBalance = sa.TotalBalance.Add(s.Amount)

		sTokenAmount := s.Amount.Div(sTokenPrice)

		sa.STokenSupply = sa.STokenSupply.Add(sTokenAmount)

		s.STokenBalance = sTokenAmount
		s.EntryPrice = ep

		sa.seniors = append(sa.seniors, s)
	}

	// logrus.WithField("count", len(sa.ExitQueueSeniors)).Info("processing senior exit queue")
	for _, s := range sa.ExitQueueSeniors {
		amount := s.STokenBalance.Mul(sTokenPrice)

		sa.STokenSupply = sa.STokenSupply.Sub(s.STokenBalance)
		sa.TotalSeniors = sa.TotalSeniors.Sub(amount)
		sa.TotalBalance = sa.TotalBalance.Sub(amount)
	}

	// empty the queues
	sa.EntryQueueJuniors = make([]*JuniorPosition, 0)
	sa.EntryQueueSeniors = make([]*SeniorPosition, 0)
	sa.ExitQueueJuniors = make([]*JuniorPosition, 0)
	sa.ExitQueueSeniors = make([]*SeniorPosition, 0)

	// update the upside and downside rates for the new pool composition
	sa.UpsideExposureRate, sa.DownsideProtectionRate = sa.CalcSeniorRates()

	logrus.WithField("up", sa.UpsideExposureRate).WithField("down", sa.DownsideProtectionRate).Info("got new rates")

	return nil
}

// DepositJunior
// junior deposit into next epoch
// - we don't know yet the jToken price at which the underlying will be converted
//   - the conversion will happen when the new epoch starts
//
// - junior must be added to a queue
//   - in the real thing, we'll probably issue a NFT which will be redeemable after the epoch starts
func (sa *SmartAlpha) DepositJunior(user string, amount decimal.Decimal) error {
	ts := sa.Clock.Now().Unix()

	p := &JuniorPosition{
		Owner:     user,
		Amount:    amount,
		EntryTime: ts,
	}

	sa.EntryQueueJuniors = append(sa.EntryQueueJuniors, p)

	return nil
}

func (sa *SmartAlpha) DepositSenior(user string, amount decimal.Decimal) error {
	ts := sa.Clock.Now().Unix()

	b := &SeniorPosition{
		Owner:     user,
		Amount:    amount,
		EntryTime: ts,
	}

	sa.EntryQueueSeniors = append(sa.EntryQueueSeniors, b)

	return nil
}

func (sa *SmartAlpha) ExitSenior(user string) error {
	var p *SeniorPosition

	for _, v := range sa.seniors {
		if v.Owner == user {
			p = v
			break
		}
	}

	if p == nil {
		return errors.New("user not found")
	}

	if p.Exited {
		return errors.New("already exited")
	}

	p.Exited = true
	sa.ExitQueueSeniors = append(sa.ExitQueueSeniors, p)

	return nil
}

func (sa *SmartAlpha) ExitJunior(user string) error {
	var p *JuniorPosition

	for _, v := range sa.juniors {
		if v.Owner == user {
			p = v
			break
		}
	}

	if p == nil {
		return errors.New("user not found")
	}

	if p.Exited {
		return errors.New("already exited")
	}

	p.Exited = true
	sa.ExitQueueJuniors = append(sa.ExitQueueJuniors, p)

	return nil
}

func (sa *SmartAlpha) CalcSeniorRates() (decimal.Decimal, decimal.Decimal) {
	juniorLiq, _ := sa.JuniorLiquidity()
	seniorLiq, _ := sa.SeniorLiquidity()

	if juniorLiq.Add(seniorLiq).Equal(decimal.Zero) {
		return OneDec, decimal.Zero
	}

	juniorDominance := juniorLiq.Div(juniorLiq.Add(seniorLiq))

	protection := decimal.NewFromFloat(0.8).Mul(juniorDominance)
	if protection.GreaterThan(decimal.NewFromFloat(0.35)) {
		protection = decimal.NewFromFloat(0.35)
	}

	var sum decimal.Decimal

	if juniorDominance.LessThan(decimal.NewFromFloat(0.05)) {
		sum = decimal.NewFromInt(-18).Mul(juniorDominance).Add(OneDec)
	} else {
		sum = decimal.NewFromInt(18).Div(decimal.NewFromInt(19)).Mul(juniorDominance).Add(OneDec.Div(decimal.NewFromInt(19)))
	}

	return sum.Sub(protection), protection
}

func (sa *SmartAlpha) CalcSeniorRate() (decimal.Decimal, error) {
	// (junior_loanable_liquidity / total_pool_liquidity)
	loanable, err := sa.JuniorLiquidity()
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
