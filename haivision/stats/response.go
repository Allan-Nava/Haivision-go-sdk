package stats

// Modelli delle statistiche.
//
// TUTTI i campi numerici sono `float64`, non `int`. La doc Haivision li dà come `number` e
// molti sono misure in Mbit/s (`bitrate`, `sendRate`, `usedBandwidth`), quindi frazionarie: con
// `int` un solo valore come 4.5 faceva fallire l'INTERA chiamata Get*Statistics con
// "cannot unmarshal number 4.5 into Go value of type int".
//
// Anche i contatori (numPackets, signalLosses, fec*Packets) sono float64: JSON non distingue
// 42 da 42.0, e un gateway che serializzasse un contatore come 42.0 romperebbe di nuovo tutto.
// Le statistiche finiscono comunque in sistemi di metriche che usano float64.
// Restano `int` solo `port` e `localPort`, che sono identificatori e non misure.

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
	Name                  string  `json:"name"`
	ID                    string  `json:"id"`
	Mode                  string  `json:"mode"`
	ElapsedRunningTime    string  `json:"elapsedRunningTime"`
	SignalLosses          float64 `json:"signalLosses"`
	SendRate              float64 `json:"sendRate"`
	NumPackets            float64 `json:"numPackets"`
	UsedBandwidth         float64 `json:"usedBandwidth"`
	Bitrate               float64 `json:"bitrate"`
	State                 string  `json:"state"`
	FecLostPackets        float64 `json:"fecLostPackets"`
	FecRecoveredPackets   float64 `json:"fecRecoveredPackets"`
	FecUnrecoveredPackets float64 `json:"fecUnrecoveredPackets"`
	FecReorderedPackets   float64 `json:"fecReorderedPackets"`
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
	Name               string  `json:"name"`
	ID                 string  `json:"id"`
	Mode               string  `json:"mode"`
	State              string  `json:"state"`
	ElapsedRunningTime string  `json:"elapsedRunningTime"`
	Bitrate            float64 `json:"bitrate"`
	SignalLosses       float64 `json:"signalLosses"`
	UsedBandwidth      float64 `json:"usedBandwidth"`
	SendRate           float64 `json:"sendRate"`
	NumPackets         float64 `json:"numPackets"`
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
	Bitrate                float64                                    `json:"bitrate"`
	SignalLosses           float64                                    `json:"signalLosses"`
	UsedBandwidth          float64                                    `json:"usedBandwidth"`
	SendRate               float64                                    `json:"sendRate"`
	NumPackets             float64                                    `json:"numPackets"`
	SrtNumLostPackets      float64                                    `json:"srtNumLostPackets"`
	SrtPacketLossRate      float64                                    `json:"srtPacketLossRate"`
	SrtNumSkippedPackets   float64                                    `json:"srtNumSkippedPackets"`
	SrtDroppedPackets      float64                                    `json:"srtDroppedPackets"`
	SrtRoundTripTime       float64                                    `json:"srtRoundTripTime"`
	SrtBufferLevel         float64                                    `json:"srtBufferLevel"`
	SrtNegotiatedLatency   float64                                    `json:"srtNegotiatedLatency"`
	SrtLatency             float64                                    `json:"srtLatency"`
	SrtDecryptionState     string                                     `json:"srtDecryptionState"`
	SrtPeerDecryptionState string                                     `json:"srtPeerDecryptionState"`
	SrtEncryption          string                                     `json:"srtEncryption"`
	SrtMaxBandwidth        float64                                    `json:"srtMaxBandwidth"`
	SrtRetransmitRate      float64                                    `json:"srtRetransmitRate"`
	SrtEstimatedBandwidth  float64                                    `json:"srtEstimatedBandwidth"`
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
	Bitrate        float64                                    `json:"bitrate"`
	SignalLosses   float64                                    `json:"signalLosses"`
	SrtVersion     string                                     `json:"srtVersion"`
	SrtPeerVersion string                                     `json:"srtPeerVersion"`
	UsedBandwidth  float64                                    `json:"usedBandwidth"`
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
	SrtPeerVersion              string  `json:"srtPeerVersion"`
	Address                     string  `json:"address"`
	Bitrate                     float64 `json:"bitrate"`
	Label                       string  `json:"label"`
	LocalAddress                string  `json:"localAddress"`
	Port                        int     `json:"port"`
	SignalLosses                float64 `json:"signalLosses"`
	LocalPort                   int     `json:"localPort"`
	NetworkInterface            string  `json:"networkInterface"`
	NumPackets                  float64 `json:"numPackets"`
	SrtBufferLevel              float64 `json:"srtBufferLevel"`
	SrtCurrentBandwidth         float64 `json:"srtCurrentBandwidth"`
	SrtDecryptionState          string  `json:"srtDecryptionState"`
	SrtDroppedPackets           float64 `json:"srtDroppedPackets"`
	SrtDroppedPacketsDiff       float64 `json:"srtDroppedPacketsDiff"`
	SrtEncryption               string  `json:"srtEncryption"`
	SrtEstimatedBandwidth       float64 `json:"srtEstimatedBandwidth"`
	SrtFec                      string  `json:"srtFec"`
	SrtFecArq                   string  `json:"srtFecArq"`
	SrtFecCols                  float64 `json:"srtFecCols"`
	SrtFecLayout                string  `json:"srtFecLayout"`
	SrtFecPacketLoss            float64 `json:"srtFecPacketLoss"`
	SrtFecRecoveredPackets      float64 `json:"srtFecRecoveredPackets"`
	SrtFecRows                  float64 `json:"srtFecRows"`
	SrtFecTotalPacketLoss       float64 `json:"srtFecTotalPacketLoss"`
	SrtFecTotalRecoveredPackets float64 `json:"srtFecTotalRecoveredPackets"`
	SrtGroupMemberStatus        string  `json:"srtGroupMemberStatus"`
	SrtGroupMemberWeight        float64 `json:"srtGroupMemberWeight"`
	SrtGroupMode                string  `json:"srtGroupMode"`
	SrtMaxBandwidth             float64 `json:"srtMaxBandwidth"`
	SrtNegotiatedLatency        float64 `json:"srtNegotiatedLatency"`
	SrtNumLostPackets           float64 `json:"srtNumLostPackets"`
	SrtNumPackets               float64 `json:"srtNumPackets"`
	SrtPacketLossRate           float64 `json:"srtPacketLossRate"`
	SrtPeerDecryptionState      string  `json:"srtPeerDecryptionState"`
	SrtRetransmitRate           float64 `json:"srtRetransmitRate"`
	SrtRoundTripTime            float64 `json:"srtRoundTripTime"`
	SrtSkippedPackets           float64 `json:"srtSkippedPackets"`
	SrtSkippedPacketsDiff       float64 `json:"srtSkippedPacketsDiff"`
	SrtVersion                  string  `json:"srtVersion"`
	State                       string  `json:"state"`
	UsedBandwidth               float64 `json:"usedBandwidth"`
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
	CollectedAt int64                 `json:"collectedAt"`
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
	CollectedAt int64       `json:"collectedAt"`
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
	CollectedAt int64         `json:"collectedAt"`
	ClientStat  []interface{} `json:"clientStat"`
}
