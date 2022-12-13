package mybigquery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// BigQuery のリソース
type resource struct {
	ProjectID  string `json:"project_id"`
	DatasetID  string `json:"dataset_id"`
	TableID01  string `json:"dable_id_01"`
	TableID02  string `json:"dable_id_02"`
	TableID03  string `json:"dable_id_03"`
	TableID04  string `json:"dable_id_04"`
	BucketName string `json:"bucket_name"`
}

func newResource() (resource, error) {
	bytes, err := os.ReadFile("./json/keypath_key.json")
	if err != nil {
		log.Fatalln(err)
	}
	var r resource
	if err := json.Unmarshal(bytes, &r); err != nil {
		log.Fatalln(err)
	}
	return r, nil
}

// BigQuery のクライアントを作成する
// アクセストークンは json から読み込む
func newBigQueryClient(ctx context.Context) (*bigquery.Client, error) {
	// key.json への path が書かれた json を読み込む
	bytes, err := os.ReadFile("./json/keypath_key.json")
	if err != nil {
		log.Fatalln(err)
	}
	type path struct {
		BigQuery  string `json:"big_query"`
		ProjectID string `json:"project_id"`
	}
	var p path
	if err := json.Unmarshal(bytes, &p); err != nil {
		log.Fatalln(err)
	}
	key := p.BigQuery
	projectID := p.ProjectID

	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(key))
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}

	return client, nil
}

// 取得
func GetColumn() {
	// BigQuery へ接続
	ctx := context.Background()
	client, err := newBigQueryClient(ctx)
	defer client.Close()

	rows, err := query(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	if err := printResults(os.Stdout, rows); err != nil {
		log.Fatal(err)
	}
}

// 取得 SQL
func query(ctx context.Context, client *bigquery.Client) (*bigquery.RowIterator, error) {
	query := client.Query(
		`
    SELECT
      name,
      count
    FROM
      ` + "`babynames.names_2014`" + `
    WHERE
      gender = 'M'
    ORDER BY
      count DESC
    LIMIT
      5
    `,
	)
	return query.Read(ctx)
}

// 取得 した レコード の構造体
type BabyNamesRow struct {
	Name  string `bigquery:"name"`
	Count int    `bigquery:"count"`
}

// 取得 した レコード の出力
func printResults(w io.Writer, iter *bigquery.RowIterator) error {
	for {
		var row BabyNamesRow
		err := iter.Next(&row)
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error iterating through results: %v", err)
		}

		fmt.Fprintf(w, "name: %s, count: %d\n", row.Name, row.Count)
	}
}

// 書き込み(バッチ)
func InsertBatchCSV() {
	r, _ := newResource()
	datasetID := r.DatasetID
	tableID := r.TableID01
	filename := "./myBigQuery/test.csv"
	err := importCSVFromFile(datasetID, tableID, filename)
	if err != nil {
		log.Fatalln(err)
	}
}

// 書き込み(バッチ) の 詳細(CSV)
func importCSVFromFile(datasetID, tableID, filename string) error {
	// BigQuery へ接続
	ctx := context.Background()
	client, err := newBigQueryClient(ctx)
	defer client.Close()

	// ファイル操作
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	source := bigquery.NewReaderSource(f)
	source.AutoDetect = true // Allow BigQuery to determine schema.
	// source.SkipLeadingRows = 1 // CSV has a single header line.

	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(source)
	// WriteTruncate(全書き換え)で書き込みする
	// loader.LoadConfig.WriteDisposition = bigquery.WriteTruncate

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

// 書き込み(バッチ)
func InsertBatchJSON() {
	r, _ := newResource()
	datasetID := r.DatasetID
	tableID := r.TableID03
	filename := "./myBigQuery/test.json"
	err := importJSONFromFile(datasetID, tableID, filename)
	if err != nil {
		log.Fatalln(err)
	}
}

// 書き込み(バッチ) の 詳細(JSON)
func importJSONFromFile(datasetID, tableID, filename string) error {
	// BigQuery へ接続
	ctx := context.Background()
	client, err := newBigQueryClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// ファイル操作
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	source := bigquery.NewReaderSource(f)
	source.SourceFormat = bigquery.JSON
	// source.AutoDetect = true  // Allow BigQuery to determine schema.
	source.Schema = bigquery.Schema{
		{Name: "title03", Type: bigquery.StringFieldType},
		{Name: "title04", Type: bigquery.StringFieldType},
	}

	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(source)
	// WriteTruncate(全書き換え)で書き込みする
	// loader.LoadConfig.WriteDisposition = bigquery.WriteTruncate

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

// json ファイルから テーブルを作成
func MakeTableFromJSON() {
	r, _ := newResource()
	datasetID := r.DatasetID
	tableID := r.TableID04
	filename := "./myBigQuery/test.json"
	err := makeJSONFromFile(datasetID, tableID, filename)
	if err != nil {
		log.Fatalln(err)
	}
}

// json ファイルから テーブルを作成 の詳細
func makeJSONFromFile(datasetID, tableID, filename string) error {
	// BigQuery へ接続
	ctx := context.Background()
	client, err := newBigQueryClient(ctx)
	defer client.Close()

	// ファイル操作
	f, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	schema, err := bigquery.SchemaFromJSON(f)
	metaData := &bigquery.TableMetadata{Schema: schema}

	loader := client.Dataset(datasetID).Table(tableID)
	if err := loader.Create(ctx, metaData); err != nil {
		return err
	}

	return nil
}

// 書き込み(ストリーミング)
func InsertStreaming() {
	r, _ := newResource()
	datasetID := r.DatasetID
	tableID := r.TableID02
	err := insertRows(datasetID, tableID)
	if err != nil {
		log.Fatalln(err)
	}
}

// 書き込み(ストリーミング) の詳細
func insertRows(datasetID, tableID string) error {
	// BigQuery へ接続
	ctx := context.Background()
	client, err := newBigQueryClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	inserter := client.Dataset(datasetID).Table(tableID).Inserter()
	items := []*Item{
		{Name: "hello", Age: 26},
		{Name: "world", Age: 18},
	}

	if err := inserter.Put(ctx, items); err != nil {
		return err
	}
	return nil
}

// ストリーミング インサートの行
type Item struct {
	Name string
	Age  int
}

func (i *Item) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"full_name": i.Name,
		"age":       i.Age,
	}, bigquery.NoDedupeID, nil
}

