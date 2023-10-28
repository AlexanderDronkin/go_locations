package main

type jsonWriter interface {
	toJson() []byte
}

type dbConnector interface {
	Init()
	Connect(dsn string, database string) interface{}
	Close()
}

type apiSender interface {
	Init()
	Auth()
	Send()
}
