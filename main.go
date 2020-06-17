package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/nikolasmelui/go-ecommerce-data-transfer-microservice/cache"
	"github.com/nikolasmelui/go-ecommerce-data-transfer-microservice/cconfig"
	"github.com/nikolasmelui/go-ecommerce-data-transfer-microservice/source"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	req, err := http.NewRequest("GET", source.ProductsURL, nil)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Accept", "application/xml")
	req.SetBasicAuth(cconfig.Config.SourceBasicAuthLogin, cconfig.Config.SourceBasicAuthPassword)

	client := &http.Client{
		Timeout: time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			log.Fatal(errors.New(errRes.Message))
		}
		log.Fatal(fmt.Errorf("Unknown error, status code: %d", res.StatusCode))
	}

	body, _ := ioutil.ReadAll(res.Body)
	var productsResponse source.ProductsResponse
	err = xml.Unmarshal(body, &productsResponse)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}

	cacher := cache.NewRedisConnection(cconfig.Config.RedisHost, cconfig.Config.RedisPassword, cconfig.Config.RedisDB, 120)

	products := productsResponse.Products

	for i, product := range products {
		fmt.Printf("%d ----------\n", i)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		var productCache cache.Cachable
		productCache = &cache.ProductCache{}
		err = cacher.Get(ctx, product.ID, &productCache)
		if err != nil {
			log.Fatalf("Error: %v", err.Error())
		}

		var productToCache cache.Cachable
		productToCache = &cache.ProductCache{
			Data: product,
			Hash: "",
		}
		productToCache.SetHash()

		productCacheHash := productCache.GetHash()
		productToCacheHash := productToCache.GetHash()

		// Chech the instance is not in the cache or hashes is different
		if (len(productCacheHash) < 1) || (productCacheHash != productToCacheHash) {
			err := cacher.Set(ctx, product.ID, &productToCache)
			if err != nil {
				log.Fatalf("Error: %v", err.Error())
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}
