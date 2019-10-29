# ffmpeg-k8s-dns-repro

TL;DR: Run `kubectl apply -f https://raw.githubusercontent.com/jeremyje/ffmpeg-k8s-dns-repro/master/deploy.yaml` and look at the containers for failures.

For some reason Kubernetes within Minikube does not like localhost DNS resolves. Here's a demo showing that.

1. Install VirtualBox from https://www.virtualbox.org/wiki/Downloads, I'm using `6.0.14`.
1. Download minikube from https://github.com/kubernetes/minikube/releases/tag/v1.5.0.

Open a shell and run the following commands.

```bash
$ minikube start
�😄  minikube v1.5.0 on Microsoft Windows 10 Enterprise 10.0.18362 Build 18362
✨  Automatically selected the 'virtualbox' driver
�💿  Downloading VM boot image .
    > minikube-v1.5.0.iso.sha256: 65 B / 65 B [--------------] 100.00% ? p/s 0s
    > minikube-v1.5.0.iso: 143.77 MiB / 143.77 MiB [] 100.00% 11.29 MiB p/s 13s
�🔥  Creating virtualbox VM (CPUs=2, Memory=2000MB, Disk=20000MB) .
�🐳  Preparing Kubernetes v1.16.2 on Docker 18.09.9 .
�💾  Downloading kubeadm v1.16
�💾  Downloading kubelet v1.16
�🚜  Pulling images .
�🚀  Launching Kubernetes ..
⌛  Waiting for: apiserver proxy etcd scheduler controller dns
�🏄  Done! kubectl is now configured to use "minikube"
$ kubectl apply -f https://raw.githubusercontent.com/jeremyje/ffmpeg-k8s-dns-repro/master/deploy.yaml
job.batch/ffmpeg-broken-dns created
$ minikube dashboard
�🤔  Verifying dashboard health ...
�🚀  Launching proxy ...
�🔥  Verifying proxy health ...
�💔  Opening http://127.0.0.1:65114/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/ in your default browser...
```

Lookup the `ffmpeg-broken-dns*` pod. There's multiple containers in it you'll notice that the containers ending in `connect-localhost` are failing while `connect-127001` are successful. It looks like telling ffmpeg to connect to `tcp://localhost:$PORT` is broken.

```bash
$ nslookup localhost
nslookup: can't resolve '(null)': Name does not resolve

Name:      localhost
Address 1: 127.0.0.1 localhost
Address 2: ::1 localhost
```

There's a sleepy container you can use to shell into and run the command manually. Each pod is running a variant of:

`/app/main -serve_via=127.0.0.1 -connect_via=localhost -port=12346`

within the `/app/` directory.


The failing pods have output that looks like this.

