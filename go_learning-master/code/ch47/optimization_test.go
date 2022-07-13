package profiling

import "testing"

// go test -bench=BenchmarkProcessRequestOld optimization_test.go optmization.go structs.go -v -cpuprofile=cpu.prof -memprofile=mem.prof
// go tool pprof mem.prof
// ...
// go tool pprof cpu.prof
// (pprof) top -cum
// Showing nodes accounting for 140ms, 7.14% of 1960ms total
// Showing top 10 nodes out of 209
//       flat  flat%   sum%        cum   cum%
//          0     0%     0%     1600ms 81.63%  command-line-arguments.BenchmarkProcessRequestOld
//          0     0%     0%     1600ms 81.63%  command-line-arguments.processRequestOld
//          0     0%     0%     1600ms 81.63%  testing.(*B).launch
//          0     0%     0%     1600ms 81.63%  testing.(*B).runN
//          0     0%     0%     1040ms 53.06%  encoding/json.Unmarshal
//          0     0%     0%      920ms 46.94%  encoding/json.(*decodeState).object
//          0     0%     0%      920ms 46.94%  encoding/json.(*decodeState).unmarshal
//       20ms  1.02%  1.02%      920ms 46.94%  encoding/json.(*decodeState).value
//       50ms  2.55%  3.57%      880ms 44.90%  encoding/json.(*decodeState).array
//       70ms  3.57%  7.14%      520ms 26.53%  encoding/json.(*decodeState).literalStore

// (pprof) list processRequestOld
// Total: 1.96s
// ROUTINE ======================== command-line-arguments.processRequestOld in /d/work/src/go/GolangComplete/go_learning-master/code/ch47/optmization.go
//          0      1.60s (flat, cum) 81.63% of Total
//          .          .     44:
//          .          .     45:func processRequestOld(reqs []string) []string {
//          .          .     46:   reps := []string{}
//          .          .     47:   for _, req := range reqs {
//          .          .     48:           reqObj := &Request{}
//          .      1.05s     49:           json.Unmarshal([]byte(req), reqObj)
//          .          .     50:           ret := ""
//          .          .     51:           for _, e := range reqObj.PayLoad {
//          .      510ms     52:                   ret += strconv.Itoa(e) + ","
//          .          .     53:           }
//          .          .     54:           repObj := &Response{reqObj.TransactionID, ret}
//          .       40ms     55:           repJson, err := json.Marshal(&repObj)
//          .          .     56:           if err != nil {
//          .          .     57:                   panic(err)
//          .          .     58:           }
//          .          .     59:           reps = append(reps, string(repJson))
//          .          .     60:   }
// (pprof) exit

func TestCreateRequest(t *testing.T) {
	str := createRequest()
	t.Log(str)
}

func TestProcessRequest(t *testing.T) {
	reqs := []string{}
	reqs = append(reqs, createRequest())
	reps := processRequest(reqs)
	t.Log(reps[0])
}

func BenchmarkProcessRequest(b *testing.B) {

	reqs := []string{}
	reqs = append(reqs, createRequest())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processRequest(reqs)
	}
	b.StopTimer()

}

func BenchmarkProcessRequestOld(b *testing.B) {
	reqs := []string{}
	reqs = append(reqs, createRequest())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processRequestOld(reqs)
	}
	b.StopTimer()

}
