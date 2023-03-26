package example_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dmitrymomot/go-api-server/svc/example"
	kitlog "github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Test example http handler without running server in docker
func TestExampleHandler(t *testing.T) {
	uid := uuid.New()
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", uid.String()), nil)
	rr := httptest.NewRecorder()
	handler := example.MakeHTTPHandler(
		example.MakeEndpoints(example.NewService()),
		kitlog.NewLogfmtLogger(os.Stderr),
	)
	handler.ServeHTTP(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	resp := &struct {
		Data string
	}{}
	err = json.Unmarshal(data, resp)
	assert.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("example: user id: %s", uid.String()), resp.Data)
}
