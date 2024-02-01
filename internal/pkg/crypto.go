package pkg

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func RandString(max int) string {
	log.Println("sadasdasdas")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d", r.Intn(max+1))
}
