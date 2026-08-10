package test

import (
	"encoding/json"
	"testing"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/stats"
)

// Statistiche di route con valori interi: questo caso funziona già.
func TestRouteStatisticsResponse(t *testing.T) {
	documented := `{
	  "collectedAt": 1675178018888,
	  "route": {
	    "name": "route-srt",
	    "elapsedRunningTime": "00:03:46",
	    "id": "route-1",
	    "state": "connected",
	    "source": {
	      "name": "src",
	      "id": "src-1",
	      "mode": "unicast",
	      "elapsedRunningTime": "00:03:46",
	      "signalLosses": 0,
	      "sendRate": 5,
	      "numPackets": 12345,
	      "usedBandwidth": 5,
	      "bitrate": 5,
	      "state": "connected",
	      "fecLostPackets": 0,
	      "fecRecoveredPackets": 0,
	      "fecUnrecoveredPackets": 0,
	      "fecReorderedPackets": 0
	    },
	    "destinations": []
	  }
	}`
	var obj stats.ResponseRouteStatistics
	if err := json.Unmarshal([]byte(documented), &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if obj.Route.ID != "route-1" {
		t.Errorf("Route.ID = %q", obj.Route.ID)
	}
	if obj.Route.Source.NumPackets != 12345 {
		t.Errorf("Source.NumPackets = %v", obj.Route.Source.NumPackets)
	}
	if obj.Route.Source.ElapsedRunningTime != "00:03:46" {
		t.Errorf("ElapsedRunningTime = %q (la doc lo dà come stringa HH:MM:SS)",
			obj.Route.Source.ElapsedRunningTime)
	}
}

// I valori frazionari sono il caso NORMALE: la doc dà bitrate/sendRate/usedBandwidth come
// `number` in Mbit/s. Fino alla v1.x erano tipizzati `int` e un solo 4.5 faceva fallire l'INTERA
// chiamata Get*Statistics con "cannot unmarshal number 4.5 into Go value of type int".
func TestStatsAcceptsFractionalValues(t *testing.T) {
	var m stats.SourceStatisticsModel
	body := `{"name":"src","bitrate":4.5,"sendRate":1.25,"usedBandwidth":0.338,"numPackets":12345}`
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		t.Fatalf("unmarshal con valori frazionari: %v", err)
	}
	if m.Bitrate != 4.5 {
		t.Errorf("Bitrate = %v, atteso 4.5", m.Bitrate)
	}
	if m.SendRate != 1.25 {
		t.Errorf("SendRate = %v, atteso 1.25", m.SendRate)
	}
	if m.UsedBandwidth != 0.338 {
		t.Errorf("UsedBandwidth = %v, atteso 0.338", m.UsedBandwidth)
	}
	if m.NumPackets != 12345 {
		t.Errorf("NumPackets = %v, atteso 12345", m.NumPackets)
	}
}

// Anche i contatori sono float64: JSON non distingue 42 da 42.0, e un gateway che serializzasse
// un contatore come 42.0 romperebbe di nuovo tutto se fossero `int`.
func TestStatsCountersAcceptTrailingZero(t *testing.T) {
	var m stats.SourceStatisticsModel
	if err := json.Unmarshal([]byte(`{"numPackets":42.0,"signalLosses":0.0}`), &m); err != nil {
		t.Fatalf("unmarshal con contatori in forma 42.0: %v", err)
	}
	if m.NumPackets != 42 {
		t.Errorf("NumPackets = %v", m.NumPackets)
	}
}

func TestSourceStatisticsResponse(t *testing.T) {
	documented := `{
	  "collectedAt": 1675178018888,
	  "source": {"name":"src","id":"src-1","mode":"unicast","state":"connected","numPackets":42}
	}`
	var obj stats.ResponseSourceStatistics
	if err := json.Unmarshal([]byte(documented), &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if obj.Source.ID != "src-1" || obj.Source.NumPackets != 42 {
		t.Errorf("Source = %+v", obj.Source)
	}
	if obj.CollectedAt != 1675178018888 {
		t.Errorf("CollectedAt = %d", obj.CollectedAt)
	}
}
