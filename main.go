package main

import (
	"fmt"

	mybigquery "github.com/ozaki-physics/gcp-training/myBigQuery"
	"github.com/ozaki-physics/gcp-training/myLineBot"
)

func main() {
	fmt.Println("hello world!")
	// helloworld.Main()
	// bigQuery_try()
}

func lineBot_try() {
	myLineBot.Main()
}

func bigQuery_try() {
	mybigquery.GetColumn()
	mybigquery.GetDatasetList()
	mybigquery.MakeDataset()
	mybigquery.InsertBatchCSV()
	mybigquery.InsertBatchJSON()
	mybigquery.ListTables()
	mybigquery.TableExists()
	mybigquery.MakeTable()
	mybigquery.InsertStreaming()
	mybigquery.ExportTableAsJSON()
}
