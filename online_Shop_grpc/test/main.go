package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

func main() {
	host := "http://192.168.220.128:9200"
	client, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	q := elastic.NewMatchQuery("company.address", "shenzhen")

	do, err := client.Search().Index("account").Query(q).Do(context.Background())
	if err != nil {
		panic(err)
	}

	var res []byte
	res, err = do.Hits.Hits[0].Source.MarshalJSON()
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	fmt.Println(string(res))
}
