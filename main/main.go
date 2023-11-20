package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type BeamData struct {
	Length          float64
	Partitions      int64
	TimeMoments     int64
	TemperatureData map[float64][]float64
	Conductivity    []float64
}

const GRIDBLOCK string = "GRID"
const TEMPBLOCK string = "TEMP"
const TUBEBLOCK string = "TUBE"

func main() {

	startTime := time.Now()

	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_file_path>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	const maxBufferSize = 1 << 20
	buf := make([]byte, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)

	var beamData BeamData
	var currentBlock string

	for scanner.Scan() {

		line := scanner.Text()

		parts := strings.Split(line, "--")

		if len(parts) > 1 {
			line = parts[0]
		}

		if line == "/" {
			currentBlock = ""
			continue
		}

		flag := false

		if strings.Contains(line, " /") {
			flag = true
			line = strings.Split(line, " /")[0]
		}

		if line == GRIDBLOCK || line == TEMPBLOCK || line == TUBEBLOCK {
			currentBlock = line
		} else {
			switch currentBlock {
			case GRIDBLOCK:
				values := strings.Split(line, " ")
				beamData.Length, err = strconv.ParseFloat(values[0], 64)
				if err != nil {
					panic(fmt.Errorf("error parsing GRID value (length of beam): %w", err))
				}
				beamData.Partitions, err = strconv.ParseInt(values[1], 10, 64)
				if err != nil {
					panic(fmt.Errorf("error parsing GRID value (number of Partitions): %w", err))
				}
				beamData.TimeMoments, err = strconv.ParseInt(values[2], 10, 64)
				if err != nil {
					panic(fmt.Errorf("error parsing GRID value (number of TimeMoments): %w", err))
				}

			case TEMPBLOCK:

				beamData.TemperatureData = make(map[float64][]float64, beamData.TimeMoments)

				values := strings.Split(line, " ")

				timeIDX := 0
				for timeIDX < len(values) {
					tempData := make([]float64, 0)
					noDataFlag := false
					for i := timeIDX + 1; i <= timeIDX+int(beamData.Length); i++ {
						if strings.Contains(values[i], "-") {
							noDataFlag = true
							break
						} else {
							c, err := strconv.ParseFloat(values[i], 64)
							if err != nil {
								panic(fmt.Errorf("error parsing TEMP value (TemperatureData at %d index): %w", i, err))
							}
							tempData = append(tempData, c)
						}
					}

					var timeValue float64
					timeValue, err := strconv.ParseFloat(values[timeIDX], 64)
					if err != nil {
						panic(fmt.Errorf("error parsing TEMP value (TemperatureData at %d index): %w", timeIDX, err))
					}

					beamData.TemperatureData[timeValue] = tempData
					if !noDataFlag {
						timeIDX += int(beamData.Length) + 1
					} else {
						timeIDX += 2
					}
				}

			case TUBEBLOCK:
				values := strings.Split(line, " ")
				beamData.Conductivity = make([]float64, int64(beamData.Length))
				for i, v := range values {
					c, err := strconv.ParseFloat(v, 64)
					if err != nil {
						panic(fmt.Errorf("error parsing TUBE value (Conductivity at %d index): %w", i, err))
					}
					beamData.Conductivity[i] = c
				}
			}

			if flag {
				currentBlock = ""
			}
		}
	}

	fmt.Printf("GRID: length: %.10f, partitions: %d, moments: %d\n", beamData.Length, beamData.Partitions, beamData.TimeMoments)
	fmt.Printf("TUBE: conductivity: %.10f\n", beamData.Conductivity)
	fmt.Printf("TEMP: temperatureData: %.10f\n", beamData.TemperatureData)

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)

	fmt.Printf("Time elapsed: %s\n", elapsed)

}
