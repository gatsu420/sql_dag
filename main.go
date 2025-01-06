package main

import (
	"log"

	"github.com/gatsu420/sql_dag/dag"
)

func main() {
	dagMap, err := dag.ParseQuery("sample_query.sql")
	if err != nil {
		log.Fatalf("failed to parse query: %v", err)
	}

	err = dag.GenerateDAG(dagMap, "dag.dot", "dag.png")
	if err != nil {
		log.Fatalf("failed to generate DAG: %v", err)
	}
}
