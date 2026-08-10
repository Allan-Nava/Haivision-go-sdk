package srt

// Modelli di richiesta SRT.
//
// Riferimento: REST API Integrator's Reference (HMG 3.7.x), «Source object model» e
// «Destination object model». I tipi seguono gli esempi LETTERALI della doc: `ttl`, `tos` e
// `mtu` sono **numeri** (`"ttl": 64, "mtu": 1496, "tos": 136`), non stringhe come erano
// tipizzati fino alla v1.x.
//
// I campi opzionali sono puntatori con `omitempty`: così un campo non valorizzato viene
// **omesso** dal body invece di essere inviato come stringa vuota o zero. Fino alla v1.x
// erano `string` obbligatorie, e la validazione rifiutava qualsiasi route che non le
// specificasse tutte.

type RequestSourceModelSRT struct {
	Name             string `json:"name" validate:"required"`
	ID               string `json:"id,omitempty"`
	Address          string `json:"address" validate:"required"`
	Protocol         string `json:"protocol" validate:"required"`
	Port             int    `json:"port" validate:"required,min=1,max=65535"`
	NetworkInterface string `json:"networkInterface"`
	// opzionali SRT
	Mode          *string `json:"mode,omitempty"`
	SrtLatency    *int    `json:"srtLatency,omitempty"`
	SrtRcvBuf     *int    `json:"srtRcvBuf,omitempty"`
	SrtStreamId   *string `json:"srtStreamId,omitempty"`
	SrtPassPhrase *string `json:"srtPassPhrase,omitempty"`
	Encryption    *string `json:"encryption,omitempty"`
	UseFec        *bool   `json:"useFec,omitempty"`
}

type RequestDestinationModelSrt struct {
	Name             string `json:"name" validate:"required"`
	ID               string `json:"id,omitempty"`
	Address          string `json:"address" validate:"required"`
	Protocol         string `json:"protocol" validate:"required"`
	Port             int    `json:"port" validate:"required,min=1,max=65535"`
	NetworkInterface string `json:"networkInterface"`
	// Action avvia o ferma QUESTA destinazione dentro una update di route
	// (route.DestinationActionStart / Stop). Omessa, la destinazione resta com'è.
	Action *string `json:"action,omitempty"`
	// opzionali: numerici come da doc
	Ttl           *int    `json:"ttl,omitempty"`
	Tos           *int    `json:"tos,omitempty"`
	Mtu           *int    `json:"mtu,omitempty"`
	RetainHeader  *bool   `json:"retainHeader,omitempty"`
	SrtLatency    *int    `json:"srtLatency,omitempty"`
	SrtPassPhrase *string `json:"srtPassPhrase,omitempty"`
	Encryption    *string `json:"encryption,omitempty"`
	UseFec        *bool   `json:"useFec,omitempty"`
	SrtFecCols    *int    `json:"srtFecCols,omitempty"`
	SrtFecRows    *int    `json:"srtFecRows,omitempty"`
}
