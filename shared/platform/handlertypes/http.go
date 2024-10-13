package handlertypes

type Response struct {
	Body       interface{}
	HttpStatus int
}

type Meta struct {
	RequestId string
}

type Request struct {
	Body    string
	Headers Headers
	Query   interface{}
	Params  interface{}
	Meta    *Meta
}

type Headers struct {
	Authorization  string `header:"Authorization"`
	Connection     string `header:"Connection"`
	Accept         string `header:"Accept"`
	ContentType    string `header:"Content-Type"`
	UserAgent      string `header:"User-Agent"`
	AcceptEncoding string `header:"Accept-Encoding"`
	ContentLength  string `header:"Content-Length"`
}
