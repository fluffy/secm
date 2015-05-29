
all: ks ks.html

clean:
	rm -f  ks ks.html

ks: ks.go
	go build ks.go

ks.html: ks.raml
	raml2html ks.raml > ks.html
	#raml2md ks.raml | pandoc > ks.html

