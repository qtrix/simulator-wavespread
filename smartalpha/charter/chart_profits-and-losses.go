package charter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/sirupsen/logrus"
)

func (c *Charter) profitsAndLosses(w http.ResponseWriter, r *http.Request) {
	params, err := getParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"start":  params.StartTime,
		"end":    params.EndTime,
		"points": params.ChartPoints,
	}).Info("Profits and losses Chart requested")

	chartPoints, actions, err := c.simulate(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var xAxis []string
	var losses, juniorProfits, price, abondEntryPrice []opts.LineData
	for _, point := range chartPoints {
		xAxis = append(xAxis, point.Timestamp.String())

		losses = append(losses, opts.LineData{Value: point.TotalLoss})
		juniorProfits = append(juniorProfits, opts.LineData{Value: point.TotalProfits})
		price = append(price, opts.LineData{Value: point.Price})
		abondEntryPrice = append(abondEntryPrice, opts.LineData{Value: point.EntryPrice})
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
			Title:    "Profits and losses",
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
		AddSeries("Junior Losses", losses).
		AddSeries("Junior profits", juniorProfits)

	line.AddSeries("Price", price,
		charts.WithLineStyleOpts(opts.LineStyle{Color: "red", Width: 2, Type: "dotted"}),
		charts.WithLineChartOpts(opts.LineChart{YAxisIndex: 1}),
	).AddSeries("Abond EntryPrice", abondEntryPrice,
		charts.WithLineStyleOpts(opts.LineStyle{Color: "green", Width: 2, Type: "dotted"}),
		charts.WithLineChartOpts(opts.LineChart{YAxisIndex: 1}),
	)

	line.Render(w)

	data, _ := json.MarshalIndent(map[string]interface{}{
		"actions": actions,
	}, "", "  ")
	w.Write([]byte("<pre>" + string(data) + "</pre>"))
}
