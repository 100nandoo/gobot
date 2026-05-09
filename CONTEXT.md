# Gobot

Gobot is a collection of Telegram bots for distinct information domains. Each bot has its own user-facing scope even when multiple bots share integration code or data sources.

## Language

**Finance Watchlist Bot**:
A Telegram bot that analyzes a saved watchlist and produces indicator-driven guidance for stocks and ETFs.
_Avoid_: quote bot, price bot

**Quote Bot**:
A Telegram bot that returns the latest market price for a requested ticker without watchlist analysis.
_Avoid_: finance bot, analyzer

**Quote Lookup**:
A user request for the latest market price and daily change for a single ticker.
_Avoid_: analysis, watchlist scan

**Exact Ticker Symbol**:
The Yahoo Finance symbol the user must provide verbatim for a **Quote Lookup** to succeed.
_Avoid_: company name, fuzzy match

**Lookup Trigger**:
The user input pattern that the **Quote Bot** treats as a request for a **Quote Lookup**.
_Avoid_: arbitrary chat text

**Ticker Token**:
An input string containing one **Exact Ticker Symbol**, optionally wrapped with lightweight punctuation in plain-text messages.
_Avoid_: sentence, multi-symbol request

**Wrapper Punctuation**:
Leading or trailing punctuation around a **Ticker Token** that the **Quote Bot** may strip without changing the ticker body.
_Avoid_: `$` prefix, embedded separators

**Invalid Symbol Response**:
The bot reply used when a **Quote Lookup** does not match a valid Yahoo Finance ticker symbol.
_Avoid_: transient error, retry message

**Provider Failure Response**:
The bot reply used when the market data provider cannot return quote data for a valid lookup attempt.
_Avoid_: invalid ticker guidance

**Canonical Symbol Reply**:
The bot response format that shows the normalized ticker symbol actually used for the lookup.
_Avoid_: raw user input echo

**Instrument Name**:
The provider-sourced human-readable stock or ETF name shown alongside a successful quote.
_Avoid_: user-supplied label, guessed description

**Stateless Quote Bot**:
A **Quote Bot** that answers each **Quote Lookup** without persisting user settings, watchlists, or lookup history.
_Avoid_: watchlist bot, personalized bot

**Dedicated Bot Identity**:
A Telegram bot token and username used only by one bot product.
_Avoid_: shared bot token, shared bot persona

**Quote Command**:
The explicit Telegram command that triggers a **Quote Lookup**.
_Avoid_: generic finance command, stock-only command

**Supported Instrument**:
An asset the **Quote Bot** is willing to quote in v1 based on provider-reported instrument type.
_Avoid_: any quoteable symbol, unsupported asset

**Unsupported Instrument Response**:
The bot reply used when a ticker resolves successfully but the instrument type is outside the bot's supported scope.
_Avoid_: invalid ticker response, provider failure response

**Shared Gobot Runtime**:
The single `gobot` process that hosts multiple Telegram bot modules, each with its own bot identity.
_Avoid_: separate deployable, separate service

**Generic Quote Interface**:
A quote bot interface that relies on the canonical **Quote Command** and plain-text ticker input rather than symbol-specific shortcut commands.
_Avoid_: watchlist shortcut set, curated command list

**English Quote UX**:
The English-only command help and response language used by the **Quote Bot** in v1.
_Avoid_: bilingual UX, Indonesian-only UX

**Long Polling Runtime**:
The Telegram polling model where the bot fetches updates directly rather than receiving webhook callbacks.
_Avoid_: webhook deployment

**SAAHAM Bot**:
The branded product name of the Quote Bot and the source of its dedicated bot token name.
_Avoid_: assuming `SAHAM` is the intended spelling

**Branded Code Module**:
A Go package and folder name that intentionally uses the product brand rather than a generic descriptive name.
_Avoid_: generic quote package, hidden brand boundary

