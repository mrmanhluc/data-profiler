package core

import (
	"math/rand"
	"os"
	"testing"
)

func TestScriptBasedModelerConfiguration(t *testing.T) {
	m := new(ScriptBasedModeler)
	err := m.Configure(map[string]string{"foo": "bar"})
	if err == nil {
		t.Log("Should have returned an error!")
		t.FailNow()
	}
	err = m.Configure(map[string]string{"script": mlScript})
	if err != nil {
		t.Log("Should have not returned an error!")
		t.FailNow()
	}
}

func TestScriptBasedModelerRun(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 50, 3)
	m, err := createModeler(datasets)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	m.Configure(map[string]string{"script": mlScriptAppx})
	err = m.Run()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if m.Samples() == nil || len(m.Samples()) == 0 {
		t.Log("Should have created samples")
		t.Fail()
	}
	if len(m.AppxValues()) != len(m.Datasets()) {
		t.Log("Wrong number of approximated values")
		t.Fail()
	}

	cleanDatasets(datasets)

}

func TestErrorMetrics(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 50, 3)
	m, err := createModeler(datasets)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	m.Configure(map[string]string{"script": mlScriptAppx})
	if m.ErrorMetrics() != nil {
		t.Log("ErrorMetrics should have been nil")
		t.Fail()
	}
	err = m.Run()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	metrics := m.ErrorMetrics()
	if metrics["MSE-all"] < 0 || metrics["R^2-all"] > 1.0 || metrics["R^2-all"] < 0.0 {
		t.Log("Wrong metric values")
		t.Fail()
	}
	cleanDatasets(datasets)
}

func TestScriptBasedModelerMissingDatasets(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 50, 3)
	noDeletedDatasets := 47
	m, err := createModeler(datasets)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	permutation := rand.Perm(len(datasets))
	for i := 0; i < noDeletedDatasets; i++ {
		path := datasets[permutation[i]].Path()
		os.Remove(path)
		t.Log("Removing", path)
	}

	m.Configure(map[string]string{"script": mlScriptAppx})
	err = m.Run()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if m.Samples() == nil || len(m.Samples()) == 0 {
		t.Log("Should have created samples")
		t.Fail()
	}
	if len(m.AppxValues()) != len(m.Datasets()) {
		t.Log("Wrong number of approximated values")
		t.Fail()
	}
	if len(m.Samples()) != (50 - noDeletedDatasets) {
		t.Log("Expected", noDeletedDatasets, "samples and got", len(m.Samples()))
		t.Fail()
	}

	cleanDatasets(datasets)
}

func createModeler(datasets []*Dataset) (Modeler, error) {
	e := NewDatasetSimilarityEstimator(SimilarityTypeBhattacharyya, datasets)
	e.Configure(map[string]string{"partitions": "256"})
	e.Compute()
	sm := e.SimilarityMatrix()

	mds := NewMDScaling(sm, 3, mdsScript)
	mds.Compute()
	coords := mds.Coordinates()

	eval, err := NewDatasetEvaluator(OnlineEval, map[string]string{
		"script":  operatorScript,
		"testset": "",
	})
	if err != nil {
		return nil, err
	}

	m := NewModeler(datasets, 0.2, coords, eval)
	return m, nil
}
