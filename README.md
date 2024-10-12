## MBP CAM

https://stackoverflow.com/questions/26999595/what-steps-are-needed-to-stream-rtsp-from-ffmpeg

```bash
ffmpeg -re -stream_loop 0 -i temporal_assets/Tobilleras.mp4 -vcodec libx264 -f rtsp -rtsp_transport tcp rtsp://localhost:8554/live.stream
```

## Diagrams

https://diagrams.mingrammer.com/

https://plantuml.com/

https://structurizr.com/

## RTSP Allowed Commands

https://www.wowza.com/blog/rtsp-the-real-time-streaming-protocol-explained -> RTSP Commands

## RTP Documentation

https://pion.ly/why-pion/