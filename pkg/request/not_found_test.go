package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"github.com/stretchr/testify/require"
)

func TestNotFoundHandler(t *testing.T) {
	// Setup logger
	l, err := logging.CommonLogger(logging.NewConfig(`tests`))
	require.NoError(t, err, "Failed to create logger")

	tests := []struct {
		name   string
		w      *httptest.ResponseRecorder
		r      *http.Request
		status int
		want   string
	}{
		{
			name:   "NotFound",
			w:      httptest.NewRecorder(),
			r:      httptest.NewRequest(http.MethodGet, "/", nil),
			status: http.StatusNotFound,
			want:   "{\"Message\":\"Not found\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NotFoundHandler(l).ServeHTTP(tt.w, tt.r)
			require.Equal(t, tt.status, tt.w.Code)
			require.Equal(t, tt.want, tt.w.Body.String())
		})
	}
}
