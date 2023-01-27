package haivision


type RequestInitSession struct {
	Username string `json:username required:"true" validate:"nonnil,min=1"`
	Password string `json:password required:"true" validate:"nonnil,min=1"`
}
//
type RequestSourceModelUDPandRTP struct {
    Name string `json:name required:"true" validate:"nonnil,min=1"`
    ID string `json:id required:"true" validate:"nonnil,min=1"`
    Address string `json:address required:"true" validate:"nonnil,min=1"`
    Protocol	 string `json:protocol required:"true" validate:"nonnil,min=1"`
    Port int `json:port required:"true" validate:"nonnil,min=1"` 
    NetworkInterface	 string `json:networkInterface required:"true" validate:"nonnil,min=1"` 
    RetainHeader string `json:retainHeader required:"true" validate:"nonnil,min=1"` 
    SourceAddress string `json:sourceAddress required:"true" validate:"nonnil,min=1"`
    Fec string `json:fec required:"true" validate:"nonnil,min=1"`
}


type RequestSourceModelSRT struct {
    Name string `json:name required:"true" validate:"nonnil,min=1"`
    ID string `json:id required:"true" validate:"nonnil,min=1"`
    Address string `json:address required:"true" validate:"nonnil,min=1"`
    Protocol	 string `json:protocol required:"true" validate:"nonnil,min=1"`
    Port int `json:port required:"true" validate:"nonnil,min=1"` 
    NetworkInterface	 string `json:networkInterface required:"true" validate:"nonnil,min=1"` 
}