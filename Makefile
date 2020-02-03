.PHONY: build-all
build-all:
	@packr
	@rsrc -manifest eqgamepatch.manifest -o rsrc.syso
	@GOOS=windows go build -o bin/eqgamepatch.exe main.go
.PHONY: sandbox
sandbox:
	@packr
	@go run client/main.go
LENGTH := 16
OFFSET := 0
.PHONY: peek
peek:
	@echo "peeking at offset $(OFFSET) length $(LENGTH)"
	@echo "abysmal"
	@hexdump -e '8/1 "0x%02X, ""\t"" "' -e '8/1 "%c""\n"' -s $(OFFSET) -n $(LENGTH) s3d/abysmal_obj.s3d
	@echo "test"
	@hexdump -e '8/1 "0x%02X, ""\t"" "' -e '8/1 "%c""\n"' -s $(OFFSET) -n $(LENGTH) s3d/test.s3d
.PHONY: test
test:
	@rm -rf s3d/out/*
	@go test -v ./...
.PHNOY: hex
hex:
	echo $((16#))