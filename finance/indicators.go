package finance

import (
	"math"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/trend"
	"github.com/cinar/indicator/v2/volatility"
)

type IndicatorResult struct {
	CurrentPrice float64
	PrevClose    float64
	PriceChange  float64
	Currency     string

	RSI       float64
	RSISignal string

	MACDLine      float64
	MACDSignalVal float64
	MACDHistogram float64
	MACDSignal    string

	SMA50    float64
	SMA200   float64
	SMASignal string

	BollingerUpper  float64
	BollingerMid    float64
	BollingerLower  float64
	BollingerSignal string

	RSIScore  int
	MACDScore int
	SMAScore  int
	BBScore   int

	CompositeScore int
	Recommendation string
}

func computeIndicators(prices *PriceData) *IndicatorResult {
	closes := prices.Close
	result := &IndicatorResult{
		CurrentPrice: prices.Price,
		PrevClose:    prices.PrevClose,
		Currency:     prices.Currency,
	}

	if prices.PrevClose > 0 {
		result.PriceChange = ((prices.Price - prices.PrevClose) / prices.PrevClose) * 100
	}

	// RSI (14-period)
	rsi := momentum.NewRsi[float64]()
	rsiValues := helper.ChanToSlice(rsi.Compute(helper.SliceToChan(closes)))
	if len(rsiValues) > 0 {
		result.RSI = rsiValues[len(rsiValues)-1]
	}
	switch {
	case result.RSI < 30:
		result.RSISignal = "Oversold"
	case result.RSI > 70:
		result.RSISignal = "Overbought"
	default:
		result.RSISignal = "Neutral"
	}

	// MACD (12, 26, 9)
	macd := trend.NewMacd[float64]()
	macdLine, signalLine := macd.Compute(helper.SliceToChan(closes))
	var macdValues, signalValues []float64
	done := make(chan struct{})
	go func() {
		signalValues = helper.ChanToSlice(signalLine)
		close(done)
	}()
	macdValues = helper.ChanToSlice(macdLine)
	<-done

	if len(macdValues) > 0 && len(signalValues) > 0 {
		result.MACDLine = macdValues[len(macdValues)-1]
		result.MACDSignalVal = signalValues[len(signalValues)-1]
		result.MACDHistogram = result.MACDLine - result.MACDSignalVal

		if len(macdValues) > 1 && len(signalValues) > 1 {
			prevHistogram := macdValues[len(macdValues)-2] - signalValues[len(signalValues)-2]
			if prevHistogram <= 0 && result.MACDHistogram > 0 {
				result.MACDSignal = "Bullish Crossover"
			} else if prevHistogram >= 0 && result.MACDHistogram < 0 {
				result.MACDSignal = "Bearish Crossover"
			} else {
				result.MACDSignal = "Neutral"
			}
		}
	}

	// SMA 50 & 200
	sma50 := trend.NewSmaWithPeriod[float64](50)
	sma50Values := helper.ChanToSlice(sma50.Compute(helper.SliceToChan(closes)))
	if len(sma50Values) > 0 {
		result.SMA50 = sma50Values[len(sma50Values)-1]
	}

	sma200 := trend.NewSmaWithPeriod[float64](200)
	sma200Values := helper.ChanToSlice(sma200.Compute(helper.SliceToChan(closes)))
	if len(sma200Values) > 0 {
		result.SMA200 = sma200Values[len(sma200Values)-1]
	}

	price := prices.Price
	switch {
	case price > result.SMA50 && result.SMA50 > result.SMA200:
		result.SMASignal = "Strong Uptrend"
	case price < result.SMA50 && result.SMA50 < result.SMA200:
		result.SMASignal = "Strong Downtrend"
	case price > result.SMA50 && price > result.SMA200:
		result.SMASignal = "Above Both MAs"
	case price < result.SMA50 && price < result.SMA200:
		result.SMASignal = "Below Both MAs"
	default:
		result.SMASignal = "Mixed"
	}

	// Bollinger Bands (20, 2)
	bb := volatility.NewBollingerBands[float64]()
	upperCh, midCh, lowerCh := bb.Compute(helper.SliceToChan(closes))
	var upperValues, midValues, lowerValues []float64
	bbDone := make(chan struct{})
	bbDone2 := make(chan struct{})
	go func() {
		midValues = helper.ChanToSlice(midCh)
		close(bbDone)
	}()
	go func() {
		lowerValues = helper.ChanToSlice(lowerCh)
		close(bbDone2)
	}()
	upperValues = helper.ChanToSlice(upperCh)
	<-bbDone
	<-bbDone2

	if len(upperValues) > 0 && len(midValues) > 0 && len(lowerValues) > 0 {
		result.BollingerUpper = upperValues[len(upperValues)-1]
		result.BollingerMid = midValues[len(midValues)-1]
		result.BollingerLower = lowerValues[len(lowerValues)-1]

		lowerDist := math.Abs(price-result.BollingerLower) / result.BollingerLower
		upperDist := math.Abs(price-result.BollingerUpper) / result.BollingerUpper

		switch {
		case lowerDist < 0.01:
			result.BollingerSignal = "Near Lower Band"
		case upperDist < 0.01:
			result.BollingerSignal = "Near Upper Band"
		case price < result.BollingerLower:
			result.BollingerSignal = "Below Lower Band"
		case price > result.BollingerUpper:
			result.BollingerSignal = "Above Upper Band"
		default:
			result.BollingerSignal = "Within Bands"
		}
	}

	calculateCompositeScore(result)

	return result
}

func calculateCompositeScore(r *IndicatorResult) {
	switch {
	case r.RSI < 30:
		r.RSIScore = 3
	case r.RSI < 40:
		r.RSIScore = 1
	case r.RSI > 70:
		r.RSIScore = -3
	case r.RSI > 60:
		r.RSIScore = -1
	}

	switch r.MACDSignal {
	case "Bullish Crossover":
		r.MACDScore = 2
	case "Bearish Crossover":
		r.MACDScore = -2
	}

	switch r.SMASignal {
	case "Strong Uptrend":
		r.SMAScore = 2
	case "Strong Downtrend":
		r.SMAScore = -2
	}

	switch r.BollingerSignal {
	case "Near Lower Band", "Below Lower Band":
		r.BBScore = 2
	case "Near Upper Band", "Above Upper Band":
		r.BBScore = -2
	}

	r.CompositeScore = r.RSIScore + r.MACDScore + r.SMAScore + r.BBScore

	switch {
	case r.CompositeScore >= 5:
		r.Recommendation = "Strong Buy"
	case r.CompositeScore >= 2:
		r.Recommendation = "Buy"
	case r.CompositeScore >= -1:
		r.Recommendation = "Neutral"
	case r.CompositeScore >= -4:
		r.Recommendation = "Sell"
	default:
		r.Recommendation = "Strong Sell"
	}
}

func AnalyzeETF(symbol string) (*IndicatorResult, error) {
	prices, err := fetchPriceData(symbol)
	if err != nil {
		return nil, err
	}

	result := computeIndicators(prices)
	return result, nil
}
