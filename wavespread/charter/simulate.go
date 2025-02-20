package charter

import (
	"time"

	"github.com/pkg/errors"

	"github.com/qtrix/simulator-wavespread/wavespread"
	"github.com/qtrix/simulator-wavespread/wavespread/charter/actions"
	"github.com/qtrix/simulator-wavespread/wavespread/mocks"
)

func minTime(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t1
	}

	return t2
}

func (c *Charter) simulate(params *Params) ([]Point, actions.Actions, error) {
	actionsList := c.generateActions(params)

	var allActions actions.Actions
	for _, a := range actionsList {
		allActions = append(allActions, a)
	}

	duration := params.EndTime.Sub(params.StartTime)
	distance := time.Duration(duration.Nanoseconds() / params.ChartPoints)

	clock := &mocks.Clock{}
	clock.SetTime(params.StartTime)
	sa := wavespread.New(c.Oracle, clock)

	epochLength := time.Duration(params.EpochLength) * time.Second

	nextPointTs := params.StartTime
	nextTaskTs := minTime(nextPointTs, actionsList[0].ExecuteAt())

	var points []Point

	lastTs := actionsList[len(actionsList)-1].ExecuteAt()

	for t := params.StartTime; t.Before(lastTs.Add(epochLength * 3)); t = t.Add(epochLength) {
		clock.SetTime(t)

		point, err := c.getPoint(sa)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not generate point")
		}

		points = append(points, *point)

		err = sa.StartNextEpoch()
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not start next epoch")
		}

		// execute all actions that are scheduled for this epoch
		// calculate all the points matching this epoch

		for {
			if nextTaskTs.After(t.Add(epochLength)) {
				break
			}

			if len(actionsList) > 0 && (actionsList[0].ExecuteAt().Before(nextPointTs) || actionsList[0].ExecuteAt().Equal(nextPointTs)) {
				action := actionsList[0]
				actionsList = actionsList[1:]

				clock.SetTime(action.ExecuteAt())
				_, err := action.Execute(sa)
				if err != nil {
					if err.Error() == "blocked" {
						continue
					}
					return nil, nil, errors.Wrap(err, "could not execute action")
				}
			} else {
				clock.SetTime(nextPointTs)

				point, err := c.getPoint(sa)
				if err != nil {
					return nil, nil, errors.Wrap(err, "could not generate point")
				}

				points = append(points, *point)
				nextPointTs = nextPointTs.Add(distance)
			}

			if len(actionsList) > 0 {
				nextTaskTs = minTime(actionsList[0].ExecuteAt(), nextPointTs)
			} else {
				nextTaskTs = nextPointTs
			}
		}
	}

	return points, allActions, nil
}
