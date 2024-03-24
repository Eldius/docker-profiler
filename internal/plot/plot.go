package plot

import (
	"fmt"
	"github.com/eldius/docker-profiler/internal/helper"
	"github.com/eldius/docker-profiler/internal/model"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"os"
	"path/filepath"
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
	f, _ := os.Create(filepath.Join(".data", file))
	defer func() {
		_ = f.Close()
	}()
	return graph.Render(chart.SVG, f)
}

func Plot(mdps []model.MetricsDatapoint) {
	count := len(mdps)
	memUsagePoints := make(plotter.XYs, count)
	memLimitPoints := make(plotter.XYs, count)
	memPercentage := make(plotter.XYs, count)
	cpuOnlinePoints := make(plotter.XYs, count)
	cpuUsagePoints := make(plotter.XYs, count)
	cpuPercentPoints := make(plotter.XYs, count)
	for i, v := range mdps {
		//memUsagePoints[i].X = float64(i) // Index as X value
		memUsagePoints[i].X = float64(v.Timestamp.Unix()) // Index as X value
		memUsagePoints[i].Y = v.MemoryUsage

		//memLimitPoints[i].X = float64(i)
		memLimitPoints[i].X = float64(v.Timestamp.Unix())
		memLimitPoints[i].Y = v.MemoryLimit

		memPercentage[i].X = float64(v.Timestamp.Unix())
		memPercentage[i].Y = helper.Percentage(uint64(v.MemoryUsage), uint64(v.MemoryLimit))

		//cpuOnlinePoints[i].X = float64(i)
		cpuOnlinePoints[i].X = float64(v.Timestamp.Unix())
		cpuOnlinePoints[i].Y = v.CPUOnlineCount

		//cpuPercentPoints[i].X = float64(i)
		cpuPercentPoints[i].X = float64(v.Timestamp.Unix())
		cpuPercentPoints[i].Y = v.CPUPercentage

		//cpuUsagePoints[i].X = float64(i)
		cpuUsagePoints[i].X = float64(v.Timestamp.Unix())
		cpuUsagePoints[i].Y = v.CPUUsage
	}

	memFormatter := newMemoryFormatter()
	percentageFormatter := newPercentageFormatter()

	draw(memUsagePoints, memFormatter, "Memory", "memory_usage.png", "Memory Usage")
	draw(memLimitPoints, memFormatter, "Memory", "memory_limit.png", "Memory Limit")
	draw(cpuUsagePoints, nil, "CPU Time", "cpu_usage.png", "CPU Usage")
	draw(cpuOnlinePoints, nil, "Number of CPUs", "cpu_online.png", "CPU Count")
	draw(cpuPercentPoints, percentageFormatter, "CPU Usage %", "cpu_percentage.png", "CPU Usage %")
}

func draw(data plotter.XYer, yFormatter plot.Ticker, yLabel, file, title string) {
	fmt.Printf("Printing chart '%s'...\n", title)

	xticks := plot.TimeTicks{Format: "2006-01-02\n15:04"}
	p := plot.New()
	p.Title.Text = title
	p.X.Tick.Marker = xticks
	if yFormatter != nil {
		p.Y.Tick.Marker = yFormatter
	}
	p.Y.Label.Text = yLabel
	p.Add(plotter.NewGrid())

	line, err := plotter.NewLine(data)
	if err != nil {
		panic(err)
	}
	p.Add(line)
	_ = plotutil.AddScatters(p, data)

	//dataCount := data.Len()
	if err := p.Save(30*vg.Inch, 10*vg.Inch, filepath.Join(".data", file)); err != nil {
		panic(err)
	}
}

func newMemoryFormatter() plot.Ticker {
	return &memoryTickerMarker{
		Ticker: plot.DefaultTicks{},
	}
}

func newPercentageFormatter() plot.Ticker {
	return &percentageTickerMarker{
		Ticker: plot.DefaultTicks{},
	}
}

type memoryTickerMarker struct {
	Ticker plot.Ticker
}

type percentageTickerMarker struct {
	Ticker plot.Ticker
}

func (m memoryTickerMarker) Ticks(min, max float64) []plot.Tick {
	ticks := m.Ticker.Ticks(min, max)
	for i := range ticks {
		tick := &ticks[i]
		if tick.Label == "" {
			continue
		}
		tick.Label = helper.FormatMemory(uint64(tick.Value))
	}
	return ticks
}

func (m percentageTickerMarker) Ticks(min, max float64) []plot.Tick {
	ticks := m.Ticker.Ticks(min, max)
	for i := range ticks {
		tick := &ticks[i]
		if tick.Label == "" {
			continue
		}
		tick.Label = fmt.Sprintf("%01.2f%%", tick.Value)
	}
	return ticks
}
