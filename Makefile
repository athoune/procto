linux:
	docker run -it --rm -v `pwd`:/go golang go build -v -o hsperfcarbon
