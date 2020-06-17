package cache

import (
	"crypto/sha256"
	"fmt"

	"github.com/nikolasmelui/go-ecommerce-data-transfer-microservice/entity"
)

// Cachable ...
type Cachable interface {
	SetHash()
	GetHash() string
}

// ProductCache ...
type ProductCache struct {
	Data entity.Product
	Hash string
}

// SetHash makes the instance sha256-hash and set it in the instance Hash field
func (productCache *ProductCache) SetHash() {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%v", productCache.Data)))
	productCache.Hash = fmt.Sprintf("%x", hash.Sum(nil))
}

// GetHash returns the instance value of the hash field
func (productCache *ProductCache) GetHash() string {
	return productCache.Hash
}
