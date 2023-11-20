package main

import (
	"bulk2023_jeen/datatransfer"
	"fmt"
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

	fmt.Printf("GRID: length: %.10f, partitions: %d, moments: %d\n", beamData.Length, beamData.Partitions, beamData.TimeMoments)
	fmt.Printf("TUBE: conductivity: %.10f\n", beamData.Conductivity)
	fmt.Printf("TEMP: temperatureData: %.10f\n", beamData.TemperatureData)

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)

	fmt.Printf("Time elapsed: %s\n", elapsed)

}
