
.PHONY: clean
clean:
	@rm -rf target

.PHONY: bin
bin: clean
	@go build -o target/timetrace .