// データセット の確認
func GetDatasetList() {
	err := getDatasetList()
	if err != nil {
		log.Fatalln(err)
	}
}

// データセット の一覧
func getDatasetList() error {
	// BigQuery へ接続
	ctx := context.Background()
	client, err := newBigQueryClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	it := client.Datasets(ctx)
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(dataset.DatasetID)
	}
	return nil
}

// データセット の作成
func MakeDataset() {
	r, _ := newResource()
	datasetID := r.DatasetID
	err := createDataset(datasetID)
	if err != nil {
		log.Fatalln(err)
	}
}

// データセット の作成 の詳細
func createDataset(datasetID string) error {
	// BigQuery へ接続
	ctx := context.Background()
	client, err := newBigQueryClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	meta := &bigquery.DatasetMetadata{
		// 前にブラウザから作った dataset(babynames) が US だったから US でいいかな
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}
	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {
		return err
	}
	return nil
}

// テーブル の一覧
func ListTables() {
	r, _ := newResource()
	datasetID := r.DatasetID
	err := listTables(datasetID)
	if err != nil {
		log.Fatalln(err)
	}
}

// テーブル の一覧 の詳細
func listTables(datasetID string) error {
	// BigQuery へ接続
	ctx := context.Background()
	client, err := newBigQueryClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	ts := client.Dataset(datasetID).Tables(ctx)
	for {
		t, err := ts.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Table: %q\n", t.TableID)
	}
	return nil
}

// テーブル の確認
func TableExists() {
	r, _ := newResource()
	datasetID := r.DatasetID
	tableID := r.TableID01
	// tableID := r.TableID02
	err := tableExists(datasetID, tableID)
	if err != nil {
		log.Fatalln(err)
	}
}

// テーブル の確認 の詳細
func tableExists(datasetID, tableID string) error {
	// BigQuery へ接続
	ctx := context.Background()
	client, err := newBigQueryClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	tableRef := client.Dataset(datasetID).Table(tableID)
	if _, err = tableRef.Metadata(ctx); err != nil {
		return err
		// if e, ok := err.(*googleapi.Error); ok {
		// 	if e.Code == http.StatusNotFound {
		// 		return errors.New("dataset or table not found")
		// 	}
		// }
	}
	fmt.Printf("%s: 存在します\n", tableID)
	return nil
}

// テーブル の作成
func MakeTable() {
	r, _ := newResource()
	datasetID := r.DatasetID
	tableID := r.TableID02
	err := createTable(datasetID, tableID)
	if err != nil {
		log.Fatalln(err)
	}
}

// テーブル の作成 の詳細
func createTable(datasetID, tableID string) error {
	// BigQuery へ接続
	ctx := context.Background()
	client, err := newBigQueryClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	sampleSchema := bigquery.Schema{
		{Name: "full_name", Type: bigquery.StringFieldType},
		{Name: "age", Type: bigquery.IntegerFieldType},
	}
	meta := &bigquery.TableMetadata{
		Schema: sampleSchema,
	}

	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, meta); err != nil {
		return err
	}
	return nil
}

// テーブルを JSON にエクスポート
func ExportTableAsJSON() {
	r, _ := newResource()
	projectID := r.ProjectID
	datasetID := r.DatasetID
	tableID := r.TableID03
	bucketName := r.BucketName
	gcsURI := "gs://" + bucketName + "/testBigQuery.json"
	err := exportTableAsJSON(projectID, datasetID, tableID, gcsURI)
	if err != nil {
		log.Fatalln(err)
	}
}

func exportTableAsJSON(srcProject, srcDataset, srcTable, gcsURI string) error {
	// BigQuery へ接続
	ctx := context.Background()
	client, err := newBigQueryClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	gcsRef := bigquery.NewGCSReference(gcsURI)
	gcsRef.DestinationFormat = bigquery.JSON

	extractor := client.DatasetInProject(srcProject, srcDataset).Table(srcTable).ExtractorTo(gcsRef)
	extractor.Location = "US"

	job, err := extractor.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}
