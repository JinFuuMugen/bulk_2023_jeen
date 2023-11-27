package main

import (
	"bulk2023_jeen/compute"
	"bulk2023_jeen/datatransfer"
	"fmt"
	"math"
	"os"
	"time"
)

func main() {

	startTime := time.Now()

	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_file_path>")
		os.Exit(1)
	}

	var beamData datatransfer.BeamData

	datatransfer.ReadParseInput(os.Args[1], &beamData)

	fmt.Printf("GRID: length: %.2f, partitions: %d, moments: %d\n", beamData.Length, beamData.Partitions, beamData.TimeMoments)
	fmt.Printf("TUBE: conductivity: %.7f\n", beamData.Conductivity)

	newVal := compute.CalculateTemperature(&beamData)

	var discrepancies float64

	for i := 0; i < int(beamData.TimeMoments); i++ {
		for j := 0; j < int(beamData.Partitions); j++ {
			discrepancies += math.Pow(newVal[float64(i)][j]-beamData.TemperatureData[float64(i)][j], 2)
		}
	}

	fmt.Printf("Calc values are: %.3f\n", newVal)
	fmt.Printf("Given values are: %.3f\n", beamData.TemperatureData)

	fmt.Printf("Discrepancies: %.7f\n", discrepancies)

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)

	fmt.Printf("Time elapsed: %s\n", elapsed)
}
