package finance

import (
	"encoding/json"
	"fmt"
	"gobot/config"
	"gobot/pkg"
	"os"
	"slices"
	"strconv"
	"strings"
)

const financeSettingsTableName = "FinanceWatchlists"

const maxWatchlistSize = 8
const defaultAlertThreshold = 5

type FinanceSettingsRow struct {
	ChatID         int64  `json:"chat_id"`
	Symbols        string `json:"symbols"`
	AlertThreshold int    `json:"alert_threshold"`
	ScoreState     string `json:"score_state"`
}

type FinanceSettings struct {
	Symbols        []ETFSymbol
	AlertThreshold int
	LastScores     map[string]int
}

func defaultWatchSymbols() []ETFSymbol {
	return append([]ETFSymbol(nil), DefaultSymbols...)
}

func defaultFinanceSettings() FinanceSettings {
	return FinanceSettings{
		Symbols:        defaultWatchSymbols(),
		AlertThreshold: defaultAlertThreshold,
		LastScores:     map[string]int{},
	}
}

func normalizeTicker(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func parseTickers(values []string) []string {
	var normalized []string
	seen := make(map[string]struct{})

	for _, value := range values {
		ticker := normalizeTicker(value)
		if ticker == "" {
			continue
		}
		if _, ok := seen[ticker]; ok {
			continue
		}
		seen[ticker] = struct{}{}
		normalized = append(normalized, ticker)
	}

	return normalized
}

func buildSymbols(tickers []string) []ETFSymbol {
	symbols := make([]ETFSymbol, 0, len(tickers))
	for _, ticker := range tickers {
		symbols = append(symbols, ETFSymbol{Ticker: ticker})
	}
	return symbols
}

func serializeTickers(tickers []string) string {
	return strings.Join(tickers, ",")
}

func deserializeTickers(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	return parseTickers(strings.Split(raw, ","))
}

func usePersistentWatchlistStore() bool {
	return os.Getenv(config.SupabaseUrl) != "" && os.Getenv(config.SupabaseKey) != ""
}

func deserializeScoreState(raw string) map[string]int {
	if strings.TrimSpace(raw) == "" {
		return map[string]int{}
	}

	var scores map[string]int
	if err := json.Unmarshal([]byte(raw), &scores); err != nil {
		pkg.LogWithTimestamp("Error parsing finance score state: %v", err)
		return map[string]int{}
	}

	normalized := make(map[string]int, len(scores))
	for ticker, score := range scores {
		normalized[normalizeTicker(ticker)] = score
	}
	return normalized
}

func serializeScoreState(scores map[string]int) string {
	if len(scores) == 0 {
		return "{}"
	}

	payload, err := json.Marshal(scores)
	if err != nil {
		pkg.LogWithTimestamp("Error serializing finance score state: %v", err)
		return "{}"
	}

	return string(payload)
}

func filterScoreStateForSymbols(scores map[string]int, symbols []ETFSymbol) map[string]int {
	filtered := make(map[string]int, len(symbols))
	for _, symbol := range symbols {
		if score, ok := scores[symbol.Ticker]; ok {
			filtered[symbol.Ticker] = score
		}
	}
	return filtered
}

func normalizeAlertThreshold(value int) int {
	switch {
	case value > 10:
		return 10
	case value < -10:
		return -10
	default:
		return value
	}
}

func settingsFromRow(row FinanceSettingsRow) FinanceSettings {
	symbols := buildSymbols(deserializeTickers(row.Symbols))
	if len(symbols) == 0 {
		symbols = defaultWatchSymbols()
	}

	return FinanceSettings{
		Symbols:        symbols,
		AlertThreshold: normalizeAlertThreshold(row.AlertThreshold),
		LastScores:     filterScoreStateForSymbols(deserializeScoreState(row.ScoreState), symbols),
	}
}

func saveFinanceSettings(chatID int64, settings FinanceSettings) error {
	settings.AlertThreshold = normalizeAlertThreshold(settings.AlertThreshold)
	if len(settings.Symbols) == 0 {
		settings.Symbols = defaultWatchSymbols()
	}
	if settings.LastScores == nil {
		settings.LastScores = map[string]int{}
	}
	settings.LastScores = filterScoreStateForSymbols(settings.LastScores, settings.Symbols)

	if !usePersistentWatchlistStore() {
		return fmt.Errorf("supabase is not configured")
	}

	var deleted []FinanceSettingsRow
	if err := pkg.SupabaseClient.DB.From(financeSettingsTableName).Delete().Eq("chat_id", strconv.FormatInt(chatID, 10)).Execute(&deleted); err != nil {
		pkg.LogWithTimestamp("Error deleting finance settings: %v", err)
	}

	tickers := make([]string, 0, len(settings.Symbols))
	for _, symbol := range settings.Symbols {
		tickers = append(tickers, symbol.Ticker)
	}

	var inserted []FinanceSettingsRow
	err := pkg.SupabaseClient.DB.From(financeSettingsTableName).Insert(FinanceSettingsRow{
		ChatID:         chatID,
		Symbols:        serializeTickers(tickers),
		AlertThreshold: settings.AlertThreshold,
		ScoreState:     serializeScoreState(settings.LastScores),
	}).Execute(&inserted)
	if err != nil {
		return fmt.Errorf("saving finance settings: %w", err)
	}

	return nil
}

func GetFinanceSettings(chatID int64) FinanceSettings {
	if !usePersistentWatchlistStore() {
		return defaultFinanceSettings()
	}

	var rows []FinanceSettingsRow
	err := pkg.SupabaseClient.DB.From(financeSettingsTableName).Select("*").Eq("chat_id", strconv.FormatInt(chatID, 10)).Execute(&rows)
	if err != nil {
		pkg.LogWithTimestamp("Error loading finance settings: %v", err)
		return defaultFinanceSettings()
	}

	if len(rows) == 0 {
		return defaultFinanceSettings()
	}

	return settingsFromRow(rows[0])
}

func GetWatchlist(chatID int64) []ETFSymbol {
	return GetFinanceSettings(chatID).Symbols
}

func GetAlertThreshold(chatID int64) int {
	return GetFinanceSettings(chatID).AlertThreshold
}

func SaveAlertThreshold(chatID int64, threshold int) error {
	settings := GetFinanceSettings(chatID)
	settings.AlertThreshold = normalizeAlertThreshold(threshold)
	return saveFinanceSettings(chatID, settings)
}

func ResetAlertThreshold(chatID int64) error {
	settings := GetFinanceSettings(chatID)
	settings.AlertThreshold = defaultAlertThreshold
	return saveFinanceSettings(chatID, settings)
}

func SaveWatchlist(chatID int64, tickers []string) error {
	normalized := parseTickers(tickers)
	if len(normalized) == 0 {
		return fmt.Errorf("watchlist cannot be empty")
	}
	if len(normalized) > maxWatchlistSize {
		return fmt.Errorf("watchlist supports up to %d tickers", maxWatchlistSize)
	}

	settings := GetFinanceSettings(chatID)
	settings.Symbols = buildSymbols(normalized)
	return saveFinanceSettings(chatID, settings)
}

func ResetWatchlist(chatID int64) error {
	settings := GetFinanceSettings(chatID)
	settings.Symbols = defaultWatchSymbols()
	return saveFinanceSettings(chatID, settings)
}

func AddToWatchlist(chatID int64, tickers []string) ([]ETFSymbol, error) {
	current := GetWatchlist(chatID)
	merged := make([]string, 0, len(current)+len(tickers))
	for _, symbol := range current {
		merged = append(merged, symbol.Ticker)
	}
	merged = append(merged, tickers...)

	if err := SaveWatchlist(chatID, merged); err != nil {
		return nil, err
	}

	return GetWatchlist(chatID), nil
}

func RemoveFromWatchlist(chatID int64, tickers []string) ([]ETFSymbol, error) {
	toRemove := parseTickers(tickers)
	if len(toRemove) == 0 {
		return GetWatchlist(chatID), nil
	}

	current := GetWatchlist(chatID)
	filtered := make([]string, 0, len(current))
	for _, symbol := range current {
		if slices.Contains(toRemove, symbol.Ticker) {
			continue
		}
		filtered = append(filtered, symbol.Ticker)
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("watchlist cannot be empty")
	}

	if err := SaveWatchlist(chatID, filtered); err != nil {
		return nil, err
	}

	return GetWatchlist(chatID), nil
}

func ValidateTickers(tickers []string) error {
	for _, ticker := range parseTickers(tickers) {
		if _, err := AnalyzeETF(ticker); err != nil {
			return fmt.Errorf("%s is invalid or has insufficient data: %w", ticker, err)
		}
	}

	return nil
}

func UpdateLastScore(chatID int64, ticker string, score int) error {
	settings := GetFinanceSettings(chatID)
	if settings.LastScores == nil {
		settings.LastScores = map[string]int{}
	}
	settings.LastScores[normalizeTicker(ticker)] = score
	return saveFinanceSettings(chatID, settings)
}
