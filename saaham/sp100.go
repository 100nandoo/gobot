package saaham

import (
	"bufio"
	_ "embed"
	"strings"
)

//go:embed sp100_snapshot.txt
var sp100SnapshotFile string

var sp100Snapshot = mustLoadSP100Snapshot()

type tickerSnapshot struct {
	symbols map[string]struct{}
}

func mustLoadSP100Snapshot() tickerSnapshot {
	return loadTickerSnapshot(sp100SnapshotFile)
}

func loadTickerSnapshot(raw string) tickerSnapshot {
	snapshot := tickerSnapshot{
		symbols: make(map[string]struct{}),
	}

	scanner := bufio.NewScanner(strings.NewReader(raw))
	for scanner.Scan() {
		symbol, ok := normalizeTickerToken(scanner.Text())
		if !ok {
			continue
		}
		snapshot.symbols[symbol] = struct{}{}
	}

	return snapshot
}

func (s tickerSnapshot) Contains(symbol string) bool {
	normalized, ok := normalizeTickerToken(symbol)
	if !ok {
		return false
	}

	_, exists := s.symbols[normalized]
	return exists
}