**Chat Scope Rule**:
The rule that determines which lookup triggers the **Quote Bot** accepts in private chats versus group chats.
_Avoid_: identical trigger behavior everywhere

**Public Group Reply**:
The normal in-chat response the **Quote Bot** posts when an explicit group command triggers a **Quote Lookup**.
_Avoid_: pseudo-private group response, disabled group command reply

**Batch Quote Request**:
A single user request that asks the **Quote Bot** to return quotes for multiple ticker symbols.
_Avoid_: saved watchlist analysis, single-symbol lookup only

**Batch Size Limit**:
The maximum number of ticker symbols the **Quote Bot** accepts in one **Batch Quote Request**.
_Avoid_: unbounded symbol list

**Command Batch Interface**:
The rule that multi-symbol quote requests are accepted only through the explicit **Quote Command**.
_Avoid_: plain-text batch parser

**Per-Symbol Batch Result**:
The rule that each symbol in a **Batch Quote Request** returns its own success or failure outcome without failing the whole batch.
_Avoid_: all-or-nothing batch failure

**Input Order Reply**:
The rule that batch quote results are returned in the same symbol order the user submitted.
_Avoid_: sorted reply, grouped reply

**Normalized Batch Deduplication**:
The rule that repeated ticker symbols in a batch are collapsed after normalization while preserving first-seen order.
_Avoid_: repeated quote blocks, duplicate upstream lookup

**Lookup Normalization Pipeline**:
The ordered normalization steps applied before symbol deduplication, lookup, and reply formatting.
_Avoid_: inconsistent parsing rules across code paths

## Relationships

- A **Quote Bot** can share market data providers with the **Finance Watchlist Bot**
- A **Quote Bot** is distinct from a **Finance Watchlist Bot** because it answers ad hoc price lookups rather than saved-watchlist analysis
- A **Quote Bot** handles one **Quote Lookup** at a time
- A **Quote Lookup** requires an **Exact Ticker Symbol** in v1
- A **Lookup Trigger** in v1 can be either a bot command or a plain-text exact ticker symbol
- A plain-text **Lookup Trigger** can contain one **Ticker Token** with lightweight punctuation, but not mixed prose
- **Wrapper Punctuation** may be stripped from a **Ticker Token**, but `$` prefixes are not part of the accepted v1 syntax
- An **Invalid Symbol Response** tells the user to provide an **Exact Ticker Symbol**
- A **Provider Failure Response** tells the user to retry later without implying the symbol is wrong
- A successful **Quote Lookup** returns a **Canonical Symbol Reply**
- A successful **Canonical Symbol Reply** also includes the **Instrument Name**
- v1 is a **Stateless Quote Bot**
- The **Quote Bot** uses a **Dedicated Bot Identity**
- The canonical **Quote Command** in v1 is `/q`
- A **Supported Instrument** in v1 must be reported by the provider as `EQUITY` or `ETF`
- An **Unsupported Instrument Response** explains when a valid symbol is outside the supported instrument scope
- The **Quote Bot** runs inside the **Shared Gobot Runtime**
- v1 uses a **Generic Quote Interface**
- v1 uses an **English Quote UX**
- The **Quote Bot** uses a **Long Polling Runtime**
- The **Quote Bot** is branded as **SAAHAM Bot**
- The **Quote Bot** uses a **Branded Code Module**
- The **Chat Scope Rule** is: private chats accept commands and plain-text ticker tokens; group chats accept commands only
- Explicit group quote commands produce a **Public Group Reply**
- v1 supports a **Batch Quote Request**
- The **Batch Size Limit** in v1 is 5 symbols per request
- The **Command Batch Interface** means plain text stays single-symbol only while `/q` accepts batches
- A **Batch Quote Request** returns a **Per-Symbol Batch Result**
- A batch reply preserves **Input Order Reply**
- A batch request applies **Normalized Batch Deduplication**
- **Normalized Batch Deduplication** uses the **Lookup Normalization Pipeline**

