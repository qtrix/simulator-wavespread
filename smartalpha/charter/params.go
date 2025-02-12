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
	NrSeniors   int64
	NrJuniors   int64
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

	seniorsInt := int64(10)
	nrSeniors := r.URL.Query().Get("seniors")
	if nrSeniors != "" {
		seniorsInt, err = strconv.ParseInt(nrSeniors, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse 'seniors'")
		}
	}

	juniorsInt := int64(10)
	nrJuniors := r.URL.Query().Get("juniors")
	if nrJuniors != "" {
		juniorsInt, err = strconv.ParseInt(nrJuniors, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse 'juniors'")
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
		NrSeniors:   seniorsInt,
		NrJuniors:   juniorsInt,
		ChartPoints: pointsInt,
		EpochLength: epochLengthInt,
	}, nil
}
