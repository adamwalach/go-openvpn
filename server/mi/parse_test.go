package mi_test

import (
	"testing"

	mi "github.com/adamwalach/go-openvpn/server/mi"
	"github.com/stretchr/testify/assert"
)

var responseVersion = `OpenVPN Version: OpenVPN 2.3.2 x86_64-pc-linux-gnu [SSL (OpenSSL)] [LZO] [EPOLL] [PKCS11] [eurephia] [MH] [IPv6] built on Dec  1 2014
Management Version: 1
END
`

var responsePid = `SUCCESS: pid=10869
`

var responseError = `Error: bad response
`
var responseLoadStats = `SUCCESS: nclients=12,bytesin=109984474,bytesout=1589364037
`
var responseLoadStatsBroken = `SUCCESS: nclients=12,bytesin,bytesout=1589364037
`
var responseStatus = `TITLE,OpenVPN 2.3.2 x86_64-pc-linux-gnu [SSL (OpenSSL)] [LZO] [EPOLL] [PKCS11] [eurephia] [MH] [IPv6] built on Dec  1 2014
TIME,Mon Dec 26 21:02:26 2016,1482782546
HEADER,CLIENT_LIST,Common Name,Real Address,Virtual Address,Bytes Received,Bytes Sent,Connected Since,Connected Since (time_t),Username
CLIENT_LIST,democlient1.example.com,11.112.113.114:58978,10.8.0.30,131180,414251,Mon Dec 26 20:57:17 2016,1482782237,UNDEF
CLIENT_LIST,rpi-v3,11.112.113.114:35331,10.8.0.14,5107,6909,Mon Dec 26 20:59:42 2016,1482782382,UNDEF
HEADER,ROUTING_TABLE,Virtual Address,Common Name,Real Address,Last Ref,Last Ref (time_t)
ROUTING_TABLE,10.8.0.30,democlient1.example.com,11.112.113.114:58978,Mon Dec 26 21:02:25 2016,1482782545
ROUTING_TABLE,10.8.0.14,rpi-v3,11.112.113.114:35331,Mon Dec 26 20:59:43 2016,1482782383
GLOBAL_STATS,Max bcast/mcast queue length,0
END
`
var responseKillSession = `SUCCESS: common name 'test1.example.com' found, 1 client(s) killed
`
var responseKillSessionError = `failure
`

func TestVersionEmptyStr(t *testing.T) {
	_, err := mi.ParseVersion("")
	assert.NotNil(t, err)
}

func TestVersionParser(t *testing.T) {
	v, err := mi.ParseVersion(responseVersion)
	assert.Nil(t, err)

	openVpnExpected := "OpenVPN 2.3.2 x86_64-pc-linux-gnu [SSL (OpenSSL)] [LZO] [EPOLL] [PKCS11] [eurephia] [MH] [IPv6] built on Dec  1 2014"
	assert.Equal(t, openVpnExpected, v.OpenVPN)
	assert.Equal(t, "1", v.Management)
}

func TestPidEmptyStr(t *testing.T) {
	_, err := mi.ParsePid("")
	assert.NotNil(t, err)
}

func TestPidErrorStr(t *testing.T) {
	_, err := mi.ParsePid(responseError)
	assert.NotNil(t, err)
}

func TestPidWrongNrOfLines(t *testing.T) {
	_, err := mi.ParsePid("a\n\n\na")
	assert.NotNil(t, err)
}

func TestPidParser(t *testing.T) {
	pid, err := mi.ParsePid(responsePid)
	assert.Nil(t, err)
	assert.Equal(t, int64(10869), pid)
}

func TestLoadStatsEmptyStr(t *testing.T) {
	_, err := mi.ParseStats("")
	assert.NotNil(t, err)
}

func TestLoadStatsWrongNrOfLines(t *testing.T) {
	_, err := mi.ParseStats("a\n\na")
	assert.NotNil(t, err)
}

func TestStatsParser(t *testing.T) {
	ls, err := mi.ParseStats(responseLoadStats)
	assert.Nil(t, err)
	assert.Equal(t, int64(12), ls.NClients)
	assert.Equal(t, int64(109984474), ls.BytesIn)
	assert.Equal(t, int64(1589364037), ls.BytesOut)
}

func TestStatsParserError(t *testing.T) {
	_, err := mi.ParseStats(responseLoadStatsBroken)
	assert.NotNil(t, err)
}

func TestStatsErrorStr(t *testing.T) {
	_, err := mi.ParseStats(responseError)
	assert.NotNil(t, err)
}

func TestStatusParser(t *testing.T) {
	s, err := mi.ParseStatus(responseStatus)
	assert.Nil(t, err)
	sExpected := &mi.Status{
		Title: "OpenVPN 2.3.2 x86_64-pc-linux-gnu [SSL (OpenSSL)] [LZO] [EPOLL] [PKCS11] [eurephia] [MH] [IPv6] built on Dec  1 2014",
		Time:  "Mon Dec 26 21:02:26 2016",
		TimeT: "1482782546",
		ClientList: []*mi.OVClient{
			{
				CommonName:      "democlient1.example.com",
				RealAddress:     "11.112.113.114:58978",
				VirtualAddress:  "10.8.0.30",
				BytesReceived:   131180,
				BytesSent:       414251,
				ConnectedSince:  "Mon Dec 26 20:57:17 2016",
				ConnectedSinceT: "1482782237",
				Username:        "UNDEF",
			},
			{
				CommonName:      "rpi-v3",
				RealAddress:     "11.112.113.114:35331",
				VirtualAddress:  "10.8.0.14",
				BytesReceived:   5107,
				BytesSent:       6909,
				ConnectedSince:  "Mon Dec 26 20:59:42 2016",
				ConnectedSinceT: "1482782382",
				Username:        "UNDEF",
			},
		},
		RoutingTable: []*mi.RoutingPath{
			{
				VirtualAddress: "10.8.0.30",
				CommonName:     "democlient1.example.com",
				RealAddress:    "11.112.113.114:58978",
				LastRef:        "Mon Dec 26 21:02:25 2016",
				LastRefT:       "1482782545",
			},
			{
				VirtualAddress: "10.8.0.14",
				CommonName:     "rpi-v3",
				RealAddress:    "11.112.113.114:35331",
				LastRef:        "Mon Dec 26 20:59:43 2016",
				LastRefT:       "1482782383",
			},
		},
	}
	assert.Equal(t, sExpected, s)
}

func TestKillSessionParser(t *testing.T) {
	m, err := mi.ParseKillSession(responseKillSession)
	assert.Nil(t, err)
	assert.Equal(t, "common name 'test1.example.com' found, 1 client(s) killed",
		m)
}

func TestKillSessionWrongNrOfLines(t *testing.T) {
	_, err := mi.ParseKillSession("a\n\na")
	assert.NotNil(t, err)
}

func TestKillSessionResponseError(t *testing.T) {
	_, err := mi.ParseKillSession(responseKillSessionError)
	assert.NotNil(t, err)
}
