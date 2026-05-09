package saaham

import (
	"strings"
	"testing"

	tele "gopkg.in/telebot.v3"
)

func TestPlainTextLookupStaysSingleSymbolOnly(t *testing.T) {
	if symbol, ok := parsePlainTextLookup("(aapl)"); !ok || symbol != "AAPL" {
		t.Fatalf("expected private plain-text lookup to accept one wrapped ticker, got %q ok=%v", symbol, ok)
	}

	if _, ok := parsePlainTextLookup("AAPL MSFT"); ok {
		t.Fatal("expected plain-text batch input to be rejected")
	}

	if _, ok := parsePlainTextLookup("show me AAPL"); ok {
		t.Fatal("expected mixed prose to be rejected as a plain-text lookup trigger")
	}
}

func TestCommandBatchBehaviorRemainsCommandOnly(t *testing.T) {
	got := parseBatchCommandLookups([]string{"aapl", "MSFT"})
	if len(got) != 2 || got[0] != "AAPL" || got[1] != "MSFT" {
		t.Fatalf("expected command batch parsing to preserve input order, got %#v", got)
	}
}

func TestChatScopeRules(t *testing.T) {
	privateChat := &tele.Chat{Type: tele.ChatPrivate}
	groupChat := &tele.Chat{Type: tele.ChatGroup}
	supergroupChat := &tele.Chat{Type: tele.ChatSuperGroup}

	if !allowsPlainTextLookup(privateChat) {
		t.Fatal("expected private chats to allow plain-text lookup triggers")
	}
	if !allowsPlainTextLookup(groupChat) {
		t.Fatal("expected group chats to allow plain-text lookup triggers")
	}
	if !allowsPlainTextLookup(supergroupChat) {
		t.Fatal("expected supergroup chats to allow plain-text lookup triggers")
	}
	if !allowsCommandLookup(groupChat) {
		t.Fatal("expected group chats to allow explicit /q commands")
	}
	if !usesQuotedGroupReply(groupChat) {
		t.Fatal("expected group chats to use quoted replies")
	}
	if !usesQuotedGroupReply(supergroupChat) {
		t.Fatal("expected supergroup chats to use quoted replies")
	}
	if usesQuotedGroupReply(privateChat) {
		t.Fatal("expected private chats to avoid quoted replies")
	}
}

func TestFormatBatchQuoteReplySupportsMixedOutcomes(t *testing.T) {
	reply := formatBatchQuoteReply([]BatchQuoteResult{
		{
			Query: "AAPL",
			Result: &QuoteResult{
				Symbol:         "AAPL",
				InstrumentName: "Apple Inc.",
				Currency:       "USD",
				Price:          210.15,
				Change:         3.14,
				ChangePercent:  1.52,
			},
		},
		{
			Query: "BAD",
			Err:   &LookupError{Kind: LookupErrorInvalidSymbol, Symbol: "BAD"},
		},
		{
			Query: "QQQJ",
			Err:   &LookupError{Kind: LookupErrorUnsupported, Symbol: "QQQJ"},
		},
		{
			Query: "MSFT",
			Err:   &LookupError{Kind: LookupErrorProviderFailed, Symbol: "MSFT"},
		},
	})

	parts := strings.Split(reply, "\n\n")
	if len(parts) != 4 {
		t.Fatalf("expected 4 result blocks, got %d", len(parts))
	}
	if !strings.Contains(parts[0], "*AAPL* - Apple Inc.") {
		t.Fatalf("expected first block to be AAPL success, got %q", parts[0])
	}
	if parts[1] != "*BAD* - invalid symbol" {
		t.Fatalf("expected invalid symbol block, got %q", parts[1])
	}
	if parts[2] != "*QQQJ* - unsupported instrument" {
		t.Fatalf("expected unsupported instrument block, got %q", parts[2])
	}
	if parts[3] != "*MSFT* - provider failure" {
		t.Fatalf("expected provider failure block, got %q", parts[3])
	}
}
