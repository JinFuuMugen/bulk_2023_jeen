package compute

import (
	"bulk2023_jeen/datatransfer"
)

func solveTDMA(a, b, c, d []float64) []float64 {
	n := len(a)

	ac := make([]float64, n)
	bc := make([]float64, n)
	cc := make([]float64, n)
	dc := make([]float64, n)
	copy(ac, a)
	copy(bc, b)
	copy(cc, c)
	copy(dc, d)

	cc[0] /= bc[0]
	dc[0] /= bc[0]
	n--
	for i := 1; i < n; i++ {
		cc[i] /= bc[i] - ac[i]*cc[i-1]
		dc[i] = (dc[i] - ac[i]*dc[i-1]) / (bc[i] - ac[i]*cc[i-1])
	}

	dc[n] = (dc[n] - ac[n]*dc[n-1]) / (bc[n] - ac[n]*cc[n-1])

	for i := n - 1; i >= 0; i-- {
		dc[i] -= cc[i] * dc[i+1]
	}

	return dc
}

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
		A[i] = -(dt * r1[i]) / (2 * dx * dx)
		B[i] = 1 + (dt*(r1[i]+r2[i]))/(2*dx*dx)
		C[i] = -(dt * r2[i]) / (2 * dx * dx)
	}

	B[0] = 1
	C[0] = 0
	B[len(B)-1] = 1
	A[len(A)-1] = 0
	A[0] = 0
	C[len(C)-1] = 0

	for i := 1; i < int(beamData.TimeMoments); i++ {
		temperatureMap[float64(i)] = solveTDMA(A, B, C, temperatureMap[float64(i-1)])
	}

	return temperatureMap
}
