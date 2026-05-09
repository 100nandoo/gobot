# Separate quote bot with dedicated Telegram identity inside the shared Gobot runtime

We will build the Quote Bot as a separate Telegram bot product with its own token and user-facing command surface, rather than extending the existing Finance Watchlist Bot. The Quote Bot will still run as a module inside the shared `gobot` process because that matches the current runtime architecture and keeps deployment simple. In v1, it is a stateless English-language bot for exact-symbol quote lookups only, and it will only serve provider-classified `EQUITY` and `ETF` instruments.
