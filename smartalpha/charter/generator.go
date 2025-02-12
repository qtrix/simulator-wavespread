package charter

import (
	"math/rand"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"github.com/thanhpk/randstr"

	"github.com/qtrix/simulator-wavespread/smartalpha/charter/actions"
)

func (c *Charter) generateActions(params *Params) actions.Actions {
	juniors, seniors := c.generateRandomPositions(params.NrJuniors, params.NrSeniors, params.StartTime, params.EndTime)

	epochLengthTime := time.Duration(params.EpochLength) * time.Second

	var list actions.Actions
	for _, a := range juniors {
		list = append(list, a)

		if viper.GetBool("junior-exit-enabled") {
			exit := &actions.JuniorExit{
				Action: "junior-exit",
				User:   a.User,
				ExitAt: randomTimestamp(a.EntryAt.Add(epochLengthTime), params.EndTime),
			}

			list = append(list, exit)
		}
	}

	for _, a := range seniors {
		list = append(list, a)

		if viper.GetBool("senior-exit-enabled") {
			exit := &actions.SeniorExit{
				Action: "senior-exit",
				User:   a.User,
				ExitAt: randomTimestamp(a.EntryAt.Add(epochLengthTime), params.EndTime),
			}

			list = append(list, exit)
		}
	}

	sort.Sort(list)

	return list
}

func (c *Charter) generateRandomPositions(numJuniors, numSeniors int64, start, end time.Time) ([]*actions.JuniorPosition, []*actions.SeniorPosition) {
	var juniors []*actions.JuniorPosition
	var seniors []*actions.SeniorPosition

	for i := int64(0); i < numJuniors; i++ {
		j := generateRandomJunior(start, end)
		if i == 0 {
			j.EntryAt = start
		}
		juniors = append(juniors, j)
	}

	for i := 0; int64(i) < numSeniors; i++ {
		s := generateRandomSenior(start, end)
		if i == 0 {
			s.EntryAt = start.Add(1 * time.Second)
		}

		seniors = append(seniors, s)
	}

	return juniors, seniors
}

func generateRandomJunior(start, end time.Time) *actions.JuniorPosition {
	amount := decimal.NewFromInt(rand.Int63n(1000))

	return &actions.JuniorPosition{
		Action:  "junior-entry",
		User:    randstr.Hex(20),
		Amount:  amount,
		EntryAt: randomTimestamp(start, end),
	}
}

func generateRandomSenior(start, end time.Time) *actions.SeniorPosition {
	amount := decimal.NewFromInt(rand.Int63n(1000))

	return &actions.SeniorPosition{
		Action:  "senior-entry",
		User:    randstr.Hex(20),
		Amount:  amount,
		EntryAt: randomTimestamp(start, end),
	}
}

func randomTimestamp(start, end time.Time) time.Time {
	if end.Before(start) {
		return start.Add(24 * time.Hour)
	}

	randomTime := rand.Int63n(end.Unix()-start.Unix()) + start.Unix()

	randomNow := time.Unix(randomTime, 0)

	return randomNow
}
