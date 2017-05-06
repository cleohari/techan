package talib4g

type smaIndicator struct {
	indicator Indicator
	window    int
}

func NewSimpleMovingAverage(indicator Indicator, window int) Indicator {
	return smaIndicator{indicator, window}
}

func (sma smaIndicator) Calculate(index int) float64 {
	sum := 0.0
	for i := Max(0, index-sma.window+1); i <= index; i++ {
		sum += sma.indicator.Calculate(i)
	}
	realwindow := Min(sma.window, index+1)

	return sum / float64(realwindow)
}

type emaIndicator struct {
	Indicator
	window      int
	resultCache []float64
}

// Returns a new Exponential Moving Average Calculator
// http://www.investopedia.com/terms/e/ema.asp
func NewEMAIndicator(indicator Indicator, window int) Indicator {
	return &emaIndicator{
		Indicator:   indicator,
		window:      window,
		resultCache: make([]float64, window),
	}
}

func (ema *emaIndicator) Calculate(index int) float64 {
	if index+1 < ema.window {
		return smaIndicator{ema.Indicator, ema.window}.Calculate(index)
	}

	if index == 0 {
		result := ema.Indicator.Calculate(index)
		return result
	}

	emaPrev := ema.Calculate(index - 1)
	mult := 2.0 / float64(ema.window+1)
	result := (ema.Indicator.Calculate(index)-emaPrev)*mult + emaPrev

	return result
}

func (ema *emaIndicator) cacheResult(index int, val float64) {
	if index < len(ema.resultCache) {
		ema.resultCache[index] = val
	} else {
		ema.resultCache = append(ema.resultCache, val)
	}
}

func (ema emaIndicator) multiplier(index int) float64 {
	return 2.0 / (float64(index) + 1)
}

// Returns a new Moving Average Convergence-Divergence indicator
// http://www.investopedia.com/terms/m/macd.asp
func NewMACDIndicator(baseIndicator Indicator, shortwindow, longwindow int) Indicator {
	return NewDifferenceIndicator(NewEMAIndicator(baseIndicator, shortwindow), NewEMAIndicator(baseIndicator, longwindow))
}

// Returns a new Moving Average Convergence-Divergence histogram incicator, the result of which is
// the macd indicator minus it's @param signalLinewindow EMA
// http://stockcharts.com/school/doku.php?id=chart_school:technical_indicators:macd-histogram
func NewMACDHistogramIndicator(macdIdicator Indicator, signalLinewindow int) Indicator {
	return NewDifferenceIndicator(macdIdicator, NewEMAIndicator(macdIdicator, signalLinewindow))
}
