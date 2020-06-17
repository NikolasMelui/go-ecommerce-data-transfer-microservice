package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nikolasmelui/go-ecommerce-data-transfer-microservice/cache"
	"github.com/nikolasmelui/go-ecommerce-data-transfer-microservice/cconfig"
	"github.com/nikolasmelui/go-ecommerce-data-transfer-microservice/source"
)

func main() {

	sourceClient := source.NewClient()

	var productsResponse source.ProductsResponse

	err := sourceClient.GetData("/products", &productsResponse)
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
