# strs

This package contains string utilities and algorithms. New `MatchSplit`(non-recursive) function is about 2x faster than recursive version.

```bash
Mike@/tmp:$ go test -bench=. t_test.go -cpuprofile=c.out
BenchmarkMatch-2         1000000              1144 ns/op
BenchmarkMatchNew-2      2000000             **607 ns/op**
PASS
ok      command-line-arguments  3.052s
Mike@/tmp:$
```
