package middleware

import (
	"log"
	"net/http"

	appErr "base-skeleton/internal/shared/errors"
	"base-skeleton/internal/shared/response"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rec := recover(); rec != nil {
				log.Println("panic:", rec)
				_ = response.JSON(
					w,
					http.StatusInternalServerError,
					"internal server error",
					nil,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// helper to write error response
func WriteError(w http.ResponseWriter, err error) {
	if e, ok := err.(*appErr.AppError); ok {
		_ = response.JSON(w, e.Code, e.Message, nil)
		return
	}

	_ = response.JSON(
		w,
		http.StatusInternalServerError,
		err.Error(),
		nil,
	)
}
