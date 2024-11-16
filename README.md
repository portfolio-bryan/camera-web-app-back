## MBP CAM

https://stackoverflow.com/questions/26999595/what-steps-are-needed-to-stream-rtsp-from-ffmpeg

```bash
ffmpeg -re -stream_loop 0 -i temporal_assets/Tobilleras.mp4 -vcodec libx264 -f rtsp -rtsp_transport tcp rtsp://localhost:8554/live.stream
```

## FFPROBE

```bash
ffprobe -v error <video> -show_format -show_streams -print_format json -select_streams v -show_entries stream=codec_name
```

To find [open videos](https://test-videos.co.uk/)

## FFPLAY

```bash
ffplay -v error <video> -x 600 -y 600 -noborder -top 0 -left 0 -fs -an -showmode waves
```

fs: fullscreen
an: with no sound
vn: with no video

## Diagrams

https://diagrams.mingrammer.com/

https://plantuml.com/

https://structurizr.com/

## RTSP Allowed Commands

https://www.wowza.com/blog/rtsp-the-real-time-streaming-protocol-explained -> RTSP Commands

## RTP Documentation

https://pion.ly/why-pion/




https://web.dev/articles/webrtc-infrastructure?utm_source=substack&utm_medium=email
https://webrtc.org/?utm_source=substack&utm_medium=email