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

func (c *Charter) comparisonChart(w http.ResponseWriter, r *http.Request) {
	params, err := getParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"start":  params.StartTime,
		"end":    params.EndTime,
		"points": params.ChartPoints,
	}).Info("Comparison chart requested")

	chartPoints, actions, err := c.simulate(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var xAxis []string
	var price, seniorValue, juniorValue, seniorValueRegular, juniorValueRegular []opts.LineData
	for _, point := range chartPoints {
		xAxis = append(xAxis, point.Timestamp.String())

		price = append(price, opts.LineData{Value: point.Price})
		seniorValue = append(seniorValue, opts.LineData{Value: point.SeniorValue})
		seniorValueRegular = append(seniorValueRegular, opts.LineData{Value: point.SeniorValueRegular})
		juniorValue = append(juniorValue, opts.LineData{Value: point.JuniorValue})
		juniorValueRegular = append(juniorValueRegular, opts.LineData{Value: point.JuniorValueRegular})
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
			Title:    "2SA || !2SA",
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

	line.AddSeries("Senior value", seniorValue).
		AddSeries("Junior value", juniorValue).
		AddSeries("Senior value (without SA)", seniorValueRegular).
		AddSeries("Junior value (without SA)", juniorValueRegular)

	line.SetXAxis(xAxis).
		AddSeries("Price", price,
			charts.WithLineStyleOpts(opts.LineStyle{Color: "red", Width: 2, Type: "dotted"}),
			charts.WithLineChartOpts(opts.LineChart{YAxisIndex: 1}),
		)

	line.Render(w)

	data, _ := json.MarshalIndent(map[string]interface{}{
		"actions": actions,
	}, "", "  ")
	w.Write([]byte("<pre>" + string(data) + "</pre"))
}
