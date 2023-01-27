package haivision


type RequestInitSession struct {
	Username string `json:username required:"true" validate:"nonnil,min=1"`
	Password string `json:password required:"true" validate:"nonnil,min=1"`
}