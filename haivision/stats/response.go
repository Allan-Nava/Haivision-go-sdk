package stats

/*
name	string	Name of the source.
id	string	Unique identifier for the source.
mode	string	Unicast or multicast.
elapsedRunningTime	string	Either an empty string (for idle routes), or a string in HH:MM:SS format (e.g. 00:03:46).
signalLosses	number	Number of signal losses.
sendRate	number	Packet send rate in Mbits/s.
numPackets	number	Number of packets sent.
usedBandwidth	number	Bandwidth used in Mbits/s.
bitrate	number
Stream bitrate in Mbits/s.
state	string	Source state: disconnected, connecting, connection established, or connected.
fecLostPackets	number	Number of lost FEC packets.
fecRecoveredPackets	number	Number of recovered FEC packets.
fecUnrecoveredPackets	number	Number of unrecovered FEC packets.
fecReorderedPackets
*/
type SourceStatisticsModel struct {
	Name                  string `json:"name"`
	ID                    string `json:"id"`
	Mode                  string `json:"mode"`
	ElapsedRunningTime    string `json:"elapsedRunningTime"`
	SignalLosses          int    `json:"signalLosses"`
	SendRate              int    `json:"sendRate"`
	NumPackets            int    `json:"numPackets"`
	UsedBandwidth         int    `json:"usedBandwidth"`
	Bitrate               int    `json:"bitrate"`
	State                 string `json:"state"`
	FecLostPackets        int    `json:"fecLostPackets"`
	FecRecoveredPackets   int    `json:"fecRecoveredPackets"`
	FecUnrecoveredPackets int    `json:"fecUnrecoveredPackets"`
	FecReorderedPackets   int    `json:"fecReorderedPackets"`
}

/*
UDP, RTP
Name	Type	Description
name	string	Name of the destination.
id	string	Unique identifier for the destination.
mode	string	Destination mode: unicast or multicast.
state	string	Destination state: disconnected, connecting, connection established, or connected.
elapsedRunningTime	string	Either an empty string (for idle routes), or a string in HH:MM:SS format (e.g. 00:03:46).
bitrate	number	Stream bitrate in Mbits/s.
signalLosses	number	Number of signal losses.
usedBandwidth	number	Bandwidth used in Mbits/s.
sendRate	number	Packet send rate in Mbits/s.
numPackets	number	Number of packets.
*/
type DestinationStatisticsUdpRtpHlsModel struct {
	Name               string `json:"name"`
	ID                 string `json:"id"`
	Mode               string `json:"mode"`
	State              string `json:"state"`
	ElapsedRunningTime string `json:"elapsedRunningTime"`
	Bitrate            int    `json:"bitrate"`
	SignalLosses       int    `json:"signalLosses"`
	UsedBandwidth      int    `json:"usedBandwidth"`
	SendRate           int    `json:"sendRate"`
	NumPackets         int    `json:"numPackets"`
}

