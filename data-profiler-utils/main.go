package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var utilsDescription = map[string]string{
	"help":          "prints this help message",
	"partition":     "split a dataset file into more smaller files",
	"heatmap":       "creates a heatmap of the datasets according to their accuracy",
	"similarities":  "calculates and stores the similarity matrix of the specified datasets",
	"train":         "exhaustively trains the specified ML job with the datasets and creates a scores matrix",
	"clustering":    "clusters the datasets based on the similarity matrix and prints their accuracy vs their cluster",
	"simcomparison": "compares a list of similarity matrices",
	"mds":           "executes Multidimensional Scaling to a similarity matrix",

	// experiments
	"exp-accuracy":       "trains an ML model and prints the error",
	"exp-ordering":       "compares the ordering of the datasets to the original ordering",
	"exp-online-indexer": "executes the online indexer experiment",

	// utils
	"print-utils": "utilility to print objects in binary form",
}

func main() {
	// consume the first command
	rand.Seed(int64(time.Now().Nanosecond()))
	if len(os.Args) < 2 || os.Args[1] == "help" {
		fmt.Fprintf(os.Stderr, "Usage: %s [command]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "List of commands:")
		for name, description := range utilsDescription {
			fmt.Fprintf(os.Stderr, "\t%s - %s\n", name, description)
		}
		os.Exit(1)
	}
	command := os.Args[1]
	os.Args = os.Args[1:]

	if command == "partition" {
		partitionerRun()
	} else if command == "heatmap" {
		heatmapRun()
	} else if command == "similarities" {
		similaritiesRun()
	} else if command == "train" {
		trainRun()
	} else if command == "clustering" {
		clusteringRun()
	} else if command == "simcomparison" {
		simcomparisonRun()
	} else if command == "mds" {
		mdsRun()
	} else if command == "indexing" {
		indexingRun()
	} else if len(command) > 3 && command[0:4] == "exp-" {
		experiment := command[4:len(command)]
		if experiment == "accuracy" {
			expAccuracyRun()
		} else if experiment == "ordering" {
			expOrderingRun()
		} else if experiment == "online-indexer" {
			expOnlineIndexerRun()
		}
	} else if command == "print-utils" {
		printUtilsRun()
	} else {
		fmt.Fprintln(os.Stderr, "Command not identified")
	}

}
