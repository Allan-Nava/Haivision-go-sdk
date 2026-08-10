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
		t.Errorf("Source.NumPackets = %d", obj.Route.Source.NumPackets)
	}
	if obj.Route.Source.ElapsedRunningTime != "00:03:46" {
		t.Errorf("ElapsedRunningTime = %q (la doc lo dà come stringa HH:MM:SS)",
			obj.Route.Source.ElapsedRunningTime)
	}
}

// TestKnownBug_StatsFractionalBitrate documenta un bug APERTO, non un comportamento atteso.
//
// La doc dà `bitrate`, `sendRate` e `usedBandwidth` come `number` in **Mbit/s**, quindi
// frazionari, ma i modelli in haivision/stats li tipizzano `int`: un solo valore non intero fa
// fallire l'INTERA chiamata Get*Statistics. Item di backlog: `stats-float64` (v2.0.0, breaking
// perché cambia il tipo di campi esportati).
//
// Quando quell'item sarà chiuso questo test FALLIRÀ: va riscritto per asserire i valori.
func TestKnownBug_StatsFractionalBitrate(t *testing.T) {
	var m stats.SourceStatisticsModel
	err := json.Unmarshal([]byte(`{"name":"src","bitrate":4.5,"sendRate":1.25}`), &m)
	if err == nil {
		t.Fatalf("l'unmarshal ora RIESCE: il bug `stats-float64` è stato corretto — " +
			"riscrivi questo test per asserire bitrate=4.5 e sendRate=1.25")
	}
	t.Logf("bug aperto confermato (`stats-float64`): %v", err)
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
