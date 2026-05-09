# Gobot

Gobot is a shared Go runtime for multiple Telegram bots and scouts.
Each module focuses on a distinct domain, while the repo reuses common integrations and deployment shape.

## User-facing bots

[Finance :simple-telegram:](finance.md){ .md-button }

[SAAHAM Bot :simple-telegram:](saaham.md){ .md-button }

[Spotifytube :simple-spotify: :simple-youtubemusic:](spotifytube.md){ .md-button }

[Free Games on :simple-steam::simple-epicgames:](freegames.md){ .md-button }

## Runtime modules in this repo

- `finance`: configurable watchlist analysis and scheduled scout alerts
- `saaham`: exact-symbol stock and ETF quote lookup bot
- `spotifytube`: Spotify and YouTube link conversion bot
- `freegames`: Reddit-driven free game scout for Telegram channels
- `antam`: gold price bot and scheduled gold price sender
- `reddit`: Reddit media scout that posts curated content to Telegram
