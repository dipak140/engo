# Engo Transcoder
Origin Server for Transcoding + HLS Packaging

## Build and run
``` bash
go build -o engo .
./engo videos/input.mp4
```

## Low Level Design
```
             ┌────────────────────┐
             │   Blob Storage     │  ← Input video
             └────────┬───────────┘
                      ↓
            ┌─────────────────────┐
            │   Origin Server     │
            │  (Transcoding Core) │
            └────────┬────────────┘
         ┌───────────┴──────────────┐
         ↓                          ↓
 [ Transcoding Manager ]      [ HLS Packager ]
         ↓                          ↓
 ┌──────────────┐          ┌─────────────────────┐
 │ Temp Storage │ ←──────→ │ Playlist Generator  │
 └──────────────┘          └────────┬────────────┘
                                    ↓
                           [ Output Uploader / HTTP Serve ]
```


## Core Components
1. Blob Fetcher: Fetch video file from blob storage, also has validation for video file. 
2. Transcoding Manager: Transcode video into multi bitrate hls chunks. Uses worker pool pattern for efficiency.
3. HLS Packager: 
    1. create folder for output videos and stores locally 
    2. Generates master.m3u8 pointing to:
        1080p/index.m3u8
        720p/index.m3u8
4. Output Uploader: Upload to blob

## Future Enhancements
1. Job Queue (Redis): Handle concurrent jobs, retries
2. Webhook: on Complete Notify client when HLS is ready
3. Presigned URLs: Secure temporary access to HLS files
4. Database: Store job metadata, status tracking
5. Health Checks: Production-grade deployment readiness
