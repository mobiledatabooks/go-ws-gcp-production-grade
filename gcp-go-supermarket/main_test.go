package main

// MIT License

// Copyright (c) 2022 Mobile Data Books, LLC

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/v3/assert"
)

// go vet .
// go test -v -cover ./...

// GOMAXPROCS=8 CGO_ENABLED=1 go test -count 1000 -race -failfast  -v

// export SDKROOT="$(xcrun --sdk macosx --show-sdk-path)"
// CGO_ENABLED=1 go test -race
// CGO_ENABLED=1 go test -race -failfast
// GOMAXPROCS=8 CGO_ENABLED=1 go test -count 100 -race -v

// go test -coverprofile=c.out
// go tool cover -func=c.out
// go tool cover -html=c.out

//
// -count n
//     Run each test, benchmark, and fuzz seed n times (default 1).
//     If -cpu is set, run n times for each GOMAXPROCS value.
//     Examples are always run once. -count does not apply to
//     fuzz tests matched by -fuzz.
//
// -cpu 1,2,4
//     Specify a list of GOMAXPROCS values for which the tests, benchmarks or
//     fuzz tests should be executed. The default is the current value
//     of GOMAXPROCS. -cpu does not apply to fuzz tests matched by -fuzz.

// go test -run TestPing -v
// CGO_ENABLED=1 go test -race -run TestItems -v
// CGO_ENABLED=1 go test -race -run TestAddCheckProduceCode -v
// CGO_ENABLED=1 go test -race -run TestAddCheckName -v
// CGO_ENABLED=1 go test -race -run TestAddCheckUnitPrice -v
// GOMAXPROCS=8 CGO_ENABLED=1 go test -count 100 -race -failfast -run TestAddSingleRecord -v
// CGO_ENABLED=1 go test -race -run TestAddMultipleRecords -v
// CGO_ENABLED=1 go test -race -run TestGetItem -v
// CGO_ENABLED=1 go test -race -run TestDelete -v
//
// https://pkg.go.dev/cmd/go#hdr-Testing_flags
// https://go.dev/src/net/http/status.go
// https://go.dev/blog/race-detector
// https://go.dev/doc/articles/race_detector
//
// sysctl -n hw.ncpu
// 8
// GOMAXPROCS=8

func routerGETReq(method, path string, router *gin.Engine) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	router.ServeHTTP(w, req)
	// log.Printf("%d - %s", w.Code, w.Body.String())
	return w
}
func routerPOSTReq(method, path string, jsonData []byte, router *gin.Engine) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	router.ServeHTTP(w, req)
	// log.Printf("%d - %s", w.Code, w.Body.String())
	return w
}

// go test -run TestPing -v

