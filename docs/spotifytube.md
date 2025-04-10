# Spotifytube
Function Chain Call Chart for Spotifytube
```mermaid
graph LR
    A --> S
    A --> Y
    subgraph Y[handleYoutubeURL]
        Ey --> V
        V --> Ss
        Ss --> Ssu
    end

    subgraph S[handleSpotifyURL]
        E --> T
        T --> SY
        SY --> Syu
    end

    A[handleTextMessage]
    E[ExtractSpotifyTrackID]
    T[GetSpotifyTrack]
    SY[SearchYoutube]
    Syu[sendYoutubeURLs]

    Ey[ExtractYoutubeVideoID]
    V[GetVideo]
    Ss[SearchSpotify]
    Ssu[sendSpotifyURLs]
```