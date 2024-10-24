package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"

	// If you don't like x264, you can also use vpx by importing as below
	// This is required to use VP8/VP9 video encoder
	// or you can also use openh264 for alternative h264 implementation
	// "github.com/pion/mediadevices/pkg/codec/openh264"
	// or if you use a raspberry pi like, you can use mmal for using its hardware encoder
	// "github.com/pion/mediadevices/pkg/codec/mmal"
	"github.com/pion/mediadevices/pkg/codec/openh264"
	"github.com/pion/mediadevices/pkg/codec/opus"

	// This is required to use opus audio encoder
	//"github.com/pion/mediadevices/pkg/codec/x264" // This is required to use h264 video encoder

	// Note: If you don't have a camera or microphone or your adapters are not supported,
	//       you can always swap your adapters with our dummy adapters below.
	// _ "github.com/pion/mediadevices/pkg/driver/audiotest"
	// _ "github.com/pion/mediadevices/pkg/driver/videotest"
	_ "github.com/pion/mediadevices/pkg/driver/camera"     // This is required to register camera adapter
	_ "github.com/pion/mediadevices/pkg/driver/microphone" // This is required to register microphone adapter
)

func main() {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	//signal.Decode(signal.MustReadStdin(), &offer)

	content, err := os.ReadFile("sdp.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("==================================", string(content))
	sdp_str := string(content)
	decode(sdp_str, &offer)
	//decode(readUntilNewline(), &offer)

	// Create a new RTCPeerConnection
	h264Params, err := openh264.NewParams()
	if err != nil {
		panic(err)
	}
	h264Params.BitRate = 1_000_000

	opusParams, err := opus.NewParams()

	if err != nil {
		panic(err)
	}
	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&h264Params),
		mediadevices.WithAudioEncoders(&opusParams),
	)

	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	s, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			//c.FrameFormat = prop.FrameFormat(frame.FormatI420)
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
		},
		Audio: func(c *mediadevices.MediaTrackConstraints) {
		},
		Codec: codecSelector,
	})
	if err != nil {
		panic(err)
	}

	for _, track := range s.GetTracks() {
		track.OnEnded(func(err error) {
			fmt.Printf("Track (ID: %s) ended with error: %v\n",
				track.ID(), err)
		})

		_, err = peerConnection.AddTransceiverFromTrack(track,
			webrtc.RtpTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionSendonly,
			},
		)
		if err != nil {
			panic(err)
		}
	}

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(encode(peerConnection.LocalDescription()))

	// Block forever
	select {}
}

// Decode a base64 and unmarshal JSON into a SessionDescription
func decode(in string, obj *webrtc.SessionDescription) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(b, obj); err != nil {
		panic(err)
	}
}

// JSON encode + base64 a SessionDescription
func encode(obj *webrtc.SessionDescription) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}