```bash

2019/10/29 08:30:47 Running -serve_via=localhost -connect_via=localhost -port=10000
2019/10/29 08:30:47 Serving via localhost:10000
2019/10/29 08:30:47 wait 5 seconds to make sure TCP is available
2019/10/29 08:30:52 calling ffmpeg with [-y -progress tcp://localhost:10000 -i /app/small.ogv -strict -2 /app/small.mp4]
2019/10/29 08:30:52 ffmpeg version N-50535-gec5d385722-static https://johnvansickle.com/ffmpeg/  Copyright (c) 2000-2019 the FFmpeg developers
  built with gcc 6.3.0 (Debian 6.3.0-18+deb9u1) 20170516
  configuration: --enable-gpl --enable-version3 --enable-static --disable-debug --disable-ffplay --disable-indev=sndio --disable-outdev=sndio --cc=gcc-6 --enable-fontconfig --enable-frei0r --enable-gnutls --enable-gmp --enable-libgme --enable-gray --enable-libaom --enable-libfribidi --enable-libass --enable-libvmaf --enable-libfreetype --enable-libmp3lame --enable-libopencore-amrnb --enable-libopencore-amrwb --enable-libopenjpeg --enable-librubberband --enable-libsoxr --enable-libspeex --enable-libsrt --enable-libvorbis --enable-libopus --enable-libtheora --enable-libvidstab --enable-libvo-amrwbenc --enable-libvpx --enable-libwebp --enable-libx264 --enable-libx265 --enable-libxml2 --enable-libdav1d --enable-libxvid --enable-libzvbi --enable-libzimg
  libavutil      56. 35.101 / 56. 35.101
  libavcodec     58. 59.102 / 58. 59.102
  libavformat    58. 33.100 / 58. 33.100
  libavdevice    58.  9.100 / 58.  9.100
  libavfilter     7. 64.100 /  7. 64.100
  libswscale      5.  6.100 /  5.  6.100
  libswresample   3.  6.100 /  3.  6.100
  libpostproc    55.  6.100 / 55.  6.100
[tcp @ 0x55555596bd00] Failed to resolve hostname localhost: Name or service not known
Failed to open progress URL "tcp://localhost:10000": Input/output error
Failed to set value 'tcp://localhost:10000' for option 'progress': Input/output error
Error parsing global options: Input/output error
2019/10/29 08:30:52 
ffmpeg has closed.
2019/10/29 08:30:52 cannot call ffmpeg exit status 1
```

Successful containers have output that looks like this:

