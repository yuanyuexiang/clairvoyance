# clairvoyance
clairvoyance demonstrates how to send video and/or audio to your browser from files saved to disk.

For an example of playing H264 from disk see [clairvoyance-h264](https://github.com/pion/example-webrtc-applications/tree/master/clairvoyance-h264)

## Instructions
### Create IVF named `output.ivf` that contains a VP8/VP9/AV1 track and/or `output.ogg` that contains a Opus track
```
ffmpeg -i $INPUT_FILE -g 30 -b:v 2M output.ivf
ffmpeg -i $INPUT_FILE -c:a libopus -page_duration 20000 -vn output.ogg
```

**Note**: In the `ffmpeg` command which produces the .ivf file, the argument `-b:v 2M` specifies the video bitrate to be 2 megabits per second. We provide this default value to produce decent video quality, but if you experience problems with this configuration (such as dropped frames etc.), you can decrease this. See the [ffmpeg documentation](https://ffmpeg.org/ffmpeg.html#Options) for more information on the format of the value.

### Download clairvoyance

```
go install github.com/pion/webrtc/v4/examples/clairvoyance@latest
```

### Open clairvoyance example page
[jsfiddle.net](https://jsfiddle.net/8kup9mvn/) you should see two text-areas, 'Start Session' button and 'Copy browser SessionDescription to clipboard'

### Run clairvoyance with your browsers Session Description as stdin
The `output.ivf` you created should be in the same directory as `clairvoyance`. In the jsfiddle press 'Copy browser Session Description to clipboard' or copy the base64 string manually.

Now use this value you just copied as the input to `clairvoyance`

#### Linux/macOS
Run `echo $BROWSER_SDP | clairvoyance`
#### Windows
1. Paste the SessionDescription into a file.
1. Run `clairvoyance < my_file`

### Input clairvoyance's Session Description into your browser
Copy the text that `clairvoyance` just emitted and copy into the second text area in the jsfiddle

### Hit 'Start Session' in jsfiddle, enjoy your video!
A video should start playing in your browser above the input boxes. `clairvoyance` will exit when the file reaches the end

Congrats, you have used Pion WebRTC! Now start building something cool

https://docs.github.com/zh/actions/writing-workflows/workflow-syntax-for-github-actions#jobsjob_idcontainer

git config --global http.proxy http://127.0.0.1:7890
git config --global https.proxy http://127.0.0.1:7890

git config --global --unset http.proxy
git config --global --unset https.proxy
