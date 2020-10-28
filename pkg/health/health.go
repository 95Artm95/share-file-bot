package health

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/friendsofgo/errors"
	"github.com/rs/zerolog/log"
)

func Check(ctx context.Context, addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return errors.Wrap(err, "split host and port")
	}

	if host == "" {
		host = "localhost"
	}

	u := fmt.Sprintf("http://%s:%s/health", host, port)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return errors.Wrap(err, "new request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("healthcheck failed")
	}

	return nil
}

// NewHandler creates a healthcheck handler
func NewHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := db.PingContext(ctx); err != nil {
			http.Error(w, "💩", http.StatusInternalServerError)
			log.Ctx(ctx).Error().Err(err).Msg("healthcheck")
		} else {
			_, _ = io.WriteString(w, "👌")
		}
	})
}
