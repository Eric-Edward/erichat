package middlewares

import (
	"EcChat/utils"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestAuthorization(t *testing.T) {
	_, _ = utils.GenerateJWT("0bc956e6-6a6c-4f35-aa9f-4cfaf0d42a53", time.Now().Add(time.Minute*20))
	fmt.Println("jmt")
	_, _ = http.Get("http://localhost:8081/complete?uuid=0bc956e6-6a6c-4f35-aa9f-4cfaf0d42a53")

}
