go test -run=^# -v -bench=BenchmarkCodec_ -benchtime=$2x -count 10 -benchmem -cpuprofile=cpu.out -memprofile=mem.out -trace=trace.out | tee $1_$2_bench.txt
#go tool pprof -http :8080 cpu.out
#go tool pprof -http :8081 mem.out
#go tool trace trace.out

#go tool pprof $FILENAME.test cpu.out
# (pprof) list <func name>

# go get -u golang.org/x/perf/cmd/benchstat
#benchstat cbor_$1_bench.txt
rm cpu.out mem.out trace.out *.test
