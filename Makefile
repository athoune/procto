linux:
	#docker run -it --rm -v `pwd`:/go golang go get
	docker run -it --rm -v `pwd`:/go golang go build -v -o procto