func TestPing(t *testing.T) {
	db := database{}
	router := db.dbInit()
	tests := map[string]struct {
		method     string
		path       string
		wantCode   int
		wantResult string
	}{
		"simple": {method: "GET", path: "/api/v1/ping", wantCode: 200, wantResult: "pong"},
		"error":  {method: "GET", path: "/ping1", wantCode: 405, wantResult: `{"error":"endpoint not found"}`},
	}
	for name, tc := range tests {
		got := routerGETReq(tc.method, tc.path, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}
}

// go test -run TestItems -v

func TestItems(t *testing.T) {
	db := database{}
	router := db.dbInit()

	testsItems := map[string]struct {
		method     string
		path       string
		wantCode   int
		wantResult string
	}{
		"simple":     {method: "GET", path: "/api/v1/items", wantCode: 200, wantResult: `[{"code":"A12T-4GH7-QPL9-3N4M","name":"Lettuce","price":"$3.41"},{"code":"E5T6-9UI3-TH15-QR88","name":"Peach","price":"$2.99"},{"code":"TQ4C-VV6T-75ZX-1RMR","name":"Gala Apple","price":"$3.59"},{"code":"YRT6-72AS-K736-L4AR","name":"Green Pepper","price":"$0.79"}]`},
		"wrong path": {method: "GET", path: "/items1", wantCode: 405, wantResult: `{"error":"endpoint not found"}`},
	}
	for name, tc := range testsItems {
		got := routerGETReq(tc.method, tc.path, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}
}

// go test -run TestAddCheckProduceCode -v

func TestAddCheckProduceCode(t *testing.T) {
	db := database{}
	router := db.dbInit() //

	tests := map[string]struct {
		method     string
		path       string
		jsonData   []byte
		wantCode   int
		wantResult string
	}{
		"bad request1: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL9-3N4M1", "name": "Lettuce", "price": "3.41"}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.ProduceCode' Error:Field validation for 'ProduceCode' failed on the 'isproducecode' tag"}`},
		"bad request2: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL93N4M", "name": "Lettuce", "price": "3.41"}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.ProduceCode' Error:Field validation for 'ProduceCode' failed on the 'isproducecode' tag"}`},
		"ok request: code":   {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code":"ZRT6-72AS-K736-L4AZ","name":"Greener Pepper","price":"9.99"}]`), wantCode: 201, wantResult: `{"status":"item added"}`},
	}
	for name, tc := range tests {
		got := routerPOSTReq(tc.method, tc.path, tc.jsonData, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}

	testsItems := map[string]struct {
		method     string
		path       string
		wantCode   int
		wantResult string
	}{
		"validate": {method: "GET", path: "/api/v1/items", wantCode: 200, wantResult: `[{"code":"A12T-4GH7-QPL9-3N4M","name":"Lettuce","price":"$3.41"},{"code":"E5T6-9UI3-TH15-QR88","name":"Peach","price":"$2.99"},{"code":"TQ4C-VV6T-75ZX-1RMR","name":"Gala Apple","price":"$3.59"},{"code":"YRT6-72AS-K736-L4AR","name":"Green Pepper","price":"$0.79"},{"code":"ZRT6-72AS-K736-L4AZ","name":"Greener Pepper","price":"$9.99"}]`},
	}
	for name, tc := range testsItems {
		got := routerGETReq(tc.method, tc.path, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}
	tests = map[string]struct {
		method     string
		path       string
		jsonData   []byte
		wantCode   int
		wantResult string
	}{
		"bad request: code":              {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "", "name": "Lettuce", "price": "3.41"}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.ProduceCode' Error:Field validation for 'ProduceCode' failed on the 'required' tag"}`},
		"bad request: code, name":        {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "", "name": "", "price": "3.41"}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.ProduceCode' Error:Field validation for 'ProduceCode' failed on the 'required' tag\nKey: 'Item.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`},
		"bad request: code, name, price": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "", "name": "", "price": ""}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.ProduceCode' Error:Field validation for 'ProduceCode' failed on the 'required' tag\nKey: 'Item.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'Item.UnitPrice' Error:Field validation for 'UnitPrice' failed on the 'required' tag"}`},
	}
	for name, tc := range tests {
		got := routerPOSTReq(tc.method, tc.path, tc.jsonData, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}
}

// go test -run TestAddCheckName -v
func TestAddCheckName(t *testing.T) {
	db := database{}
	router := db.dbInit() //

	tests := map[string]struct {
		method     string
		path       string
		jsonData   []byte
		wantCode   int
		wantResult string
	}{
		"bad request1: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL9-3N4M", "name": "Lettuce-", "price": "3.41"}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.Name' Error:Field validation for 'Name' failed on the 'alphanumandspace' tag"}`},
		"bad request2: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "ZRT6-72AS-K736-L4AZ", "name": "Greener Pepper", "price": "9.99"}]`), wantCode: 201, wantResult: `{"status":"item added"}`},
	}
	for name, tc := range tests {
		got := routerPOSTReq(tc.method, tc.path, tc.jsonData, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}
	testsItems := map[string]struct {
		method     string
		path       string
		wantCode   int
		wantResult string
	}{
		"validate": {method: "GET", path: "/api/v1/items", wantCode: 200, wantResult: `[{"code":"A12T-4GH7-QPL9-3N4M","name":"Lettuce","price":"$3.41"},{"code":"E5T6-9UI3-TH15-QR88","name":"Peach","price":"$2.99"},{"code":"TQ4C-VV6T-75ZX-1RMR","name":"Gala Apple","price":"$3.59"},{"code":"YRT6-72AS-K736-L4AR","name":"Green Pepper","price":"$0.79"},{"code":"ZRT6-72AS-K736-L4AZ","name":"Greener Pepper","price":"$9.99"}]`},
	}
	for name, tc := range testsItems {
		got := routerGETReq(tc.method, tc.path, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}
}

// go test -run TestAddCheckUnitPrice -v
func TestAddCheckUnitPrice(t *testing.T) {
	db := database{}
	router := db.dbInit() //

	tests := map[string]struct {
		method     string
		path       string
		jsonData   []byte
		wantCode   int
		wantResult string
	}{
		"bad request1: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL9-3N4M", "name": "Lettuce", "price": "3.41-"}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.UnitPrice' Error:Field validation for 'UnitPrice' failed on the 'isunitprice' tag"}`},
		"ok request1: code":  {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "X12T-4GH7-QPL9-3N4X", "name": "Lettuces", "price": "9.41"}]`), wantCode: 201, wantResult: `{"status":"item added"}`},
		"bad request2: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "X12T-4GH7-QPL9-3N4X", "name": "Lettuces", "price": "9.411"}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.UnitPrice' Error:Field validation for 'UnitPrice' failed on the 'isunitprice' tag"}`},
		"bad request3: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "X12T-4GH7-QPL9-3N4X", "name": "Lettuces", "price": "9411"}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.UnitPrice' Error:Field validation for 'UnitPrice' failed on the 'isunitprice' tag"}`},
		"ok request2: code":  {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL9-3N4A", "name": "Lettuces1", "price": "9.4"}]`), wantCode: 201, wantResult: `{"status":"item added"}`},
		"bad request4: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL9-3N4A", "name": "Lettuces1", "price": "9."}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.UnitPrice' Error:Field validation for 'UnitPrice' failed on the 'isunitprice' tag"}`},
		"bad request5: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL9-3N4A", "name": "Lettuces1", "price": "9"}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.UnitPrice' Error:Field validation for 'UnitPrice' failed on the 'isunitprice' tag"}`},
		"bad request6: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL9-3N4A", "name": "Lettuces1", "price": ""}]`), wantCode: 400, wantResult: `{"error":"[0]: Key: 'Item.UnitPrice' Error:Field validation for 'UnitPrice' failed on the 'required' tag"}`},
		"bad request7: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL9-3N4A", "name": "Lettuces1", "price": 9.41}]`), wantCode: 400, wantResult: `{"error":"json: cannot unmarshal number into Go struct field Item.price of type string"}`},
	}
	for name, tc := range tests {
		got := routerPOSTReq(tc.method, tc.path, tc.jsonData, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}

}

