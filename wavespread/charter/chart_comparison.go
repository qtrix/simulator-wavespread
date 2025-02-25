package charter

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (c *Charter) comparisonChart(w http.ResponseWriter, r *http.Request) {
	params, err := getParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	spew.Dump("--------------xxxxxxxx---------------")
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

	//var xAxis []string
	//var price, anchorValue, surferValue, anchorValueRegular, surferValueRegular []opts.LineData
	//for _, point := range chartPoints {
	//	xAxis = append(xAxis, point.Timestamp.String())
	//
	//	price = append(price, opts.LineData{Value: point.Price})
	//	anchorValue = append(anchorValue, opts.LineData{Value: point.AnchorValue})
	//	anchorValueRegular = append(anchorValueRegular, opts.LineData{Value: point.AnchorValueRegular})
	//	surferValue = append(surferValue, opts.LineData{Value: point.SurferValue})
	//	surferValueRegular = append(surferValueRegular, opts.LineData{Value: point.SurferValueRegular})
	//}
	//
	//// create a new line instance
	//line := charts.NewLine()
	//// set some global options like Title/Legend/ToolTip or anything else
	//line.SetGlobalOptions(
	//	charts.WithInitializationOpts(opts.Initialization{
	//		Theme:  types.ThemeWesteros,
	//		Width:  "100%",
	//		Height: "800px",
	//	}),
	//	charts.WithTitleOpts(opts.Title{
	//		Title:    "2SA || !2SA",
	//		Subtitle: fmt.Sprintf("showing data from %s to %s", params.StartTime, params.EndTime),
	//	}),
	//	charts.WithLegendOpts(opts.Legend{
	//		Show:   true,
	//		Orient: "horizontal",
	//		Top:    "0",
	//		Right:  "0",
	//	}),
	//	charts.WithTooltipOpts(opts.Tooltip{
	//		Show: true,
	//	}),
	//	charts.WithDataZoomOpts(opts.DataZoom{
	//		Type:       "inside",
	//		Start:      0,
	//		End:        100,
	//		XAxisIndex: []int{0},
	//	}),
	//	charts.WithDataZoomOpts(opts.DataZoom{
	//		Type:       "slider",
	//		Start:      0,
	//		End:        100,
	//		XAxisIndex: []int{0},
	//	}),
	//)
	//
	//line.ExtendYAxis(opts.YAxis{})
	//
	//line.AddSeries("Anchor value", anchorValue).
	//	AddSeries("Surfer value", surferValue).
	//	AddSeries("Anchor value (without SA)", anchorValueRegular).
	//	AddSeries("Surfer value (without SA)", surferValueRegular)
	//
	//line.SetXAxis(xAxis).
	//	AddSeries("Price", price,
	//		charts.WithLineStyleOpts(opts.LineStyle{Color: "red", Width: 2, Type: "dotted"}),
	//		charts.WithLineChartOpts(opts.LineChart{YAxisIndex: 1}),
	//	)
	//
	//line.Render(w)

	data, _ := json.MarshalIndent(map[string]interface{}{
		"data":    chartPoints,
		"actions": actions,
	}, "", "  ")
	spew.Dump("--------------------------------------------")
	spew.Dump(data)
	w.Write(data)
}