```bash
2019/10/29 08:30:48 Running -serve_via=localhost -connect_via=127.0.0.1 -port=10001
2019/10/29 08:30:48 Serving via localhost:10001
2019/10/29 08:30:48 wait 5 seconds to make sure TCP is available
2019/10/29 08:30:53 calling ffmpeg with [-y -progress tcp://127.0.0.1:10001 -i /app/small.ogv -strict -2 /app/small.mp4]
2019/10/29 08:30:54 ffmpeg progress: frame=142
2019/10/29 08:30:54 ffmpeg progress: fps=0.00
2019/10/29 08:30:54 ffmpeg progress: stream_0_0_q=29.0
2019/10/29 08:30:54 ffmpeg progress: bitrate=   0.1kbits/s
2019/10/29 08:30:54 ffmpeg progress: total_size=48
2019/10/29 08:30:54 ffmpeg progress: out_time_us=4693333
2019/10/29 08:30:54 ffmpeg progress: out_time_ms=4693333
2019/10/29 08:30:54 ffmpeg progress: out_time=00:00:04.693333
2019/10/29 08:30:54 ffmpeg progress: dup_frames=0
2019/10/29 08:30:54 ffmpeg progress: drop_frames=0
2019/10/29 08:30:54 ffmpeg progress: speed=9.38x
2019/10/29 08:30:54 ffmpeg progress: progress=continue
2019/10/29 08:30:54 ffmpeg progress: frame=166
2019/10/29 08:30:54 ffmpeg progress: fps=0.00
2019/10/29 08:30:54 ffmpeg progress: stream_0_0_q=-1.0
2019/10/29 08:30:54 ffmpeg progress: bitrate= 290.7kbits/s
2019/10/29 08:30:54 ffmpeg progress: total_size=202310
2019/10/29 08:30:54 ffmpeg progress: out_time_us=5568000
2019/10/29 08:30:54 ffmpeg progress: out_time_ms=5568000
2019/10/29 08:30:54 ffmpeg progress: out_time=00:00:05.568000
2019/10/29 08:30:54 ffmpeg progress: dup_frames=0
2019/10/29 08:30:54 ffmpeg progress: drop_frames=0
2019/10/29 08:30:54 ffmpeg progress: speed=7.32x
2019/10/29 08:30:54 ffmpeg progress: progress=end
2019/10/29 08:30:54 ffmpeg version N-50535-gec5d385722-static https://johnvansickle.com/ffmpeg/  Copyright (c) 2000-2019 the FFmpeg developers
  built with gcc 6.3.0 (Debian 6.3.0-18+deb9u1) 20170516
  configuration: --enable-gpl --enable-version3 --enable-static --disable-debug --disable-ffplay --disable-indev=sndio --disable-outdev=sndio --cc=gcc-6 --enable-fontconfig --enable-frei0r --enable-gnutls --enable-gmp --enable-libgme --enable-gray --enable-libaom --enable-libfribidi --enable-libass --enable-libvmaf --enable-libfreetype --enable-libmp3lame --enable-libopencore-amrnb --enable-libopencore-amrwb --enable-libopenjpeg --enable-librubberband --enable-libsoxr --enable-libspeex --enable-libsrt --enable-libvorbis --enable-libopus --enable-libtheora --enable-libvidstab --enable-libvo-amrwbenc --enable-libvpx --enable-libwebp --enable-libx264 --enable-libx265 --enable-libxml2 --enable-libdav1d --enable-libxvid --enable-libzvbi --enable-libzimg
  libavutil      56. 35.101 / 56. 35.101
  libavcodec     58. 59.102 / 58. 59.102
  libavformat    58. 33.100 / 58. 33.100
  libavdevice    58.  9.100 / 58.  9.100
  libavfilter     7. 64.100 /  7. 64.100
  libswscale      5.  6.100 /  5.  6.100
  libswresample   3.  6.100 /  3.  6.100
  libpostproc    55.  6.100 / 55.  6.100
[ogg @ 0x5555556f1000] Broken file, keyframe not correctly marked.
Input #0, ogg, from '/app/small.ogv':
  Duration: 00:00:05.55, start: 0.000000, bitrate: 632 kb/s
    Stream #0:0: Data: none
    Stream #0:1: Video: theora, yuv420p, 560x320, 30 fps, 30 tbr, 30 tbn, 30 tbc
    Metadata:
      ENCODER         : ffmpeg2theora-0.26
      SOURCE_OSHASH   : d1af78a82e61d18f
    Stream #0:2: Audio: vorbis, 48000 Hz, mono, fltp, 110 kb/s
    Metadata:
      ENCODER         : ffmpeg2theora-0.26
      SOURCE_OSHASH   : d1af78a82e61d18f
Stream mapping:
  Stream #0:1 -> #0:0 (theora (native) -> h264 (libx264))
  Stream #0:2 -> #0:1 (vorbis (native) -> aac (native))
Press [q] to stop, [?] for help
[libx264 @ 0x5555556f7a00] using cpu capabilities: MMX2 SSE2Fast LZCNT SSSE3 SSE4.2 AVX
[libx264 @ 0x5555556f7a00] profile Progressive High, level 3.0, 4:2:0, 8-bit
[libx264 @ 0x5555556f7a00] 264 - core 157 r2969 d4099dd - H.264/MPEG-4 AVC codec - Copyleft 2003-2019 - http://www.videolan.org/x264.html - options: cabac=1 ref=3 deblock=1:0:0 analyse=0x3:0x113 me=hex subme=7 psy=1 psy_rd=1.00:0.00 mixed_ref=1 me_range=16 chroma_me=1 trellis=1 8x8dct=1 cqm=0 deadzone=21,11 fast_pskip=1 chroma_qp_offset=-2 threads=3 lookahead_threads=1 sliced_threads=0 nr=0 decimate=1 interlaced=0 bluray_compat=0 constrained_intra=0 bframes=3 b_pyramid=2 b_adapt=1 b_bias=0 direct=1 weightb=1 open_gop=0 weightp=2 keyint=250 keyint_min=25 scenecut=40 intra_refresh=0 rc_lookahead=40 rc=crf mbtree=1 crf=23.0 qcomp=0.60 qpmin=0 qpmax=69 qpstep=4 ip_ratio=1.40 aq=1:1.00
Output #0, mp4, to '/app/small.mp4':
  Metadata:
    encoder         : Lavf58.33.100
    Stream #0:0: Video: h264 (libx264) (avc1 / 0x31637661), yuv420p(progressive), 560x320, q=-1--1, 30 fps, 15360 tbn, 30 tbc
    Metadata:
      SOURCE_OSHASH   : d1af78a82e61d18f
      encoder         : Lavc58.59.102 libx264
    Side data:
      cpb: bitrate max/min/avg: 0/0/0 buffer size: 0 vbv_delay: N/A
    Stream #0:1: Audio: aac (LC) (mp4a / 0x6134706D), 48000 Hz, mono, fltp, 69 kb/s
    Metadata:
      SOURCE_OSHASH   : d1af78a82e61d18f
      encoder         : Lavc58.59.102 aac
[ogg @ 0x5555556f1000] Broken file, keyframe not correctly marked.
    Last message repeated 1 times
frame=  142 fps=0.0 q=29.0 size=       0kB time=00:00:04.69 bitrate=   0.1kbits/s speed=9.38x    
frame=  166 fps=0.0 q=-1.0 Lsize=     198kB time=00:00:05.56 bitrate= 290.7kbits/s speed=7.32x    
video:143kB audio:47kB subtitle:0kB other streams:0kB global headers:0kB muxing overhead: 3.800884%
[libx264 @ 0x5555556f7a00] frame I:1     Avg QP:18.15  size: 18737
[libx264 @ 0x5555556f7a00] frame P:62    Avg QP:21.48  size:  1580
[libx264 @ 0x5555556f7a00] frame B:103   Avg QP:26.33  size:   283
[libx264 @ 0x5555556f7a00] consecutive B-frames: 13.9%  7.2%  9.0% 69.9%
[libx264 @ 0x5555556f7a00] mb I  I16..4: 38.4% 46.0% 15.6%
[libx264 @ 0x5555556f7a00] mb P  I16..4:  1.1%  1.4%  0.1%  P16..4: 15.7%  7.2%  6.1%  0.0%  0.0%    skip:68.4%
[libx264 @ 0x5555556f7a00] mb B  I16..4:  0.0%  0.2%  0.0%  B16..8: 14.5%  1.5%  0.2%  direct: 0.3%  skip:83.1%  L0:47.6% L1:46.3% BI: 6.1%
[libx264 @ 0x5555556f7a00] 8x8 transform intra:53.9% inter:65.9%
[libx264 @ 0x5555556f7a00] coded y,uvDC,uvAC intra: 45.5% 73.8% 31.4% inter: 3.8% 6.7% 1.9%
[libx264 @ 0x5555556f7a00] i16 v,h,dc,p:  1% 42%  4% 54%
[libx264 @ 0x5555556f7a00] i8 v,h,dc,ddl,ddr,vr,hd,vl,hu:  9% 40% 29%  2%  4%  2%  7%  2%  5%
[libx264 @ 0x5555556f7a00] i4 v,h,dc,ddl,ddr,vr,hd,vl,hu: 15% 26% 13%  5%  8%  8% 12%  4%  9%
[libx264 @ 0x5555556f7a00] i8c dc,h,v,p: 36% 46% 11%  7%
[libx264 @ 0x5555556f7a00] Weighted P-Frames: Y:0.0% UV:0.0%
[libx264 @ 0x5555556f7a00] ref P L0: 72.2%  8.4% 11.9%  7.5%
[libx264 @ 0x5555556f7a00] ref B L0: 82.2% 12.9%  4.9%
[libx264 @ 0x5555556f7a00] ref B L1: 95.9%  4.1%
[libx264 @ 0x5555556f7a00] kb/s:210.79
[aac @ 0x55555589a640] Qavg: 366.779
2019/10/29 08:30:54 
ffmpeg has closed.
2019/10/29 08:30:54 error handling TCP accept tcp 127.0.0.1:10001: use of closed network connection, ignoring since this happens on ffmpeg close
```