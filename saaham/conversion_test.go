package saaham

import (
	"strings"
	"testing"
)

func TestParseConversionCommand(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		want  *ConversionRequest
		isErr bool
	}{
		{
			name: "three-part command",
			args: []string{"100", "sgd", "idr"},
			want: &ConversionRequest{Amount: 100, BaseCurrency: "SGD", QuoteCurrency: "IDR"},
		},
		{
			name: "four-part command with to",
			args: []string{"100", "sgd", "to", "idr"},
			want: &ConversionRequest{Amount: 100, BaseCurrency: "SGD", QuoteCurrency: "IDR"},
		},
		{
			name:  "missing arguments",
			args:  nil,
			want:  nil,
			isErr: false,
		},
		{
			name:  "invalid amount",
			args:  []string{"abc", "sgd", "idr"},
			isErr: true,
		},
		{
			name:  "invalid code",
			args:  []string{"100", "singapore", "idr"},
			isErr: true,
		},
		{
			name:  "invalid shape",
			args:  []string{"100", "sgd", "into", "idr"},
			isErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseConversionCommand(tc.args)
			if tc.isErr {
				if err == nil {
					t.Fatal("expected parse error")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.want == nil {
				if got != nil {
					t.Fatalf("expected nil request, got %#v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("expected request, got nil")
			}
			if *got != *tc.want {
				t.Fatalf("expected %#v, got %#v", *tc.want, *got)
			}
		})
	}
}

func TestParsePlainTextConversion(t *testing.T) {
	request, ok := parsePlainTextConversion("100 sgd idr")
	if !ok {
		t.Fatal("expected plain-text conversion to parse")
	}
	if request.Amount != 100 || request.BaseCurrency != "SGD" || request.QuoteCurrency != "IDR" {
		t.Fatalf("unexpected request %#v", request)
	}

	if _, ok := parsePlainTextConversion("show me 100 sgd idr"); ok {
		t.Fatal("expected mixed prose to be rejected as plain-text conversion")
	}
}

func TestLookupConversionUsesDirectYahooPair(t *testing.T) {
	service := stubQuoteService{
		results: map[string]*QuoteResult{
			"SGDIDR=X": {Symbol: "SGDIDR=X", InstrumentName: "SGD/IDR", Price: 12000},
		},
	}

	result, convertedAmount, err := lookupConversion(service, &ConversionRequest{
		Amount:        100,
		BaseCurrency:  "SGD",
		QuoteCurrency: "IDR",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || result.Symbol != "SGDIDR=X" {
		t.Fatalf("expected SGDIDR=X result, got %#v", result)
	}
	if convertedAmount != 1200000 {
		t.Fatalf("expected converted amount 1200000, got %v", convertedAmount)
	}
}

func TestFormatConversionReplyIncludesAmountAndPair(t *testing.T) {
	reply := formatConversionReply(
		&ConversionRequest{Amount: 100, BaseCurrency: "SGD", QuoteCurrency: "IDR"},
		&QuoteResult{Symbol: "SGDIDR=X", Price: 12000},
		1200000,
	)

	if !strings.Contains(reply, "*100.00 SGD = 1,200,000 IDR*") {
		t.Fatalf("expected converted amount in reply, got %q", reply)
	}
	if !strings.Contains(reply, "Rate: `SGDIDR=X 12,000.0000`") {
		t.Fatalf("expected rate line in reply, got %q", reply)
	}
}

func TestConversionAmountDisplayDecimals(t *testing.T) {
	if got := conversionAmountDisplayDecimals("IDR"); got != 0 {
		t.Fatalf("expected IDR to use 0 decimals, got %d", got)
	}
	if got := conversionAmountDisplayDecimals("USD"); got != 2 {
		t.Fatalf("expected USD to use 2 decimals, got %d", got)
	}
}
