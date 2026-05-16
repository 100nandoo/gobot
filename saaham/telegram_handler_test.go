package saaham

import (
	"testing"
	"time"

	tele "gopkg.in/telebot.v3"
)

type fakeContext struct {
	chat       *tele.Chat
	sender     *tele.User
	text       string
	args       []string
	sendCalls  []string
	replyCalls []string
	store      map[string]interface{}
}

func (f *fakeContext) Bot() *tele.Bot                     { return nil }
func (f *fakeContext) Update() tele.Update                { return tele.Update{} }
func (f *fakeContext) Message() *tele.Message             { return nil }
func (f *fakeContext) Callback() *tele.Callback           { return nil }
func (f *fakeContext) Query() *tele.Query                 { return nil }
func (f *fakeContext) InlineResult() *tele.InlineResult   { return nil }
func (f *fakeContext) ShippingQuery() *tele.ShippingQuery { return nil }
func (f *fakeContext) PreCheckoutQuery() *tele.PreCheckoutQuery {
	return nil
}
func (f *fakeContext) Poll() *tele.Poll                       { return nil }
func (f *fakeContext) PollAnswer() *tele.PollAnswer           { return nil }
func (f *fakeContext) ChatMember() *tele.ChatMemberUpdate     { return nil }
func (f *fakeContext) ChatJoinRequest() *tele.ChatJoinRequest { return nil }
func (f *fakeContext) Migration() (int64, int64)              { return 0, 0 }
func (f *fakeContext) Topic() *tele.Topic                     { return nil }
func (f *fakeContext) Boost() *tele.BoostUpdated              { return nil }
func (f *fakeContext) BoostRemoved() *tele.BoostRemoved       { return nil }
func (f *fakeContext) Sender() *tele.User                     { return f.sender }
func (f *fakeContext) Chat() *tele.Chat                       { return f.chat }
func (f *fakeContext) Recipient() tele.Recipient {
	if f.chat != nil {
		return f.chat
	}
	return f.sender
}
func (f *fakeContext) Text() string            { return f.text }
func (f *fakeContext) Entities() tele.Entities { return nil }
func (f *fakeContext) Data() string            { return "" }
func (f *fakeContext) Args() []string          { return f.args }
func (f *fakeContext) Send(what interface{}, _ ...interface{}) error {
	f.sendCalls = append(f.sendCalls, what.(string))
	return nil
}
func (f *fakeContext) SendAlbum(_ tele.Album, _ ...interface{}) error { return nil }
func (f *fakeContext) Reply(what interface{}, _ ...interface{}) error {
	f.replyCalls = append(f.replyCalls, what.(string))
	return nil
}
func (f *fakeContext) Forward(_ tele.Editable, _ ...interface{}) error    { return nil }
func (f *fakeContext) ForwardTo(_ tele.Recipient, _ ...interface{}) error { return nil }
func (f *fakeContext) Edit(_ interface{}, _ ...interface{}) error         { return nil }
func (f *fakeContext) EditCaption(_ string, _ ...interface{}) error       { return nil }
func (f *fakeContext) EditOrSend(_ interface{}, _ ...interface{}) error   { return nil }
func (f *fakeContext) EditOrReply(_ interface{}, _ ...interface{}) error  { return nil }
func (f *fakeContext) Delete() error                                      { return nil }
func (f *fakeContext) DeleteAfter(_ time.Duration) *time.Timer            { return nil }
func (f *fakeContext) Notify(_ tele.ChatAction) error                     { return nil }
func (f *fakeContext) Ship(_ ...interface{}) error                        { return nil }
func (f *fakeContext) Accept(_ ...string) error                           { return nil }
func (f *fakeContext) Answer(_ *tele.QueryResponse) error                 { return nil }
func (f *fakeContext) Respond(_ ...*tele.CallbackResponse) error          { return nil }
func (f *fakeContext) RespondText(_ string) error                         { return nil }
func (f *fakeContext) RespondAlert(_ string) error                        { return nil }
func (f *fakeContext) Get(key string) interface{} {
	if f.store == nil {
		return nil
	}
	return f.store[key]
}
func (f *fakeContext) Set(key string, val interface{}) {
	if f.store == nil {
		f.store = map[string]interface{}{}
	}
	f.store[key] = val
}

type trackingQuoteService struct {
	lookups []string
	result  *QuoteResult
	err     error
}

func (s *trackingQuoteService) Lookup(symbol string) (*QuoteResult, error) {
	s.lookups = append(s.lookups, symbol)
	return s.result, s.err
}

