.PHONY: test
test:
	go test -v -coverprofile cover.out ./... && go tool cover -html=cover.out

.PHONY: bench
bench:
	go test -v -bench . -benchmem -count=4 -cpuprofile=cpu.out -memprofile=mem.out && go tool pprof -text cpu.out && go tool pprof -text mem.out