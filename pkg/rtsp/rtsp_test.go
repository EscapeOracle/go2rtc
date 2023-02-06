package rtsp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLParse(t *testing.T) {
	// https://github.com/AlexxIT/WebRTC/issues/395
	base := "rtsp://::ffff:192.168.1.123/onvif/profile.1/"
	u, err := urlParse(base)
	assert.Empty(t, err)
	assert.Equal(t, "::ffff:192.168.1.123:", u.Host)

	// https://github.com/AlexxIT/go2rtc/issues/208
	base = "rtsp://rtsp://turret2-cam.lan:554/stream1/"
	u, err = urlParse(base)
	assert.Empty(t, err)
	assert.Equal(t, "turret2-cam.lan:554", u.Host)
}

func TestMultipleSinSDP(t *testing.T) {
	// https://github.com/AlexxIT/WebRTC/issues/417
	s := `v=0
o=- 91674849066 1 IN IP4 192.168.1.123
s=RtspServer
i=live
t=0 0
a=control:*
a=range:npt=0-
m=video 0 RTP/AVP 96
c=IN IP4 0.0.0.0
s=RtspServer
i=live
a=control:track0
a=range:npt=0-
a=rtpmap:96 H264/90000
a=fmtp:96 packetization-mode=1;profile-level-id=42001E;sprop-parameter-sets=Z0IAHvQCgC3I,aM48gA==
a=control:track0
m=audio 0 RTP/AVP 97
c=IN IP4 0.0.0.0
s=RtspServer
i=live
a=control:track1
a=range:npt=0-
a=rtpmap:97 MPEG4-GENERIC/8000/1
a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3;config=1588
a=control:track1
`
	medias, err := UnmarshalSDP([]byte(s))
	assert.Nil(t, err)
	assert.NotNil(t, medias)
}

func TestFindFmtp(t *testing.T) {
	// https://github.com/AlexxIT/WebRTC/issues/419
	s := `v=0
o=- 1675628282 1675628283 IN IP4 192.168.1.123
s=streamed by the RTSP server
t=0 0
m=video 0 RTP/AVP 96
a=rtpmap:96 H264/90000
a=control:track0
m=audio 0 RTP/AVP 8
a=rtpmap:0 pcma/8000/1
a=control:track1
a=framerate:25
a=range:npt=now-
a=fmtp:96 packetization-mode=1;profile-level-id=64001F;sprop-parameter-sets=Z0IAH5WoFAFuQA==,aM48gA==
`
	medias, err := UnmarshalSDP([]byte(s))
	assert.Nil(t, err)
	assert.NotNil(t, medias)
	assert.NotEqual(t, "", medias[0].Codecs[0].FmtpLine)
}