func TestShouldIgnoreTriggerForBotSender(t *testing.T) {
	ctx := &fakeContext{
		sender: &tele.User{IsBot: true},
	}

	if !shouldIgnoreTrigger(ctx) {
		t.Fatal("expected bot-authored messages to be ignored")
	}
}

func TestShouldIgnoreTriggerAllowsHumanSender(t *testing.T) {
	ctx := &fakeContext{
		sender: &tele.User{IsBot: false},
	}

	if shouldIgnoreTrigger(ctx) {
		t.Fatal("expected human-authored messages to be handled")
	}
}

func TestSendQuoteMessageUsesReplyInGroups(t *testing.T) {
	ctx := &fakeContext{
		chat: &tele.Chat{Type: tele.ChatGroup},
	}

	if err := sendQuoteMessage(ctx, "group quote"); err != nil {
		t.Fatalf("sendQuoteMessage returned error: %v", err)
	}
	if len(ctx.replyCalls) != 1 || ctx.replyCalls[0] != "group quote" {
		t.Fatalf("expected one group reply call, got %#v", ctx.replyCalls)
	}
	if len(ctx.sendCalls) != 0 {
		t.Fatalf("expected no direct send calls for groups, got %#v", ctx.sendCalls)
	}
}

func TestSendQuoteMessageUsesSendInPrivateChats(t *testing.T) {
	ctx := &fakeContext{
		chat: &tele.Chat{Type: tele.ChatPrivate},
	}

	if err := sendQuoteMessage(ctx, "private quote"); err != nil {
		t.Fatalf("sendQuoteMessage returned error: %v", err)
	}
	if len(ctx.sendCalls) != 1 || ctx.sendCalls[0] != "private quote" {
		t.Fatalf("expected one private send call, got %#v", ctx.sendCalls)
	}
	if len(ctx.replyCalls) != 0 {
		t.Fatalf("expected no reply calls for private chats, got %#v", ctx.replyCalls)
	}
}

func TestPlainTextQuoteLookupIgnoresBotAuthoredGroupTickerToken(t *testing.T) {
	originalService := quoteService
	service := &trackingQuoteService{
		result: &QuoteResult{Symbol: "AAPL", InstrumentName: "Apple Inc."},
	}
	quoteService = service
	t.Cleanup(func() {
		quoteService = originalService
	})

	ctx := &fakeContext{
		chat:   &tele.Chat{Type: tele.ChatGroup},
		sender: &tele.User{IsBot: true},
		text:   "AAPL",
	}

	if err := plainTextQuoteLookup(ctx); err != nil {
		t.Fatalf("plainTextQuoteLookup returned error: %v", err)
	}
	if len(service.lookups) != 0 {
		t.Fatalf("expected no lookup for bot-authored plain-text trigger, got %#v", service.lookups)
	}
	if len(ctx.sendCalls) != 0 || len(ctx.replyCalls) != 0 {
		t.Fatalf("expected no reply activity for ignored trigger, send=%#v reply=%#v", ctx.sendCalls, ctx.replyCalls)
	}
}

func TestQuoteCommandIgnoresBotAuthoredGroupCommand(t *testing.T) {
	originalService := quoteService
	service := &trackingQuoteService{
		result: &QuoteResult{Symbol: "AAPL", InstrumentName: "Apple Inc."},
	}
	quoteService = service
	t.Cleanup(func() {
		quoteService = originalService
	})

	ctx := &fakeContext{
		chat:   &tele.Chat{Type: tele.ChatGroup},
		sender: &tele.User{IsBot: true},
		args:   []string{"AAPL"},
	}

	if err := quoteCommand(ctx); err != nil {
		t.Fatalf("quoteCommand returned error: %v", err)
	}
	if len(service.lookups) != 0 {
		t.Fatalf("expected no lookup for bot-authored command trigger, got %#v", service.lookups)
	}
	if len(ctx.sendCalls) != 0 || len(ctx.replyCalls) != 0 {
		t.Fatalf("expected no reply activity for ignored command, send=%#v reply=%#v", ctx.sendCalls, ctx.replyCalls)
	}
}

func TestConversionCommandIgnoresBotAuthoredGroupCommand(t *testing.T) {
	originalService := quoteService
	service := &trackingQuoteService{
		result: &QuoteResult{Symbol: "SGDIDR=X", InstrumentName: "SGD/IDR", Price: 12000},
	}
	quoteService = service
	t.Cleanup(func() {
		quoteService = originalService
	})

	ctx := &fakeContext{
		chat:   &tele.Chat{Type: tele.ChatGroup},
		sender: &tele.User{IsBot: true},
		args:   []string{"100", "SGD", "IDR"},
	}

	if err := conversionCommand(ctx); err != nil {
		t.Fatalf("conversionCommand returned error: %v", err)
	}
	if len(service.lookups) != 0 {
		t.Fatalf("expected no lookup for bot-authored conversion command, got %#v", service.lookups)
	}
	if len(ctx.sendCalls) != 0 || len(ctx.replyCalls) != 0 {
		t.Fatalf("expected no reply activity for ignored command, send=%#v reply=%#v", ctx.sendCalls, ctx.replyCalls)
	}
}

