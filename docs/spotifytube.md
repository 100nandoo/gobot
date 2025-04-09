# Spotify Tube Bot

```mermaid
graph LR
    A --> S
    A --> Y

    S --> I
    Y --> I
    I --> Se[Search]
    Se --> U

    A[handleTextMessage]
    S[handleSpotifyURL]
    Y[handleYoutubeURL]
    I[InspectUrl]
    Se[Search]
    U[sendTrackURLs]
```