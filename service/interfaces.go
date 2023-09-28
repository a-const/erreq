package service

type Get interface {
	Request(params []string, port string) any
}

type Post interface {
	Request(body any, port string) any
}
