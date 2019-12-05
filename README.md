# video-picker
自动从虎牙下载视频,并剪辑精彩片段

## 原理
docker run -d --rm --name conv -v /home/tangxu/opensource/huya-live-recorder:/hecate/examples/  creactiviti/hecate hecate -i /hecate/examples/test.mp4 -o /hecate/examples/output/ --generate_mov --lmov 60 -w 1080
docker logs -f conv


## 参照
https://github.com/yahoo/hecate

https://github.com/leeeboo/huya-stream/blob/master/main.go

>需要知道视频时间,否则使用mp3的时间作为结束
ffmpeg -i /home/tangxu/test.mp3 -i /home/tangxu/dota1_sum.mp4  -t 7.1 -y /home/tangxu/new1.mp4

>这样又没有配音了...
ffmpeg -i /home/tangxu/dota1_sum.mp4 -i /home/tangxu/test.mp3  -y /home/tangxu/new1.mp4
