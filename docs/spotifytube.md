# Spotify Tube Bot

```mermaid
graph LR
    A --> S
    A --> Sp
    I --> Se
    Se --> U
    subgraph S[handleYoutubeURL]
        I
        Se
        U
    end

    subgraph Sp[handleSpotifyURL]
        E --> G
        G --> T
        T --> SY
        SY --> Su
    end

    A[handleTextMessage]
    E[ExtractSpotifyTrackID]
    G[GetSpotifyAccessToken]
    T[GetSpotifyTrack]
    SY[SearchYoutube]
    Su[sendYoutubeURLs]
```