// CGO_ENABLED=1 go test -race -run TestAddSingleRecord -v

func TestAddSingleRecord(t *testing.T) {
	db := database{}
	router := db.dbInit()
	//

	tests := map[string]struct {
		method     string
		path       string
		jsonData   []byte
		wantCode   int
		wantResult string
	}{
		"ok request1: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL9-3N4M", "name": "Lettuce", "price": "3.41"}]`), wantCode: 200, wantResult: `{"status":"item exist, not added"}`},
		"ok request2: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code":"A12T-4GH7-QPL9-3N4M","name":"Lettuce","price":"3.41"},{"code":"E5T6-9UI3-TH15-QR88","name":"Peach","price":"2.99"},{"code":"TQ4C-VV6T-75ZX-1RMR","name":"Gala Apple","price":"3.59"}]`), wantCode: 200, wantResult: `{"status":"item exist, not added"}`},
		"ok request3: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code":"ZRT6-72AS-K736-L4AZ","name":"Greener Pepper","price":"9.99"}]`), wantCode: 201, wantResult: `{"status":"item added"}`},
		// "bad request1: code": {method: "POST", path: "/add", jsonData: []byte(`[{"code":"ZRT6-72AS-K736-L4AZ","name":"Greener Pepper","price":"9.99"}]`), wantCode: 200, wantResult: `{"status":"item exist, not added"}`},
	}
	for name, tc := range tests {
		got := routerPOSTReq(tc.method, tc.path, tc.jsonData, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}
	//

	testsItems := map[string]struct {
		method     string
		path       string
		wantCode   int
		wantResult string
	}{
		"simple": {method: "GET", path: "/api/v1/items", wantCode: 200, wantResult: `[{"code":"A12T-4GH7-QPL9-3N4M","name":"Lettuce","price":"$3.41"},{"code":"E5T6-9UI3-TH15-QR88","name":"Peach","price":"$2.99"},{"code":"TQ4C-VV6T-75ZX-1RMR","name":"Gala Apple","price":"$3.59"},{"code":"YRT6-72AS-K736-L4AR","name":"Green Pepper","price":"$0.79"},{"code":"ZRT6-72AS-K736-L4AZ","name":"Greener Pepper","price":"$9.99"}]`},
	}
	for name, tc := range testsItems {
		got := routerGETReq(tc.method, tc.path, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}
	w := httptest.NewRecorder()
	// "bad request1: code": {method: "POST", path: "/add", jsonData: []byte(`[{"code":"ZRT6-72AS-K736-L4AZ","name":"Greener Pepper","price":"9.99"}]`), wantCode: 200, wantResult: `{"status":"item exist, not added"}`},

	jsonData := []byte(`[{"code":"ZRT6-72AS-K736-L4AZ","name":"Greener Pepper","price":"9.99"}]`)
	req, _ := http.NewRequest("POST", "/api/v1/add", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	router.ServeHTTP(w, req)
	// log.Printf("%d - %s", w.Code, w.Body.String())

	// ok request
	assert.Equal(t, 200, w.Code)
	res := `{"status":"item exist, not added"}`
	assert.Equal(t, res, w.Body.String())
}
func TestAddMultipleRecords(t *testing.T) {
	db := database{}
	router := db.dbInit() //

	tests := map[string]struct {
		method     string
		path       string
		jsonData   []byte
		wantCode   int
		wantResult string
	}{
		"ok request1: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL9-3N4M", "name": "Lettuce", "price": "3.41"}]`), wantCode: 200, wantResult: `{"status":"item exist, not added"}`},
		"ok request2: code": {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code": "A12T-4GH7-QPL9-3N4M", "name": "Lettuce", "price": "3.41"}]`), wantCode: 200, wantResult: `{"status":"item exist, not added"}`},
		"ok request: code":  {method: "POST", path: "/api/v1/add", jsonData: []byte(`[{"code":"ZRT6-72AS-K736-L4AZ","name":"Greener Pepper","price":"9.99"}]`), wantCode: 201, wantResult: `{"status":"item added"}`},
	}
	for name, tc := range tests {
		got := routerPOSTReq(tc.method, tc.path, tc.jsonData, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}

}
func TestGetItem(t *testing.T) {
	db := database{}
	router := db.dbInit()

	tests := map[string]struct {
		method     string
		path       string
		jsonData   []byte
		wantCode   int
		wantResult string
	}{
		"ok request1: code":  {method: "GET", path: "/api/v1/item/" + "A12T-4GH7-QPL9-3N4M", wantCode: 200, wantResult: `{"code":"A12T-4GH7-QPL9-3N4M","name":"Lettuce","price":"$3.41"}`},
		"ok request2: code":  {method: "GET", path: "/api/v1/item/" + "E5T6-9UI3-TH15-QR88", wantCode: 200, wantResult: `{"code":"E5T6-9UI3-TH15-QR88","name":"Peach","price":"$2.99"}`},
		"ok request3: code":  {method: "GET", path: "/api/v1/item/" + "YRT6-72AS-K736-L4AR", wantCode: 200, wantResult: `{"code":"YRT6-72AS-K736-L4AR","name":"Green Pepper","price":"$0.79"}`},
		"ok request4: code":  {method: "GET", path: "/api/v1/item/" + "TQ4C-VV6T-75ZX-1RMR", wantCode: 200, wantResult: `{"code":"TQ4C-VV6T-75ZX-1RMR","name":"Gala Apple","price":"$3.59"}`},
		"ok request5: code":  {method: "GET", path: "/api/v1/item/" + "TQ4C-VV6T-75ZX-1RMR1", wantCode: 400, wantResult: `{"error":"Key: 'ProduceId.ProduceCode' Error:Field validation for 'ProduceCode' failed on the 'isproducecode' tag"}`},
		"ok request6: code":  {method: "GET", path: "/api/v1/item/" + "Z5T6-9UI3-TH15-QR88", wantCode: 200, wantResult: `{"error":"code not found"}`},
		"bad request1: code": {method: "GET", path: "/api/v1/item/" + "Z5T6-9UI3-TH15-QR88", wantCode: 200, wantResult: `{"error":"code not found"}`},
	}
	for name, tc := range tests {
		got := routerGETReq(tc.method, tc.path, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}
}
func TestDelete(t *testing.T) {
	db := database{}
	router := db.dbInit()

	tests := map[string]struct {
		method     string
		path       string
		jsonData   []byte
		wantCode   int
		wantResult string
	}{
		"ok request1: code":  {method: "GET", path: "/api/v1/delete/" + "A12T-4GH7-QPL9-3N4M", wantCode: 200, wantResult: `{"status":"item deleted"}`},
		"ok request2: code":  {method: "GET", path: "/api/v1/delete/" + "TQ4C-VV6T-75ZX-1RMR", wantCode: 200, wantResult: `{"status":"item deleted"}`},
		"bad request1: code": {method: "GET", path: "/api/v1/delete/" + "TQ4C-VV6T-75ZX-1RMR1", wantCode: 400, wantResult: `{"error":"Key: 'ProduceId.ProduceCode' Error:Field validation for 'ProduceCode' failed on the 'isproducecode' tag"}`},
		// "bad request2: code": {method: "GET", path: "/api/v1/delete/" + "TQ4C-VV6T-75ZX-1RMR", wantCode: 200, wantResult: `{"error":"code not found"}`},
		"bad request3: code": {method: "GET", path: "/api/v1/delete/" + "", wantCode: 405, wantResult: `{"error":"endpoint not found"}`},
		"bad request4: code": {method: "GET", path: "/delet/" + "2", wantCode: 405, wantResult: `{"error":"endpoint not found"}`},
	}
	for name, tc := range tests {
		got := routerGETReq(tc.method, tc.path, router)
		if tc.wantCode != got.Code || tc.wantResult != got.Body.String() {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.wantCode, tc.wantResult)
		}
	}

}