/*
name	string	Name of the destination.
id	string	Unique identifier for the destination.
mode	string	Destination SRT mode: caller, listener, or rendezvous.
protocol	string	Protocol: srt.
state	string	Destination state: disconnected, connecting, connection established, or connected.
elapsedRunningTime	string	Either an empty string (for idle routes), or a string in HH:MM:SS format (e.g. 00:03:46).
bitrate	number	Stream bitrate in Mbits/s.
signalLosses	number	Number of signal losses.
usedBandwidth	number	Bandwidth used in Mbits/s.
sendRate	number	Packet send rate in Mbits/s.
numPackets	number	Number of packets.
srtNumLostPackets	number	SRT number of lost (but recovered) packages.
srtPacketLossRate	number	SRT packet loss rate in percent.
srtNumSkippedPackets	number	(Receiver only) Missing packets skipped, because they were not recovered in time.
srtDroppedPackets	number	(Sender only) Number of dropped packets.
srtRoundTripTime	number	SRT round trip time in ms.
srtBufferLevel	number	SRT buffer time in ms.
srtNegotiatedLatency	number	SRT maximum latency in ms.
srtLatency	number	(Receiver only) SRT latency.
srtDecryptionState	string	(Receiver only) SRT receiver decryption state: <empty>, active, initializing, inactive (no passphrase), or inactive (invalid passphrase).
srtPeerDecryptionState	string	(Sender only) SRT peer decryption state: <empty>, active, initializing, inactive (no passphrase), or inactive (invalid passphrase).
srtEncryption	string	(Receiver only) Indicates the cipher used in the received stream: <empty>, AES128, or AES256.
srtMaxBandwidth	number	(Sender only) SRT maximum bandwidth used in Mbits/s.
srtRetransmitRate	number	SRT retransmit rate in bits/s.
srtEstimatedBandwidth	number	(Sender only) SRT estimated path max bandwidth in bits/s.
clientStat	object list	(SRT listener only) Array of route client statistics objects. See below for the client statistics model.
connections	object list	(SRT caller and rendezvous only). Array of destination connections objects. See SRT Statistics Connections Object Model for the definition.
*/
type DestinationStatisticsSrtModel struct {
	Name                   string                                     `json:"name"`
	ID                     string                                     `json:"id"`
	Mode                   string                                     `json:"mode"`
	Protocol               string                                     `json:"protocol"`
	State                  string                                     `json:"state"`
	ElapsedRunningTime     string                                     `json:"elapsedRunningTime"`
	Bitrate                int                                        `json:"bitrate"`
	SignalLosses           int                                        `json:"signalLosses"`
	UsedBandwidth          int                                        `json:"usedBandwidth"`
	SendRate               int                                        `json:"sendRate"`
	NumPackets             int                                        `json:"numPackets"`
	SrtNumLostPackets      int                                        `json:"srtNumLostPackets"`
	SrtPacketLossRate      int                                        `json:"srtPacketLossRate"`
	SrtNumSkippedPackets   int                                        `json:"srtNumSkippedPackets"`
	SrtDroppedPackets      int                                        `json:"srtDroppedPackets"`
	SrtRoundTripTime       int                                        `json:"srtRoundTripTime"`
	SrtBufferLevel         int                                        `json:"srtBufferLevel"`
	SrtNegotiatedLatency   int                                        `json:"srtNegotiatedLatency"`
	SrtLatency             int                                        `json:"srtLatency"`
	SrtDecryptionState     string                                     `json:"srtDecryptionState"`
	SrtPeerDecryptionState string                                     `json:"srtPeerDecryptionState"`
	SrtEncryption          string                                     `json:"srtEncryption"`
	SrtMaxBandwidth        int                                        `json:"srtMaxBandwidth"`
	SrtRetransmitRate      int                                        `json:"srtRetransmitRate"`
	SrtEstimatedBandwidth  int                                        `json:"srtEstimatedBandwidth"`
	ClientStat             []DestinationStatisticsSrtClientStatModel  `json:"clientStat"`
	Connections            []DestinationStatisticsSrtConnectionsModel `json:"connections"`
}

/*
label	string	Listener output label.
address	string
Client address.

port	number
Client connection port.

bitrate	number
Stream bitrate.

signalLosses	number
Number of signal losses.

srtVersion	string
SRT protocol version of the listener.

SRTPeerVersion	string	SRT protocol version of the client.
usedBandwidth	number
Bandwidth used in Mbits/s.

connections	object list	Array of destination connections objects. See SRT Statistics Connections Object Model for the definition.

*/

type DestinationStatisticsSrtClientStatModel struct {
	Label          string                                     `json:"label"`
	Address        string                                     `json:"address"`
	Port           int                                        `json:"port"`
	Bitrate        int                                        `json:"bitrate"`
	SignalLosses   int                                        `json:"signalLosses"`
	SrtVersion     string                                     `json:"srtVersion"`
	SrtPeerVersion string                                     `json:"srtPeerVersion"`
	UsedBandwidth  int                                        `json:"usedBandwidth"`
	Connections    []DestinationStatisticsSrtConnectionsModel `json:"connections"`
}

