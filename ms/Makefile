
all: ms ms.html

clean:
	rm -f  ms ms.html

ms: ms.go
	go build ms.go

ms.html: ms.raml
	raml2html ms.raml > ms.html
	#raml2md ms.raml | pandoc > ms.html

