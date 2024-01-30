# the data file must be cleaned without empty lines.
# also no " to avoid escape work
#
# sed '/^$/d' < in > out
# sed '/"/d' < in >out
compile:
	go build -o /dev/null cmd/embedding.go

fmt:
	gofmt -w cmd/embedding.go
