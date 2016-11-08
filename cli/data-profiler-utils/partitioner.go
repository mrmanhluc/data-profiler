package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/giagiannis/data-profiler/core"
)

type PartitionerParams struct {
	input         *string
	output        *string
	splits        *int
	partitionType *core.PartitionerType
}

func partitionerParseParams() *PartitionerParams {
	params := new(PartitionerParams)
	params.input = flag.String("i", "", "Input file to partition")
	params.output = flag.String("o", "", "Input file to partition")
	params.splits = flag.Int("c", 0, "Number of splits to create")
	part := flag.String("t", "UNIFORM", "Type of partitioning")
	if *part == "UNIFORM" {
		params.partitionType = new(core.PartitionerType)
		*params.partitionType = core.UNIFORM
	}
	flag.Parse()
	if *params.input == "" || *params.output == "" || *params.splits == 0 || *part == "" {
		fmt.Println("Please type -h to see usage")
		os.Exit(1)
	}
	return params
}

func partitionerRun() {
	params := partitionerParseParams()
	partitioner := core.NewDatasetPartitioner(
		*params.input,
		*params.output,
		*params.splits,
		*params.partitionType)
	partitioner.Partition()
	fmt.Println("Partitioning finished")
}
