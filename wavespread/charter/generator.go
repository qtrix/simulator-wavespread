package charter

import (
	"github.com/davecgh/go-spew/spew"
	"math/rand"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"github.com/thanhpk/randstr"

	"github.com/qtrix/simulator-wavespread/wavespread/charter/actions"
)

func (c *Charter) generateActions(params *Params) actions.Actions {
	surfers, anchors := c.generateRandomPositions(params.NrSurfers, params.NrAnchors, params.StartTime, params.EndTime)
	spew.Dump("surfers:", surfers)
	spew.Dump("anchors:", anchors)
	epochLengthTime := time.Duration(params.EpochLength) * time.Second

	var list actions.Actions
	for _, a := range surfers {
		//add surfer entry
		list = append(list, a)

		//if surfer-exit didnt happen, add surfer switch side anchor
		//need to make sure that t.ExitAt is after SwitchedAt
		if viper.GetBool("surfer-switch-side-anchor-enabled") {
			//little bit of a crap, but we need to make sure that we will have the user in anchor side after the switch
			var t actions.Action
			if viper.GetBool("surfer-exit-enabled") {
				t = &actions.AnchorExit{
					Action: "anchor-exit",
					User:   a.User,
					ExitAt: randomTimestamp(a.EntryAt.Add(epochLengthTime), params.EndTime),
				}

				list = append(list, t)
			}

			var switchedAt time.Time
			if t != nil {
				switchedAt = t.ExecuteAt()
			} else {
				switchedAt = params.EndTime
			}

			switchSide := &actions.SurferSwitchSideAnchor{
				Action:     "surfer-switch-side-anchor",
				User:       a.User,
				SwitchedAt: randomTimestamp(a.EntryAt.Add(epochLengthTime), switchedAt),
			}

			list = append(list, switchSide)
		} else {
			var t actions.Action
			if viper.GetBool("surfer-exit-enabled") {
				t = &actions.SurferExit{
					Action: "surfer-exit",
					User:   a.User,
					ExitAt: randomTimestamp(a.EntryAt.Add(epochLengthTime), params.EndTime),
				}

				list = append(list, t)
			}
		}
	}

	for _, a := range anchors {
		list = append(list, a)

		//if anchor-exit didnt happen, add anchor switch side surfer
		//need to make sure that t.ExitAt is after SwitchedAt
		if viper.GetBool("anchor-switch-side-surfer-enabled") {
			//little bit of a crap, but we need to make sure that we will have the user in surfer side after the switch
			var t actions.Action
			if viper.GetBool("anchor-exit-enabled") {
				t = &actions.SurferExit{
					Action: "surfer-exit",
					User:   a.User,
					ExitAt: randomTimestamp(a.EntryAt.Add(epochLengthTime), params.EndTime),
				}

				list = append(list, t)
			}

			var switchedAt time.Time
			if t != nil {
				switchedAt = t.ExecuteAt()
			} else {
				switchedAt = params.EndTime
			}

			switchSide := &actions.AnchorSwitchSideSurfer{
				Action:     "anchor-switch-side-surfer",
				User:       a.User,
				SwitchedAt: randomTimestamp(a.EntryAt.Add(epochLengthTime), switchedAt),
			}

			list = append(list, switchSide)
		} else {
			var t actions.Action
			if viper.GetBool("anchor-exit-enabled") {
				t = &actions.AnchorExit{
					Action: "anchor-exit",
					User:   a.User,
					ExitAt: randomTimestamp(a.EntryAt.Add(epochLengthTime), params.EndTime),
				}

				list = append(list, t)
			}
		}
	}

	sort.Sort(list)

	return list
}

func (c *Charter) generateRandomPositions(numSurfers, numAnchors int64, start, end time.Time) ([]*actions.SurferPosition, []*actions.AnchorPosition) {
	var surfers []*actions.SurferPosition
	var anchors []*actions.AnchorPosition

	for i := int64(0); i < numSurfers; i++ {
		j := generateRandomSurfer(start, end)
		if i == 0 {
			j.EntryAt = start
		}
		surfers = append(surfers, j)
	}

	for i := 0; int64(i) < numAnchors; i++ {
		s := generateRandomAnchor(start, end)
		if i == 0 {
			s.EntryAt = start.Add(1 * time.Second)
		}

		anchors = append(anchors, s)
	}

	return surfers, anchors
}

func generateRandomSurfer(start, end time.Time) *actions.SurferPosition {
	amount := decimal.NewFromInt(rand.Int63n(5) + 1)

	return &actions.SurferPosition{
		Action:  "surfer-entry",
		User:    randstr.Hex(20),
		Amount:  amount,
		EntryAt: randomTimestamp(start, end),
	}
}

func generateRandomAnchor(start, end time.Time) *actions.AnchorPosition {
	amount := decimal.NewFromInt(rand.Int63n(5) + 1)

	return &actions.AnchorPosition{
		Action:  "anchor-entry",
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
