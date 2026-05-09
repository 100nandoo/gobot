package saaham

import (
	"testing"

	tele "gopkg.in/telebot.v3"
)

func TestAllowsPlainTextLookup(t *testing.T) {
	if !allowsPlainTextLookup(&tele.Chat{Type: tele.ChatPrivate}) {
		t.Fatal("expected private chats to allow plain-text lookups")
	}
	if !allowsPlainTextLookup(&tele.Chat{Type: tele.ChatGroup}) {
		t.Fatal("expected group chats to allow plain-text lookups")
	}
	if !allowsPlainTextLookup(&tele.Chat{Type: tele.ChatSuperGroup}) {
		t.Fatal("expected supergroup chats to allow plain-text lookups")
	}
}

func TestAllowsCommandLookup(t *testing.T) {
	if !allowsCommandLookup(&tele.Chat{Type: tele.ChatPrivate}) {
		t.Fatal("expected private chats to allow command lookups")
	}
	if !allowsCommandLookup(&tele.Chat{Type: tele.ChatGroup}) {
		t.Fatal("expected group chats to allow command lookups")
	}
	if !allowsCommandLookup(&tele.Chat{Type: tele.ChatSuperGroup}) {
		t.Fatal("expected supergroup chats to allow command lookups")
	}
}
