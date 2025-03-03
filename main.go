package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type record struct {
	relativeTime int64 // Relative Time(sec)
	memory       int64 // Memory(RSS Bytes)
}

func main() {
	var data1File, data2File string

	var rootCmd = &cobra.Command{
		Use:   "plot",
		Short: "Generate a line chart from one or two CSV files",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				data1, data2 []record
			)
			data1 = readCSV(data1File)
			if data2File != "" {
				data2 = readCSV(data2File)
			}

			// Create the plot
			p := plot.New()

			p.Title.Text = "Memory usage over time"
			p.X.Label.Text = "Relative Time (Seconds)"
			p.Y.Label.Text = "Memory (RSS KB)"
			p.Y.Tick.Marker = memoryTicks{}

			// generate line1
			line1, err := generateLineFromData(data1)
			if err != nil {
				log.Fatalf("Error creating line for data1: %v", err)
			}
			line1.LineStyle.Color = plotutil.Color(0)

			// Add the line to the plot
			p.Add(line1)
			p.Legend.Add(filepath.Base(data1File), line1)

			// If a second file was provided, add the second line plot
			if len(data2) > 0 {
				// generate line2
				line2, err := generateLineFromData(data2)
				if err != nil {
					log.Fatalf("Error creating line for data2: %v", err)
				}
				line2.LineStyle.Color = plotutil.Color(1)
				p.Add(line2)
				p.Legend.Add(filepath.Base(data2File), line2)
			}

			// Save the plot to a PNG file
			if err := p.Save(10*vg.Inch, 6*vg.Inch, "line_chart.png"); err != nil {
				log.Fatalf("Error saving plot: %v", err)
			}

			fmt.Println("Chart saved to 'line_chart.png'")
		},
	}

	// Define flags for the file inputs
	rootCmd.Flags().StringVarP(&data1File, "data1", "", "", "Path to the first CSV file (required)")
	rootCmd.Flags().StringVarP(&data2File, "data2", "", "", "Path to the second CSV file (optional)")

	// Mark flags as required
	rootCmd.MarkFlagRequired("data1")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func generateLineFromData(recs []record) (*plotter.Line, error) {
	values := make(plotter.XYs, len(recs))
	for i, rec := range recs {
		values[i].X = float64(rec.relativeTime)
		values[i].Y = float64(rec.memory)
	}

	line, err := plotter.NewLine(values)
	if err != nil {
		return nil, err
	}
	return line, nil
}

func readCSV(fileName string) []record {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening CSV file %s: %v", fileName, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var data []record

	// Skip the header
	_, err = reader.Read()

	for {
		rec, err := reader.Read()
		if err != nil {
			break
		}

		relativeTime, err := strconv.ParseInt(strings.TrimSpace(rec[1]), 10, 0)
		if err != nil {
			fmt.Printf("Invalid ratative time: %s, error: %v\n", rec[1], err)
			continue
		}

		memory, err := strconv.ParseInt(strings.TrimSpace(rec[len(rec)-1]), 10, 0)
		if err != nil {
			fmt.Printf("Invalid memory: %s, error: %v\n", rec[len(rec)-1], err)
			continue
		}

		data = append(data, record{relativeTime: relativeTime, memory: memory * 1024})
	}

	return data
}

func formatMemory(bytes float64) string {
	const (
		KiB = 1024
		MiB = KiB * 1024
		GiB = MiB * 1024
	)

	switch {
	case bytes >= GiB:
		return fmt.Sprintf("%.2fGiB", bytes/GiB)
	case bytes >= MiB:
		return fmt.Sprintf("%.2fMiB", bytes/MiB)
	case bytes >= KiB:
		return fmt.Sprintf("%.2fKiB", bytes/KiB)
	default:
		return fmt.Sprintf("%.0fB", bytes)
	}
}

type memoryTicks struct{}

func (memoryTicks) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	step := (max - min) / 10
	for y := min; y <= max; y += step {
		ticks = append(ticks, plot.Tick{Value: y, Label: formatMemory(y)})
	}
	return ticks
}
