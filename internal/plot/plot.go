package plot

import (
	"fmt"
	"github.com/eldius/docker-profiler/internal/helper"
	"github.com/eldius/docker-profiler/internal/model"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"os"
	"time"
)

func PlotToFile(mdps []model.MetricsDatapoint) error {
	count := len(mdps)
	mups := make([]float64, count)
	mlps := make([]float64, count)
	cops := make([]float64, count)
	cups := make([]float64, count)
	cpps := make([]float64, count)
	xvalues := make([]time.Time, count)

	for i, m := range mdps {
		mups[i] = m.MemoryUsage
		mlps[i] = m.MemoryLimit
		cops[i] = m.CPUOnlineCount
		cups[i] = m.CPUUsage
		cpps[i] = m.CPUPercentage
		xvalues[i] = m.Timestamp

		fmt.Println("---")
		fmt.Printf("memory usage: %s\n", m.MemoryUsageStr())
		fmt.Printf("memory limit: %s\n", m.MemoryLimitStr())
		fmt.Printf("timestamp:    %v\n", m.Timestamp.Format(time.RFC3339))
	}
	defaults := chart.StyleTextDefaults()
	defaults.DotColor = drawing.ColorBlack
	defaults.StrokeColor = drawing.ColorGreen
	defaults.FontColor = drawing.ColorBlue
	defaults.DotWidth = 1
	defaults.TextWrap = chart.TextWrapWord
	defaults.FontSize = 18
	defaults.Font, _ = chart.GetDefaultFont()

	memoryFormatter := func(v interface{}) string {
		if t, ok := v.(float64); ok {
			return helper.FormatMemory(uint64(t))
		}
		return fmt.Sprintf("%v", v)
	}
	plotGraph("memory_usage.svg", chart.TimeSeries{
		Name:    "Memory Usage",
		XValues: xvalues,
		YValues: mups,
		Style:   defaults,
	}, memoryFormatter)

	plotGraph("memory_limit.svg", chart.TimeSeries{
		Name:    "Memory Limit",
		XValues: xvalues,
		YValues: mlps,
		Style:   defaults,
	}, memoryFormatter)

	plotGraph("cpu_online_count.svg", chart.TimeSeries{
		Name:    "CP Online Count",
		XValues: xvalues,
		YValues: cops,
		Style:   defaults,
	}, func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	})

	plotGraph("cpu_usage.svg", chart.TimeSeries{
		Name:    "CPU Usage",
		XValues: xvalues,
		YValues: cups,
		Style:   defaults,
	}, func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	})

	plotGraph("cpu_percentage.svg", chart.TimeSeries{
		Name:    "CPU Percentage",
		XValues: xvalues,
		YValues: cpps,
		Style:   defaults,
	}, func(v interface{}) string {
		return fmt.Sprintf("%v%%", v)
	})

	return nil
}

func plotGraph(file string, c chart.Series, formatter chart.ValueFormatter) error {
	style := c.GetStyle()

	gridLineStyle := c.GetStyle()
	gridLineStyle.StrokeColor = drawing.ColorBlack
	gridLineStyle.DotColor = drawing.ColorRed

	graph := chart.Chart{
		Title:        c.GetName(),
		TitleStyle:   style,
		ColorPalette: nil,
		Width:        1024,
		Height:       1024,
		Background:   style,
		Canvas:       style,
		Font:         style.Font,
		XAxis: chart.XAxis{
			Name:  "Time",
			Style: style,
			GridLines: []chart.GridLine{{
				IsMinor: false,
				Style:   gridLineStyle,
			}},
		},
		YAxis: chart.YAxis{
			Name:  "Memory",
			Style: style,
			GridLines: []chart.GridLine{{
				IsMinor: false,
				Style:   gridLineStyle,
			}},
			ValueFormatter: formatter,
		},
		Series: []chart.Series{
			c,
		},
	}
	f, _ := os.Create(file)
	defer func() {
		_ = f.Close()
	}()
	return graph.Render(chart.SVG, f)
}
