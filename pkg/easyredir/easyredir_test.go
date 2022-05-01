package easyredir

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRules(t *testing.T) {
	tests := []struct {
		name       string
		file       string
		wantStatus int
		wantRules  int
		wantIDs    []string
		wantError  error
	}{
		{
			name:      "should return no rules",
			file:      "testdata/get_rules_no_results.json",
			wantRules: 0,
			wantIDs:   []string{},
			wantError: nil,
		}, {
			name:      "should return one rule",
			file:      "testdata/get_rules_one_result.json",
			wantRules: 1,
			wantIDs:   []string{"abc-def"},
			wantError: nil,
		}, {
			name:      "should return two rules",
			file:      "testdata/get_rules_two_results.json",
			wantRules: 2,
			wantIDs:   []string{"abc-def", "def-abc"},
			wantError: nil,
		}, {
			name:      "should return a get rules error",
			file:      "testdata/get_rules_error_invalid_results.json",
			wantRules: 0,
			wantIDs:   []string{},
			wantError: errors.New("unable to get rules: unable to parse json: invalid character 'h' in literal true (expecting 'r')"),
		},
	}

	assert := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := ioutil.ReadFile(tt.file)
			if err != nil {
				log.Fatal(err)
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(body)
			}))
			defer server.Close()

			e := New(&Options{
				APIKey:    "",
				APISecret: "",
			})
			e.Client.baseURL = server.URL
			err = e.GetRules()

			IDs := []string{}
			for _, rd := range e.Rules.Data {
				IDs = append(IDs, rd.ID)
			}

			assert.Equal(tt.wantRules, len(e.Rules.Data))
			assert.Equal(tt.wantIDs, IDs)
			if tt.wantError != nil {
				assert.EqualError(err, tt.wantError.Error())
			}
		})
	}
}
