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
	workflow := "random-tree-workflow.bpmn"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bpmnEngine := bpmn_engine.New()
	process, err := bpmnEngine.LoadFromFile(workflow)
	if err != nil {
		panic("file '" + workflow + "' can't be read.")
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
	println(fmt.Sprintf("Workflow: '%s'\n"+
		"Executing %d times\n"+
		"Average time spent per execution is %d ms", workflow, noOfRuns, avgExecTime))
}
