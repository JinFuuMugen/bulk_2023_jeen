package main

import (
	"bulk2023_jeen/compute"
	"bulk2023_jeen/datatransfer"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	logFile, err := os.Create("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	if err := os.Truncate("log.txt", 0); err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)

	startTime := time.Now()

	if len(os.Args) != 2 {
		log.Println("Usage: go run main.go <input_file_path>")
		os.Exit(1)
	}

	var beamData datatransfer.BeamData

	datatransfer.ReadParseInput(os.Args[1], &beamData)

	log.Printf("GRID: length: %.2f, partitions: %d, moments: %d\n\n\n", beamData.Length, beamData.Partitions, beamData.TimeMoments)

	log.Printf("TUBE: conductivity: %.7f\n\n\n", beamData.Conductivity)

	newVal := compute.CalculateTemperature(&beamData)

	var discrepancies float64

	for i := 0; i < int(beamData.TimeMoments); i++ {
		for j := 0; j < int(beamData.Partitions); j++ {
			if beamData.TemperatureData[float64(i)][j] != -1.0 {
				discrepancies += math.Pow(newVal[float64(i)][j]-beamData.TemperatureData[float64(i)][j], 2)
			}

		}
	}

	log.Printf("Calc values are: %.3f\n\n\n", newVal)

	log.Printf("Given values are: %.3f\n\n\n", beamData.TemperatureData)

	log.Printf("Discrepancies: %.7f\n\n\n", discrepancies)

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)

	log.Printf("Time elapsed: %s\n\n\n", elapsed)
}