func TestConversionCommandFormatsReply(t *testing.T) {
	originalService := quoteService
	service := &trackingQuoteService{
		result: &QuoteResult{Symbol: "SGDIDR=X", InstrumentName: "SGD/IDR", Price: 12000},
	}
	quoteService = service
	t.Cleanup(func() {
		quoteService = originalService
	})

	ctx := &fakeContext{
		chat:   &tele.Chat{Type: tele.ChatPrivate},
		sender: &tele.User{IsBot: false},
		args:   []string{"100", "SGD", "to", "IDR"},
	}

	if err := conversionCommand(ctx); err != nil {
		t.Fatalf("conversionCommand returned error: %v", err)
	}
	if len(service.lookups) != 1 || service.lookups[0] != "SGDIDR=X" {
		t.Fatalf("expected SGDIDR=X lookup, got %#v", service.lookups)
	}
	if len(ctx.sendCalls) != 1 {
		t.Fatalf("expected one conversion reply, got %#v", ctx.sendCalls)
	}
	if ctx.sendCalls[0] != "*100.00 SGD = 1,200,000 IDR*\nRate: `SGDIDR=X 12,000.0000`" {
		t.Fatalf("unexpected conversion reply %q", ctx.sendCalls[0])
	}
}

func TestConversionCommandReturnsUsageForInvalidInput(t *testing.T) {
	ctx := &fakeContext{
		chat:   &tele.Chat{Type: tele.ChatPrivate},
		sender: &tele.User{IsBot: false},
		args:   []string{"hello", "world"},
	}

	if err := conversionCommand(ctx); err != nil {
		t.Fatalf("conversionCommand returned error: %v", err)
	}
	if len(ctx.sendCalls) != 1 || ctx.sendCalls[0] != conversionUsage {
		t.Fatalf("expected usage reply, got %#v", ctx.sendCalls)
	}
}

func TestPlainTextQuoteLookupSupportsPrivateChatConversion(t *testing.T) {
	originalService := quoteService
	service := &trackingQuoteService{
		result: &QuoteResult{Symbol: "SGDIDR=X", InstrumentName: "SGD/IDR", Price: 12000},
	}
	quoteService = service
	t.Cleanup(func() {
		quoteService = originalService
	})

	ctx := &fakeContext{
		chat:   &tele.Chat{Type: tele.ChatPrivate},
		sender: &tele.User{IsBot: false},
		text:   "100 sgd idr",
	}

	if err := plainTextQuoteLookup(ctx); err != nil {
		t.Fatalf("plainTextQuoteLookup returned error: %v", err)
	}
	if len(service.lookups) != 1 || service.lookups[0] != "SGDIDR=X" {
		t.Fatalf("expected SGDIDR=X lookup, got %#v", service.lookups)
	}
	if len(ctx.sendCalls) != 1 || ctx.sendCalls[0] != "*100.00 SGD = 1,200,000 IDR*\nRate: `SGDIDR=X 12,000.0000`" {
		t.Fatalf("unexpected private conversion reply %#v", ctx.sendCalls)
	}
}

func TestPlainTextQuoteLookupKeepsGroupConversionCommandOnly(t *testing.T) {
	originalService := quoteService
	service := &trackingQuoteService{
		result: &QuoteResult{Symbol: "SGDIDR=X", InstrumentName: "SGD/IDR", Price: 12000},
	}
	quoteService = service
	t.Cleanup(func() {
		quoteService = originalService
	})

	ctx := &fakeContext{
		chat:   &tele.Chat{Type: tele.ChatGroup},
		sender: &tele.User{IsBot: false},
		text:   "100 sgd idr",
	}

	if err := plainTextQuoteLookup(ctx); err != nil {
		t.Fatalf("plainTextQuoteLookup returned error: %v", err)
	}
	if len(service.lookups) != 0 {
		t.Fatalf("expected no lookup for group plain-text conversion, got %#v", service.lookups)
	}
	if len(ctx.sendCalls) != 0 || len(ctx.replyCalls) != 0 {
		t.Fatalf("expected no group reply activity, send=%#v reply=%#v", ctx.sendCalls, ctx.replyCalls)
	}
}
