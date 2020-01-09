.PHONY: build-all
build-all:
	@packr
	@fyne-cross --targets=linux/amd64,windows/amd64,darwin/amd64 github.com/xackery/eqgamepatch