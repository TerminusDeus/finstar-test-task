package postgres

import (
	cRand "crypto/rand"
	"finstar-test-task/proto"
	"fmt"
	"github.com/infobloxopen/atlas-app-toolkit/errors"
	"gorm.io/gorm"
	"log"
	"math/big"
	"math/rand"
	"unsafe"
)

type seed struct {
	Name string
	Run  func(db *gorm.DB) error
}

func All(count int) []seed {
	var seeds []seed

	log.Printf("Will generate %d users", count)

	for i := 0; i < count; i++ {
		title := GenerateRandString(5, Letter)
		userId := uint64(rand.Int63n(1000000))
		balance := GetRandFloat(10, 10000)

		seeds = append(seeds, seed{
			Name: fmt.Sprintf("%s_%v_%v", title, userId, balance),
			Run: func(db *gorm.DB) error {
				return db.Create(&proto.UserORM{Id: userId, Balance: balance}).Error
			},
		})
	}
	return seeds
}

func generateUsers(database *gorm.DB, count int) {
	for _, seed := range All(count) {
		if err := seed.Run(database); err != nil {
			if errors.CondHasPrefix("ERROR: duplicate key value violates unique constraint \"users_pkey\"")(err) {
				generateUsers(database, 1)
			} else {
				log.Printf("Running seed '%s' failed with error: %s", seed.Name, err)
			}
		}
	}
}

func GenerateRandString(length int, t randContentType) string {
	var chars string
	switch t {
	case Letter:
		chars = letterBytes
	case Mixed:
		chars = mixedBytes
	case Digit:
		chars = digitBytes
	}
	b := make([]byte, length)
	for i, cache, remain := length-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(chars) {
			b[i] = chars[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

type randContentType int32

const (
	letterIdxBits                 = 6                    // 6 bits to represent a letter index
	letterIdxMask                 = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax                  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes                   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitBytes                    = "0123456789"
	mixedBytes                    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Letter        randContentType = 0
	Digit         randContentType = 1
	Mixed         randContentType = 2
)

const floatPrecision = 100

func GetRandInt(min, max int) int {
	nBig, _ := cRand.Int(cRand.Reader, big.NewInt(int64(max+1-min)))
	n := nBig.Int64()
	return int(n) + min
}

func GetRandFloat(min, max float32) float32 {
	minInt := int(min * floatPrecision)
	maxInt := int(max * floatPrecision)

	return float32(GetRandInt(minInt, maxInt)) / floatPrecision
}
