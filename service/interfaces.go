package service

type Get interface {
	Request(params []string) any
}

type Post interface {
	Request(body string)
}
