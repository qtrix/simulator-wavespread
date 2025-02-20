package charter

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Params struct {
	StartTime   time.Time
	EndTime     time.Time
	NrAnchors   int64
	NrSurfers   int64
	ChartPoints int64
	EpochLength int64
}

func getParams(r *http.Request) (*Params, error) {
	start := r.URL.Query().Get("start")
	if start == "" {
		return nil, errors.New("missing start")
	}

	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse start time")
	}

	end := r.URL.Query().Get("end")
	if end == "" {
		return nil, errors.New("missing end")
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse end time")
	}

	anchorsInt := int64(10)
	nrAnchors := r.URL.Query().Get("anchors")
	if nrAnchors != "" {
		anchorsInt, err = strconv.ParseInt(nrAnchors, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse 'anchors'")
		}
	}

	surfersInt := int64(10)
	nrSurfers := r.URL.Query().Get("surfers")
	if nrSurfers != "" {
		surfersInt, err = strconv.ParseInt(nrSurfers, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse 'surfers'")
		}
	}

	points := r.URL.Query().Get("points")
	pointsInt, err := strconv.ParseInt(points, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse 'points'")
	}

	epochLength := r.URL.Query().Get("epoch")
	if epochLength == "" {
		epochLength = "604800"
	}
	epochLengthInt, err := strconv.ParseInt(epochLength, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse 'epoch'")
	}

	return &Params{
		StartTime:   startTime,
		EndTime:     endTime,
		NrAnchors:   anchorsInt,
		NrSurfers:   surfersInt,
		ChartPoints: pointsInt,
		EpochLength: epochLengthInt,
	}, nil
}