var r []byte
var result []byte

func benchmRouterPing(b *testing.B) { // <2>
	b.ReportAllocs()
	db := database{}
	router := db.dbInit()
	for n := 0; n < b.N; n++ { // <5>
		// log.Println("ping")
		routerGETReq("GET", "/api/v1/ping", router)
		// r = out.Bytes() // <10>
	}

	// result = r // <11>
}

func Benchmark_RouterPing(b *testing.B) {
	benchmRouterPing(b)
}

func benchmGetItem(b *testing.B) { // <2>
	b.ReportAllocs()
	db := database{}
	router := db.dbInit()
	for n := 0; n < b.N; n++ { // <5>
		// log.Println("ping")
		routerGETReq("GET", "/api/v1/item/"+"A12T-4GH7-QPL9-3N4M", router)
		// r = out.Bytes() // <10>
	}

	// result = r // <11>
}
func Benchmark_GetItem(b *testing.B) {
	benchmRouterPing(b)
}
func benchmItems(b *testing.B) { // <2>
	b.ReportAllocs()
	db := database{}
	router := db.dbInit()
	for n := 0; n < b.N; n++ { // <5>
		// log.Println("ping")
		routerGETReq("GET", "/items", router)
		// r = out.Bytes() // <10>
	}

	// result = r // <11>
}
func Benchmark_Items(b *testing.B) {
	benchmItems(b)
}