## Example dialogue

> **Dev:** "Should this ticker lookup live in the **Finance Watchlist Bot**?"
> **Domain expert:** "No. The **Quote Bot** is a separate bot for quick price checks, while the **Finance Watchlist Bot** is for deeper watchlist analysis."

## Flagged ambiguities

- "new telegram bot" was ambiguous between extending the **Finance Watchlist Bot** and creating a separate **Quote Bot** — resolved: create a separate **Quote Bot**
- "check stock/etf ticker price" was ambiguous between exact lookup and symbol discovery — resolved: v1 uses **Exact Ticker Symbol** only; fuzzy lookup is deferred
- "how users ask for a quote" was ambiguous between command and free text — resolved: v1 accepts both bot commands and plain-text exact ticker symbols
- "plain text" was ambiguous between exact-token-only and mixed text parsing — resolved: v1 accepts one **Ticker Token** with lightweight punctuation, but not mixed prose
- "lightweight punctuation" was ambiguous about social-symbol syntax — resolved: strip **Wrapper Punctuation**, but reject `$`-prefixed symbols
- "lookup failure" was ambiguous between bad input and upstream outage — resolved: use separate **Invalid Symbol Response** and **Provider Failure Response**
- "reply ticker text" was ambiguous between user input and normalized symbol — resolved: successful lookups use a **Canonical Symbol Reply**
- "success payload" was ambiguous about identity fields — resolved: successful quote replies include both canonical symbol and **Instrument Name**
- "new bot scope" was ambiguous about storage and user state — resolved: v1 is a **Stateless Quote Bot**
- "separate bot" was ambiguous between module split and Telegram identity split — resolved: the **Quote Bot** uses a **Dedicated Bot Identity**
- "explicit command name" was ambiguous across several finance terms — resolved: the canonical **Quote Command** is `/q`
- "stocks and ETFs" was ambiguous between product language and runtime enforcement — resolved: a **Supported Instrument** must be provider-typed as `EQUITY` or `ETF`
- "valid but unsupported symbol" was ambiguous with bad ticker input — resolved: use a separate **Unsupported Instrument Response**
- "separate bot" was ambiguous between product identity and deployable shape — resolved: the **Quote Bot** runs in the **Shared Gobot Runtime**
- "command surface" was ambiguous between generic lookup and curated shortcuts — resolved: v1 uses a **Generic Quote Interface**
- "bot language" was ambiguous across existing repo patterns — resolved: v1 uses an **English Quote UX**
- "Telegram transport" was ambiguous between current repo conventions and extra infra — resolved: the **Quote Bot** uses a **Long Polling Runtime**
- "`SAAHAM_BOT`" was ambiguous between typo and branding — resolved: **SAAHAM Bot** is the intentional brand spelling
- "where branding applies" was ambiguous between product/config only and code structure too — resolved: the **Quote Bot** uses a **Branded Code Module**
- "where plain-text lookup is allowed" was ambiguous across chat types — resolved: follow the **Chat Scope Rule**
- "group reply visibility" was ambiguous after enabling group commands — resolved: explicit group lookups produce a **Public Group Reply**
- "one vs many tickers per request" was ambiguous — resolved: v1 supports a **Batch Quote Request**
- "how many tickers count as a quote request" was ambiguous — resolved: the **Batch Size Limit** is 5
- "whether plain text can batch" was ambiguous after enabling multiple tickers — resolved: use the **Command Batch Interface**
- "how batch failures behave" was ambiguous — resolved: a **Batch Quote Request** returns a **Per-Symbol Batch Result**
- "batch output ordering" was ambiguous — resolved: preserve **Input Order Reply**
- "duplicate symbols in batch input" was ambiguous — resolved: apply **Normalized Batch Deduplication**
- "what counts as the same symbol" was ambiguous — resolved: deduplicate using the **Lookup Normalization Pipeline**
