package main

import "fmt"
import ("net/http")
import ("encoding/json")
import ("github.com/PuerkitoBio/goquery")
import ("github.com/olivere/elastic")

func create_client() (*elastic.Client) {
	// Obtain a client. You can provide your own HTTP client here.
	es_client, err := elastic.NewClient(http.DefaultClient)
	if err != nil {
		// Handle error
		panic(err)
	}
	return es_client
}

func check_es_server(client *elastic.Client) {
	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping().Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s \n", code, info.Version.Number)
}

func main() {
	fmt.Printf("Hello, we are searching.\n")
	es_client := create_client()
	check_es_server(es_client)
	create_index(es_client)
	ExampleScrape()
	//index_flight_info(es_client)
}

//
// can't scrape google flights, need to use their API
//
func ExampleScrape() {
	doc, err := goquery.NewDocument("https://www.google.com/flights/#search;f=DEN;t=SGN;d=2015-02-17;r=2015-02-21;mc=m") 
	if err != nil {
		panic(err)
	}
	fmt.Printf("got doc: %s", doc.Html())
	
	doc.Find("._Wm").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("text %s\n", s.Text())
		city := s.Find("_v2").Text()
		price := s.Find("._Fs").Text()
		fmt.Printf("Review %d: %s - %s\n", i, city, price)
	})
}

type flight struct {
	Route string
	Info string
	Price float64
}

func index_flight_info(client *elastic.Client) {
	// Index a flight (using JSON serialization)
	flight1 := &flight{Route: "denver-vietnam", Info: "Take Five", Price: 700.00}
	flight_json, err := json.Marshal(flight1)
	fmt.Printf("json %s \n", string(flight_json))
	put1, err := client.Index().
		Index("flight_costs").
		Type("flight_info").
		BodyString(string(flight_json)).
		Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed flight %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}

func create_index(client *elastic.Client) {
	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("flight_costs").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("flight_costs").Do()
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		} else {
			fmt.Printf("index created.\n")
		}
	}
}
