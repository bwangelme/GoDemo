package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/olivere/elastic/v7"
)

const (
	ES_URL     = "http://192.168.56.23:9200"
	INDEX_NAME = "movie"
)

type Movie struct {
	ID          int64   `json:"id"`
	Year        int64   `json:"year"`
	Score       float64 `json:"score"`
	RatingCount int64   `json:"rating_count"`
}

func getMovies(filename string) []Movie {
	fd, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	var result []Movie
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		var m Movie
		err := json.Unmarshal([]byte(line), &m)
		if err != nil {
			log.Println(err)
			log.Println(line)
			continue
		}
		result = append(result, m)
	}

	return result
}

func main() {
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(ES_URL),
	)
	if err != nil {
		log.Fatal(err)
	}
	res, code, err := client.Ping(ES_URL).Do(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res.Version.Number, code)

	movies := getMovies("data/top100.txt")
	fmt.Println(len(movies))

	chunkSize := 50
	for i := 0; i < len(movies); i += chunkSize {
		end := i + chunkSize
		bulkRequest := client.Bulk()
		for j := i; j < end && j < len(movies); j++ {
			m := movies[j]
			req := elastic.NewBulkIndexRequest().Index(INDEX_NAME).Id(strconv.Itoa(int(m.ID))).Doc(m)
			bulkRequest.Add(req)
		}
		resp, err := bulkRequest.Do(ctx)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(resp.Took)
	}

}
