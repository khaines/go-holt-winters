package holtwinters

// Creates an array of forcasted data points based on the given time series
func TripleExponentialSmoothing(series []float64, alpha, beta, gamma float64, seasonLength, nPredictions int) []float64 {
	// initialize the predicted series
	predictedValues := make([]float64, 0, len(series)+nPredictions)

	previousLevel := series[0]
	currentTrend := initialTrend(series, seasonLength)
	currentSeason := initialSeasonality(series, seasonLength)

	seriesLength := len(series)
	for i := 0; i < seriesLength+nPredictions; i++ {

		if i >= seriesLength {
			// calculating forecasted predictions using the generated model
			m := i - seriesLength + 1
			predictedValues = append(predictedValues, (previousLevel+float64(m)*currentTrend)+currentSeason[i%seasonLength])
		} else {
			// calculating predicted values based on existing data
			currentLevel := alpha*(series[i]-currentSeason[i%seasonLength]) + (1-alpha)*(previousLevel+currentTrend)
			currentTrend = beta*(currentLevel-previousLevel) + (1-beta)*currentTrend
			previousLevel = currentLevel
			currentSeason[i%seasonLength] = gamma*(series[i]-currentLevel) + (1-gamma)*currentSeason[i%seasonLength]
			predictedValues = append(predictedValues, currentLevel+currentTrend+currentSeason[i%seasonLength])
		}

	}

	return predictedValues
}

// Calculates the initial trend for the holt/winters forcast
func initialTrend(series []float64, seasonLength int) float64 {
	// initialize the trend
	trend := 0.0

	// loop through the time series
	for i := 0; i < seasonLength; i++ {
		trend += (series[i+seasonLength] - series[i]) / float64(seasonLength)
	}

	return trend / float64(seasonLength)
}

// Calculates the initial season for the holt/winters forcast
func initialSeasonality(series []float64, seasonLength int) []float64 {
	initialSeason := make([]float64, 0, seasonLength)
	numSeasons := int(len(series) / seasonLength)

	// calculates the average for each season in the data
	seasonAvgs := make([]float64, numSeasons)
	for i := 0; i < numSeasons; i++ {
		for j := i * seasonLength; j < (i*seasonLength)+seasonLength; j++ {
			seasonAvgs[i] += series[j]
		}
		seasonAvgs[i] /= float64(seasonLength)
	}

	for i := 0; i < seasonLength; i++ {
		seasonVector := 0.0
		for j := 0; j < numSeasons; j++ {
			seasonVector += series[seasonLength*j+i] - seasonAvgs[j]
		}
		initialSeason = append(initialSeason, seasonVector/float64(numSeasons))
	}

	return initialSeason
}
