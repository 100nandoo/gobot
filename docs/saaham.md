# SAAHAM Bot
![image](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)
![image](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)

SAAHAM Bot is a fast quote bot for Yahoo Finance ticker symbols, currency pairs, and supported market shortcuts.
It is separate from the finance watchlist bot: SAAHAM does ad hoc quote lookups, while `finance` does saved-watchlist analysis.

## What it does

- Returns the latest quote for one stock, ETF, index, or currency pair
- Accepts `/q` command lookups in private chats and groups
- Accepts plain-text ticker lookups in private chats and groups
- Supports batch lookups through `/q` for up to `5` symbols
- Normalizes ticker input and deduplicates repeated symbols in batch requests after canonical symbol resolution
- Expands supported shortcuts such as `STI -> ^STI`, `JKSE -> ^JKSE`, `IHSG -> ^JKSE`, `VWRA -> VWRA.L`, `CSPX -> CSPX.L`, `BJBR -> BJBR.JK`, and `USDSGD -> USDSGD=X`
- Uses a checked-in manual S&P 100 snapshot so symbols such as `AAPL` keep bare-symbol-first lookup

## Commands

| Command | Description |
|---|---|
| `/start` | Show help |
| `/help` | Show help |
| `/q AAPL` | Get the latest quote for one symbol |
| `/q STI` | Resolve a supported shortcut such as `STI -> ^STI` |
| `/q USDSGD` | Resolve a currency pair shortcut such as `USDSGD -> USDSGD=X` |
| `/q AAPL MSFT NVDA` | Get quotes for multiple symbols in one request |

## Chat behavior

- Private chats accept `/q` and plain-text ticker inputs such as `AAPL`, `STI`, `CSPX`, or `USDSGD`
- Group chats accept `/q` and a single plain-text ticker input
- Plain-text mixed prose is ignored
- `$AAPL` is rejected; use ticker inputs without the `$` prefix
- Wrapper punctuation such as `(AAPL)` is stripped before lookup
- Explicit aliases such as `IHSG -> ^JKSE` run before generic suffixless lookup rules
- Inputs that already contain `^`, `.`, or `=` skip shortcut expansion
- Six-letter currency pair shortcuts such as `USDSGD` try `USDSGD=X` before other shortcut candidates
- S&P 100 snapshot members such as `AAPL` try the bare symbol first, then `^SYMBOL`, then `.L`, then `.JK`
- Other suffixless inputs use shortcut-first lookup order: `^SYMBOL`, then `.L`, then `.JK`, then bare `SYMBOL`

## Supported instruments

SAAHAM Bot currently supports provider-classified:

- `EQUITY`
- `ETF`
- `INDEX`
- `CURRENCY`

If the symbol is valid but outside that scope, the bot replies that the instrument is unsupported.

## Failure handling

Single-symbol lookups use three user-facing outcomes:

- invalid symbol
- unsupported instrument
- provider failure

Batch lookups return per-symbol results, so one failure does not block the rest of the request.

## Runtime notes

- SAAHAM runs inside the shared Gobot process
- It uses long polling through `telebot.v3`
- It is stateless: there is no watchlist, settings store, or lookup history
- Quote replies show the canonical symbol, instrument name, current price, currency, and daily change
- The S&P 100 snapshot is a manual checked-in file; runtime does not fetch constituents live

## Setup

### Environment Variables

| Name | Desc |
|---|---|
| `SAAHAM_BOT` | Telegram bot API token |

## Related docs

- [Finance](finance.md)
- [Architecture decision: separate quote bot in shared runtime](adr/0001-separate-quote-bot-in-shared-runtime.md)
