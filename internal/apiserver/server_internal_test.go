package apiserver

//
//import (
//	"bytes"
//	"encoding/json"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"restapi_langparser/internal/model"
//	"restapi_langparser/internal/store/teststore"
//	"testing"
//)
//
//func TestServer_HandleParse(t *testing.T) {
//	testCases := []struct {
//		name         string
//		payload      interface{}
//		expectedCode int
//	}{
//		{
//			name: "valid request",
//			payload: []model.Proxy{
//				{
//					"https://validurl.com/",
//				},
//			},
//			expectedCode: http.StatusOK,
//		},
//	}
//
//	store := teststore.New()
//	s := newServer(store)
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			rec := httptest.NewRecorder()
//			b := &bytes.Buffer{}
//			json.NewEncoder(b).Encode(tc.payload)
//			req, _ := http.NewRequest(http.MethodGet, "/parse", b)
//			s.router.ServeHTTP(rec, req)
//			assert.Equalf(t, tc.expectedCode, rec.Code, "Body of failed test case: %s", rec.Body)
//			_, err := CreateFromJSON(rec.Body.String())
//			assert.NoErrorf(t, err, "Bad response format: %s [%v]", rec.Body.String(), err)
//		})
//	}
//}
//
//func TestServer_HandleAddProxy(t *testing.T) {
//	testCases := []struct {
//		name         string
//		payload      interface{}
//		expectedCode int
//	}{
//		{
//			name: "valid request",
//			payload: []model.Proxy{
//				{
//					"https://validurl.com/",
//				},
//			},
//			expectedCode: http.StatusOK,
//		},
//		{
//			name: "validation fail",
//			payload: []model.Proxy{
//				{
//					"invalid url",
//				},
//			},
//			expectedCode: http.StatusInternalServerError,
//		},
//		{
//			name: "duplicate url",
//			payload: []model.Proxy{
//				{
//					"https://duplicateurl.com/",
//				},
//				{
//					"https://duplicateurl.com/",
//				},
//			},
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "bad request",
//			payload:      map[string]string{"bad": "request"},
//			expectedCode: http.StatusBadRequest,
//		},
//	}
//
//	store := teststore.New()
//	s := newServer(store)
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			rec := httptest.NewRecorder()
//			b := &bytes.Buffer{}
//			json.NewEncoder(b).Encode(tc.payload)
//			req, _ := http.NewRequest(http.MethodPost, "/proxy", b)
//			s.router.ServeHTTP(rec, req)
//			assert.Equalf(t, tc.expectedCode, rec.Code, "Body of failed test case: %s", rec.Body)
//			_, err := CreateFromJSON(rec.Body.String())
//			assert.NoErrorf(t, err, "Bad response format: %s [%v]", rec.Body.String(), err)
//		})
//	}
//}