/*
SRTPeerVersion	string

(tick)	Peer SRT protocol version.
address	string	(tick)	(tick)	(tick)	Network address.
bitrate	number

(tick)	Stream bitrate in Mbits/s.
label	string

(tick)	Label.
localAddress	string	(tick)	(tick)	(tick)	Network local address.
port	number	(tick)	(tick)	(tick)	Network port.
signalLosses	number

(tick)	Number of signal losses.
localPort	number	(tick)	(tick)	(tick)	Network local port.
networkInterface	string	(tick)	(tick)	(tick)	Network interface.
numPackets	number
(tick)	(tick)	Number of packets.
srtBufferLevel	number	(tick)	(tick)	(tick)	SRT buffer time in ms.
srtCurrentBandwidth	number	(tick)	(tick)	(tick)	SRT current bandwidth.
srtDecryptionState	string	(tick)

SRT decryption state: <empty>, active, initializing, inactive (no passphrase), or inactive (invalid passphrase).
srtDroppedPackets	number
(tick)	(tick)	Number of dropped packets.
srtDroppedPacketsDiff	number
(tick)	(tick)	Number of dropped packets diff.
srtEncryption	string	(tick)	(tick)	(tick)	Indicates the cipher used in the stream: None, none, AES128, or AES256.
srtEstimatedBandwidth	number	(tick)	(tick)	(tick)	SRT estimated bandwidth.
srtFec	string	(tick)	(tick)	(tick)	SRT FEC.
srtFecArq	string	(tick)	(tick)	(tick)	SRT FEC ARQ: always, onreq, or never.
srtFecCols	number	(tick)	(tick)	(tick)	SRT FEC columns.
srtFecLayout	string	(tick)	(tick)	(tick)	SRT FEC layout: even or staircase.
srtFecPacketLoss	number	(tick)

SRT FEC packet loss.
srtFecRecoveredPackets	number	(tick)

SRT FEC recovered packets.
srtFecRows	number	(tick)	(tick)	(tick)	SRT FEC rows.
srtFecTotalPacketLoss	number	(tick)

SRT FEC total packet loss.
srtFecTotalRecoveredPackets	number	(tick)

SRT FEC total recoverd packets.
srtGroupMemberStatus	string	(tick)	(tick)	(tick)	SRT group member status.
srtGroupMemberWeight	number	(tick)	(tick)	(tick)	SRT group member weight.
srtGroupMode	string	(tick)	(tick)
SRT group mode: <empty>, broadcast, backup, balance, or any.
srtMaxBandwidth	number	(tick)	(tick)	(tick)	SRT maximum bandwidth used in Mbits/s.
srtNegotiatedLatency	number	(tick)	(tick)	(tick)	SRT negotiated latency in ms.
srtNumLostPackets	number	(tick)	(tick)	(tick)	SRT number of lost (but recovered) packets.
srtNumPackets	number	(tick)	(tick)	(tick)	SRT number of pacakets.
srtPacketLossRate	number	(tick)	(tick)	(tick)	SRT packet loss rate in percent.
srtPeerDecryptionState	string
(tick)	(tick)	SRT peer decryption state: <empty>, active, initializing, inactive (no passphrase), or inactive (invalid passphrase).
srtRetransmitRate	number	(tick)	(tick)	(tick)	SRT retransmit rate in bits/s.
srtRoundTripTime	number	(tick)	(tick)	(tick)	SRT round trip time in ms.
srtSkippedPackets	number	(tick)

SRT number of skipped packets, because they were not recovered in time.
srtSkippedPacketsDiff	number	(tick)

SRT number of skipped packets diff.
srtVersion	string

(tick)	SRT protocol version.
state	string	(tick)	(tick)	(tick)	Source connection state.
usedBandwidth	number

(tick)	Bandwidth used in Mbits/s.
*/
type DestinationStatisticsSrtConnectionsModel struct {
	SrtPeerVersion              string `json:"srtPeerVersion"`
	Address                     string `json:"address"`
	Bitrate                     int    `json:"bitrate"`
	Label                       string `json:"label"`
	LocalAddress                string `json:"localAddress"`
	Port                        int    `json:"port"`
	SignalLosses                int    `json:"signalLosses"`
	LocalPort                   int    `json:"localPort"`
	NetworkInterface            string `json:"networkInterface"`
	NumPackets                  int    `json:"numPackets"`
	SrtBufferLevel              int    `json:"srtBufferLevel"`
	SrtCurrentBandwidth         int    `json:"srtCurrentBandwidth"`
	SrtDecryptionState          string `json:"srtDecryptionState"`
	SrtDroppedPackets           int    `json:"srtDroppedPackets"`
	SrtDroppedPacketsDiff       int    `json:"srtDroppedPacketsDiff"`
	SrtEncryption               string `json:"srtEncryption"`
	SrtEstimatedBandwidth       int    `json:"srtEstimatedBandwidth"`
	SrtFec                      string `json:"srtFec"`
	SrtFecArq                   string `json:"srtFecArq"`
	SrtFecCols                  int    `json:"srtFecCols"`
	SrtFecLayout                string `json:"srtFecLayout"`
	SrtFecPacketLoss            int    `json:"srtFecPacketLoss"`
	SrtFecRecoveredPackets      int    `json:"srtFecRecoveredPackets"`
	SrtFecRows                  int    `json:"srtFecRows"`
	SrtFecTotalPacketLoss       int    `json:"srtFecTotalPacketLoss"`
	SrtFecTotalRecoveredPackets int    `json:"srtFecTotalRecoveredPackets"`
	SrtGroupMemberStatus        string `json:"srtGroupMemberStatus"`
	SrtGroupMemberWeight        int    `json:"srtGroupMemberWeight"`
	SrtGroupMode                string `json:"srtGroupMode"`
	SrtMaxBandwidth             int    `json:"srtMaxBandwidth"`
	SrtNegotiatedLatency        int    `json:"srtNegotiatedLatency"`
	SrtNumLostPackets           int    `json:"srtNumLostPackets"`
	SrtNumPackets               int    `json:"srtNumPackets"`
	SrtPacketLossRate           int    `json:"srtPacketLossRate"`
	SrtPeerDecryptionState      string `json:"srtPeerDecryptionState"`
	SrtRetransmitRate           int    `json:"srtRetransmitRate"`
	SrtRoundTripTime            int    `json:"srtRoundTripTime"`
	SrtSkippedPackets           int    `json:"srtSkippedPackets"`
	SrtSkippedPacketsDiff       int    `json:"srtSkippedPacketsDiff"`
	SrtVersion                  string `json:"srtVersion"`
	State                       string `json:"state"`
	UsedBandwidth               int    `json:"usedBandwidth"`
}

