package saaham

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

const conversionUsage = "Usage: `/fx 100 SGD IDR`"

type ConversionRequest struct {
	Amount        float64
	BaseCurrency  string
	QuoteCurrency string
}

func (r ConversionRequest) PairSymbol() string {
	return r.BaseCurrency + r.QuoteCurrency + "=X"
}

func parseConversionCommand(args []string) (*ConversionRequest, error) {
	if len(args) == 0 {
		return nil, nil
	}
	return parseConversionParts(args)
}

func parsePlainTextConversion(text string) (*ConversionRequest, bool) {
	request, err := parseConversionParts(strings.Fields(text))
	if err != nil || request == nil {
		return nil, false
	}
	return request, true
}

func parseConversionParts(args []string) (*ConversionRequest, error) {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		trimmed := strings.TrimSpace(arg)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}

	switch {
	case len(parts) == 3:
	case len(parts) == 4 && strings.EqualFold(parts[2], "to"):
		parts = []string{parts[0], parts[1], parts[3]}
	default:
		return nil, fmt.Errorf("invalid conversion command")
	}

	amount, err := strconv.ParseFloat(parts[0], 64)
	if err != nil || amount <= 0 {
		return nil, fmt.Errorf("invalid conversion amount")
	}

	base := strings.ToUpper(parts[1])
	quote := strings.ToUpper(parts[2])
	if !isISOCurrencyCode(base) || !isISOCurrencyCode(quote) {
		return nil, fmt.Errorf("invalid currency code")
	}

	return &ConversionRequest{
		Amount:        amount,
		BaseCurrency:  base,
		QuoteCurrency: quote,
	}, nil
}

func isISOCurrencyCode(value string) bool {
	if len(value) != 3 {
		return false
	}
	for _, r := range value {
		if !unicode.IsLetter(r) || !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func lookupConversion(service QuoteService, request *ConversionRequest) (*QuoteResult, float64, error) {
	result, err := service.Lookup(request.PairSymbol())
	if err != nil {
		return nil, 0, err
	}

	return result, request.Amount * result.Price, nil
}

func formatConversionReply(request *ConversionRequest, result *QuoteResult, convertedAmount float64) string {
	convertedDecimals := conversionAmountDisplayDecimals(request.QuoteCurrency)
	return fmt.Sprintf(
		"*%s %s = %s %s*\nRate: `%s %s`",
		formatDisplayNumber(request.Amount, 2),
		request.BaseCurrency,
		formatDisplayNumber(convertedAmount, convertedDecimals),
		request.QuoteCurrency,
		result.Symbol,
		formatDisplayNumber(result.Price, 4),
	)
}

func conversionAmountDisplayDecimals(quoteCurrency string) int {
	switch strings.ToUpper(strings.TrimSpace(quoteCurrency)) {
	case "IDR":
		return 0
	default:
		return 2
	}
}

func formatDisplayNumber(value float64, decimals int) string {
	negative := value < 0
	if negative {
		value = math.Abs(value)
	}

	raw := strconv.FormatFloat(value, 'f', decimals, 64)
	parts := strings.SplitN(raw, ".", 2)
	intPart := addThousandsSeparators(parts[0])
	if negative {
		intPart = "-" + intPart
	}
	if len(parts) == 1 {
		return intPart
	}
	return intPart + "." + parts[1]
}

func addThousandsSeparators(value string) string {
	if len(value) <= 3 {
		return value
	}

	var b strings.Builder
	prefixLen := len(value) % 3
	if prefixLen == 0 {
		prefixLen = 3
	}
	b.WriteString(value[:prefixLen])
	for i := prefixLen; i < len(value); i += 3 {
		b.WriteByte(',')
		b.WriteString(value[i : i+3])
	}
	return b.String()
}
