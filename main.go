package main

import (
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine"
	"math/rand"
	"time"
)

func execWorkflow(r *rand.Rand, bpmnEngine bpmn_engine.BpmnEngineState, process *bpmn_engine.ProcessInfo) int64 {
	variables := map[string]interface{}{}
	variables["value"] = r.Intn(256)
	startTime := time.Now().UnixNano()
	_, err := bpmnEngine.CreateAndRunInstance(process.ProcessKey, variables)
	endTime := time.Now().UnixNano()
	if err != nil {
		panic(err)
	}
	executionTimeInMillis := (endTime - startTime) / 1000
	return executionTimeInMillis
}

func justCompleteNothingHandler(job bpmn_engine.ActivatedJob) {
	job.Complete()
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	bpmnEngine := bpmn_engine.New()
	process, err := bpmnEngine.LoadFromFile("random-tree-workflow.bpmn")
	if err != nil {
		panic("file \"random-tree-workflow.bpmn\" can't be read.")
	}

	bpmnEngine.NewTaskHandler().Type("task").Handler(justCompleteNothingHandler)
	const noOfRuns = 100
	var results = make([]int64, noOfRuns)
	for i := 0; i < noOfRuns; i++ {
		executionTimeInMillis := execWorkflow(r, bpmnEngine, process)
		results[i] = executionTimeInMillis
	}

	var avgExecTime int64 = 0
	for i := 0; i < noOfRuns; i++ {
		avgExecTime += results[i]
	}
	avgExecTime = avgExecTime / int64(noOfRuns)
	println(fmt.Sprintf("Average execution of %d runs is %d ms", noOfRuns, avgExecTime))
}
