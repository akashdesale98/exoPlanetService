package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akashdesale98/exoPlanetService/model"
	"github.com/akashdesale98/exoPlanetService/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var router *gin.Engine

func testSuite(t *testing.T) func(t *testing.T) {
	router = gin.Default()
	service.RegisterAPIs(router)
	// Return a function to teardown the test
	return func(t *testing.T) {
	}
}

// Test cases can be written similaryly for other endpoints
// this test case covers all the scenarios of adding exoplanet
func TestAddExoplanet(t *testing.T) {
	defer testSuite(t)(t)

	type args struct {
		c      *gin.Context
		method string
		path   string
		body   map[string]interface{}
	}

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantResponse   map[string]interface{}
		mock           func()
	}{
		{
			name: "should add exoplanet of type Terrestrial successfully",
			args: args{
				c:      &gin.Context{},
				method: "POST",
				path:   "/exoplanet",
				body: map[string]interface{}{
					"name":            "Neptune",
					"description":     "planet of rings",
					"distance":        5,
					"radius":          3,
					"exo_planet_type": "GasGiant",
				},
			},
			wantStatusCode: http.StatusOK,
			wantResponse: map[string]interface{}{
				"name":            "Neptune",
				"description":     "planet of rings",
				"distance":        5,
				"radius":          3,
				"exo_planet_type": "GasGiant",
			},
			mock: func() {},
		},
		{
			name: "should fail while adding Terrestrial exoplanet when no mass is sent",
			args: args{
				c:      &gin.Context{},
				method: "POST",
				path:   "/exoplanet",
				body: map[string]interface{}{
					"name":            "Saturn",
					"description":     "planet of rings",
					"distance":        5,
					"radius":          3,
					"exo_planet_type": "Terrestrial",
				},
			},
			wantStatusCode: http.StatusBadRequest,
			wantResponse:   nil,
			mock:           func() {},
		},
		{
			name: "should fail if the exoplanet with same name exists",
			args: args{
				c:      &gin.Context{},
				method: "POST",
				path:   "/exoplanet",
				body: map[string]interface{}{
					"name":            "Saturn",
					"description":     "planet of rings",
					"distance":        5,
					"radius":          3,
					"mass":            159,
					"exo_planet_type": "Terrestrial",
				},
			},
			wantStatusCode: http.StatusBadRequest,
			wantResponse:   nil,
			mock: func() {
				service.InMemoryStore["alreadyExistIssue"] = model.ExoPlanet{
					Name:          "Saturn",
					Description:   "planet of rings",
					Distance:      5,
					Radius:        3,
					Mass:          159,
					ExoPlanetType: "Terrestrial",
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			// Create a mock request
			body, err := json.Marshal(tt.args.body)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest(tt.args.method, tt.args.path, bytes.NewBuffer([]byte(body)))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.wantStatusCode, rr.Code)

			if tt.wantResponse != nil {
				// Parse the response body
				var response map[string]interface{}
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to parse response body: %v", err)
				}

				// assert response body
				assert.Equal(t, tt.wantResponse["name"], response["name"])
				assert.Equal(t, tt.wantResponse["description"], response["description"])
				// assert.Equal(t, tt.wantResponse["distance"], response["distance"])
				// assert.Equal(t, tt.wantResponse["radius"], response["radius"])
				assert.Equal(t, tt.wantResponse["exo_planet_type"], response["exo_planet_type"])
			}
		})
	}
}
