# Gobot

Gobot is a shared Go runtime for several Telegram bots and background scouts.
Each bot has a distinct user-facing scope, while the process reuses common integrations such as Telegram, Yahoo Finance, Reddit, Supabase, Spotify, and YouTube.

## Bots in this repo

- `finance`: configurable watchlist analysis bot for stocks and ETFs
- `saaham`: stateless quote bot for exact-symbol stock and ETF lookups
- `freegames`: scout that posts free Steam and Epic finds to Telegram
- `spotifytube`: converts between Spotify tracks and YouTube links
- `antam`: gold price bot for Antam and Pluang pricing
- `reddit`: Reddit media scout that posts curated content to Telegram

## Runtime shape

The entrypoint is [`main.go`](./main.go).

- Background scouts start for `freegames`, `antam`, `reddit`, and `finance`
- Interactive Telegram bots start for `antam`, `spotifytube`, `finance`, and `saaham`
- Shared helpers live under [`pkg`](./pkg) and configuration keys live in [`config/config.go`](./config/config.go)

## Documentation

- Project docs site: <https://100nandoo.github.io/gobot/>
- Docs index: [`docs/index.md`](./docs/index.md)
- Domain language: [`CONTEXT.md`](./CONTEXT.md)
- Architecture decision for the quote bot: [`docs/adr/0001-separate-quote-bot-in-shared-runtime.md`](./docs/adr/0001-separate-quote-bot-in-shared-runtime.md)

## Bot docs

- [`docs/finance.md`](./docs/finance.md)
- [`docs/freegames.md`](./docs/freegames.md)
- [`docs/spotifytube.md`](./docs/spotifytube.md)
- [`docs/saaham.md`](./docs/saaham.md)

## Configuration

Environment variable names are defined in [`config/config.go`](./config/config.go).

Common keys used across the repo include:

- `SUPABASE_URL`
- `SUPABASE_KEY`
- `TELEGRAM_BOT`
- `ANTAM_TELEGRAM_BOT`
- `SPOTIFYTUBE_BOT`
- `SAAHAM_BOT`
- `TELEGRAM_CHANNEL_FREE_GAMES`
- `TELEGRAM_CHANNEL_ANTAM`
- `TELEGRAM_CHANNEL_REDDIT`
- `TELEGRAM_CHANNEL_FINANCE`
- `SPOTIFY_ID`
- `SPOTIFY_SECRET`
- `YOUTUBE_TOKEN`

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=100nandoo/gobot&type=Date)](https://www.star-history.com/#100nandoo/gobot&Date)
