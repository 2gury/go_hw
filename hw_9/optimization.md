go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 
go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 -benchtime=10x 
go tool pprof -http=:8083 optimization.test mem.out  
go tool pprof -http=:8083 optimization.test cpu.out  