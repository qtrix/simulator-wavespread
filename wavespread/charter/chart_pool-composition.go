package charter

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (c *Charter) poolCompositionChart(w http.ResponseWriter, r *http.Request) {
	params, err := getParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.WithFields(logrus.Fields{
		"start":  params.StartTime,
		"end":    params.EndTime,
		"points": params.ChartPoints,
	}).Info("Pool composition chart requested")

	chartPoints, actions, err := c.simulate(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//spew.Dump(chartPoints)
	//var xAxis []string
	//var anchorLiq, surferLiq, surferProfits, abondPaid, price []opts.LineData
	//for _, point := range chartPoints {
	//	xAxis = append(xAxis, point.Timestamp.String())
	//
	//	anchorLiq = append(anchorLiq, opts.LineData{Value: point.AnchorLiq})
	//	surferLiq = append(surferLiq, opts.LineData{Value: point.SurferLiq})
	//	surferProfits = append(surferProfits, opts.LineData{Value: point.TotalProfits})
	//	abondPaid = append(abondPaid, opts.LineData{Value: point.TotalLoss})
	//	price = append(price, opts.LineData{Value: point.Price})
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
	//		Title:    "Pool composition",
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
	//// Put data into instance
	//line.SetXAxis(xAxis).
	//	AddSeries("Anchor liquidity", anchorLiq).
	//	AddSeries("Surfer liquidity", surferLiq).
	//	AddSeries("Surfer profits", surferProfits).
	//	AddSeries("Abond paid", abondPaid)
	//
	//line.AddSeries("Price", price,
	//	charts.WithLineStyleOpts(opts.LineStyle{Color: "red", Width: 2, Type: "dotted"}),
	//	charts.WithLineChartOpts(opts.LineChart{YAxisIndex: 1}),
	//)
	//
	//line.Render(w)

	data, _ := json.MarshalIndent(map[string]interface{}{
		"data":    chartPoints,
		"actions": actions,
	}, "", "  ")
	//spew.Dump(actions)
	w.Write(data)
}
