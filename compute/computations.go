package compute

import (
	"bulk2023_jeen/datatransfer"

	"github.com/gonum/matrix/mat64"
)

func CalculateTemperature(beamData *datatransfer.BeamData) map[float64][]float64 {
	dx := beamData.Length / float64(beamData.Partitions)
	dt := 1.0

	Nx := int(beamData.Partitions)

	temperatureMap := make(map[float64][]float64)

	initialValue := beamData.TemperatureData[0]

	temperatureMap[0] = initialValue

	r1 := make([]float64, Nx)
	r2 := make([]float64, Nx)

	for i := 1; i < Nx; i++ {
		r1[i] = beamData.Conductivity[i-1] + beamData.Conductivity[i]
	}

	for i := 0; i < Nx-1; i++ {
		r2[i] = beamData.Conductivity[i] + beamData.Conductivity[i+1]
	}

	A := make([]float64, Nx)
	B := make([]float64, Nx)
	C := make([]float64, Nx)

	for i := 0; i < Nx; i++ {
		A[i] = (dt * r1[i]) / (2 * dx * dx)
		B[i] = 1 + (dt*(r1[i]+r2[i]))/(2*dx*dx)
		C[i] = (dt * r2[i]) / (2 * dx * dx)
	}

	coeffMatrix := mat64.NewDense(Nx, Nx, nil)

	for i := 0; i < Nx; i++ {
		coeffMatrix.Set(i, i, B[i])
		if i > 0 {
			coeffMatrix.Set(i, i-1, -A[i])
		}
		if i < Nx-1 {
			coeffMatrix.Set(i, i+1, -C[i])
		}
	}

	coeffMatrix.Set(0, 0, 1)
	coeffMatrix.Set(Nx-1, Nx-1, 1)
	coeffMatrix.Set(0, 1, 0)
	coeffMatrix.Set(Nx-1, Nx-2, 0)

	for i := 1; i < int(beamData.TimeMoments); i++ {
		prevData := temperatureMap[float64(i-1)]
		t := mat64.NewDense(Nx, 1, prevData)

		var D mat64.Dense

		D.Solve(coeffMatrix, t)
		temperatureMap[float64(i)] = D.ColView(0).RawVector().Data
	}

	return temperatureMap
}
