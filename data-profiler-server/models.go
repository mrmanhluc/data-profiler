package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Models

// Model is the interface returned by all the model functions
type Model interface{}

// ModelDataset is the struct that represents a set of datasets
type ModelDataset struct {
	ID          string
	Name        string
	Path        string
	Description string
	Files       []string
	Operators   []*ModelOperator
	Matrices    []*ModelSimilarityMatrix
	Estimators  []*ModelEstimator
}

// ModelOperator is the struct that represents an operator for a given dataset
type ModelOperator struct {
	ID          string
	Path        string
	Description string
}

// ModelEstimator represents an estimator object
type ModelEstimator struct {
	ID            string
	Path          string
	Filename      string
	Configuration map[string]string
}

// ModelSimilarityMatrix represents a similarity matrix
type ModelSimilarityMatrix struct {
	ID            string
	Path          string
	Filename      string
	Configuration map[string]string
}

// FUNCTIONS

// dbConnect is responsible to establish the connection with the DB backend.
// Written as separate function to increase modularity between the different
// DB backends.
func dbConnect() *sql.DB {
	db, err := sql.Open("sqlite3", Conf.Database)
	if err != nil {
		log.Println(err)
	}
	return db
}

func modelDatasetsList() []*ModelDataset {
	db := dbConnect()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM datasets")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	result := make([]*ModelDataset, 0)
	for rows.Next() {
		obj := new(ModelDataset)
		rows.Scan(&obj.ID, &obj.Path, &obj.Name, &obj.Description)
		result = append(result, obj)
	}
	return result
}

func modelDatasetGetInfo(id string) *ModelDataset {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM datasets WHERE id == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		obj := new(ModelDataset)
		rows.Scan(&obj.ID, &obj.Path, &obj.Name, &obj.Description)
		return obj
	}
	return nil
}

func modelDatasetGetFiles(id string) []string {
	var results []string
	path := modelDatasetGetInfo(id).Path
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
	}
	for _, f := range fs {
		if !f.IsDir() {
			results = append(results, f.Name())
		}
	}
	return results
}

func modelDatasetGetMatrices(id string) []*ModelSimilarityMatrix {
	db := dbConnect()
	defer db.Close()
	var results []*ModelSimilarityMatrix
	rows, err := db.Query("SELECT id, path, filename,configuration " +
		" FROM matrices WHERE datasetid == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		obj := new(ModelSimilarityMatrix)
		confString := ""
		rows.Scan(&obj.ID, &obj.Path, &obj.Filename, &confString)
		conf := make(map[string]string)
		err := json.Unmarshal([]byte(confString), &conf)
		if err != nil {
			log.Println(err)
		}
		obj.Configuration = conf
		results = append(results, obj)
	}
	return results
}

func modelDatasetGetEstimators(id string) []*ModelEstimator {
	db := dbConnect()
	defer db.Close()
	var results []*ModelEstimator
	rows, err := db.Query("SELECT id, path, filename,configuration " +
		" FROM estimators WHERE datasetid == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		obj := new(ModelEstimator)
		confString := ""
		rows.Scan(&obj.ID, &obj.Path, &obj.Filename, &confString)
		conf := make(map[string]string)
		err := json.Unmarshal([]byte(confString), &conf)
		if err != nil {
			log.Println(err)
		}
		obj.Configuration = conf
		results = append(results, obj)
	}
	return results
}

// modelSimilarityMatrixInsert inserts a new SM and returns the newly created Id
func modelSimilarityMatrixInsert(datasetID string, buffer []byte, conf map[string]string) string {
	dts := modelDatasetGetInfo(datasetID)
	filePath := writeBufferToFile(dts, "matrices", buffer)
	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare(
		"INSERT INTO matrices(path,filename,configuration,datasetid) " +
			"VALUES(?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
	}
	confString, err := json.Marshal(conf)
	if err != nil {
		log.Println(err)
	}
	res, err := stmt.Exec(filePath,
		path.Base(filePath),
		confString,
		dts.ID)
	if err != nil {
		log.Println(err)
	}
	resultInt, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%d", resultInt)
}

func modelSimilarityMatrixGet(id string) *ModelSimilarityMatrix {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query("SELECT id,path,filename,configuration" +
		" FROM matrices WHERE id == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		obj := new(ModelSimilarityMatrix)
		rows.Scan(&obj.ID, &obj.Path, &obj.Filename,
			&obj.Configuration)
		return obj
	}
	return nil
}

func modelEstimatorInsert(datasetID, matrixID string, buffer []byte, conf map[string]string) {
	dts := modelDatasetGetInfo(datasetID)
	filePath := writeBufferToFile(dts, "estimators", buffer)
	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare(
		"INSERT INTO estimators(path,filename,configuration,datasetid,matrixid) " +
			"VALUES(?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
	}
	confString, err := json.Marshal(conf)
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(filePath,
		path.Base(filePath),
		confString,
		dts.ID,
		matrixID)
	if err != nil {
		log.Println(err)
	}
}

func modelEstimatorGet(id string) *ModelEstimator {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query("SELECT id,path,filename,configuration" +
		" FROM estimator WHERE id == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		obj := new(ModelEstimator)
		rows.Scan(&obj.ID, &obj.Path, &obj.Filename,
			&obj.Configuration)
		return obj
	}
	return nil
}

func modelEstimatorGetBySM(id string) *ModelEstimator {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query("SELECT id,path,filename,configuration" +
		" FROM estimator WHERE matrixid == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		obj := new(ModelEstimator)
		rows.Scan(&obj.ID, &obj.Path, &obj.Filename,
			&obj.Configuration)
		return obj
	}
	return nil
}

// utility functions

func writeBufferToFile(dts *ModelDataset, prefix string, buffer []byte) string {
	dstDir := dts.Path + "/" + prefix
	err := os.Mkdir(dstDir, 0777)
	if err != nil {
		log.Println(err)
	}
	dstPath := dstDir + "/" + prefix + currentTimeSuffix() + ".mat"
	err = ioutil.WriteFile(dstPath, buffer, 0777)
	if err != nil {
		log.Println(err)
	}
	return dstPath
}

func currentTimeSuffix() string {
	t := time.Now()
	y, m, d := t.Year(), int(t.Month()), t.Day()
	h, min, sec := t.Hour(), t.Minute(), t.Second()
	return fmt.Sprintf("%d%d%d%d%d%d", y, m, d, h, min, sec)
}