/*
	{
	  "collectedAt": [Date/time in Unix time],
	  "route": {
	    "name": "[Route Name]",
	    "elapsedRunningTime": "00:00:14",
	    "id": "[Route ID]",
	    "state": "running",
	    "source": {
	      <Source Statistics object>
	    },
	    "destinations": [
	      {
	        <Destination Statistics object>
	      }
	    ]
	  }
	}
*/
type ResponseRouteStatistics struct {
	CollectedAt int64 `json:"collectedAt"`
	Route       struct {
		Name               string                `json:"name"`
		ElapsedRunningTime string                `json:"elapsedRunningTime"`
		ID                 string                `json:"id"`
		State              string                `json:"state"`
		Source             SourceStatisticsModel `json:"source"`
		Destinations       []interface{}         `json:"destinations"`
	} `json:"route"`
}

//

/*
	{
	  "collectedAt": [Date/time in Unix time],
	  "source": {
	    <Source Statistics Object>
	  }
	}
*/
type ResponseSourceStatistics struct {
	CollectedAt int64                `json:"collectedAt"`
	Source      SourceStatisticsModel `json:"source"`
}

/*
	{
	  "collectedAt": [Date/time in Unix time],
	  "destination": {
	    <Destination Statistics Object>
	  }
	}
*/
type ResponseDestinationStatistics struct {
	CollectedAt string      `json:"collectedAt"`
	Destination interface{} `json:"destination"`
}

/*
	{
	  "collectedAt": [Date/time in Unix time],
	  "clientStat": [
	    <Client Statistics Object>
	  ]
	}
*/
type ResponseSrtClientStatistics struct {
	CollectedAt int64        `json:"collectedAt"`
	ClientStat  []interface{} `json:"clientStat"`
}
