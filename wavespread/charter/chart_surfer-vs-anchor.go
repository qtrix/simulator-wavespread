package charter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (c *Charter) surferVsAnchor(w http.ResponseWriter, r *http.Request) {
	params, err := getParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"start":  params.StartTime,
		"end":    params.EndTime,
		"points": params.ChartPoints,
	}).Info("SurferVsAnchors Chart requested")

	chartPoints, actions, err := c.simulate(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var xAxis []string
	var anchorLiq, surferLiq, price []opts.LineData
	for _, point := range chartPoints {
		xAxis = append(xAxis, point.Timestamp.String())

		anchorPercent := decimal.NewFromFloat(0)
		surferPercent := decimal.NewFromFloat(0)
		sum := point.AnchorLiq.Add(point.SurferLiq)
		if sum.GreaterThan(decimal.Zero) {
			anchorPercent = point.AnchorLiq.Mul(decimal.NewFromInt(100)).Div(sum)
			surferPercent = point.SurferLiq.Mul(decimal.NewFromInt(100)).Div(sum)
		}

		anchorLiq = append(anchorLiq, opts.LineData{Value: anchorPercent})
		surferLiq = append(surferLiq, opts.LineData{Value: surferPercent})
		price = append(price, opts.LineData{Value: point.Price})
	}

	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  types.ThemeWesteros,
			Width:  "100%",
			Height: "800px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Surfer vs Anchor",
			Subtitle: fmt.Sprintf("showing data from %s to %s", params.StartTime, params.EndTime),
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:   true,
			Orient: "horizontal",
			Top:    "0",
			Right:  "0",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	line.ExtendYAxis(opts.YAxis{})

	// Put data into instance
	line.SetXAxis(xAxis).
		AddSeries("Anchor liquidity", anchorLiq, charts.WithAreaStyleOpts(opts.AreaStyle{Opacity: 0.5})).
		AddSeries("Surfer liquidity", surferLiq, charts.WithAreaStyleOpts(opts.AreaStyle{Opacity: 0.5})).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
			charts.WithLineChartOpts(opts.LineChart{Stack: "stack"}),
		)

	line.AddSeries("Price", price,
		charts.WithLineStyleOpts(opts.LineStyle{Color: "red", Width: 2, Type: "dotted"}),
		charts.WithLineChartOpts(opts.LineChart{YAxisIndex: 1}),
	)

	line.Render(w)

	data, _ := json.MarshalIndent(map[string]interface{}{
		"data":    chartPoints,
		"actions": actions,
	}, "", "  ")
	w.Write([]byte(string(data)))
}
