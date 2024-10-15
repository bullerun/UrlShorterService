package links

import (
	"UrlShorterService/internal/http_server/response"
	"context"
	"github.com/go-chi/render"
	_ "github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	id    string
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias" validate:"omitempty"`
}

type LinkRepositoryInterface interface {
	Save(ctx context.Context, url, alias string, id int64) error
}

func Save(log *slog.Logger, rep LinkRepositoryInterface) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "links.Save"
		var request Request
		if err := render.DecodeJSON(r.Body, &request); err != nil {
			log.Error(op, "failed to decode body", err)
			render.JSON(w, r, "failed to decode body")
			return
		}
		err := rep.Save(r.Context(), request.Url, request.Alias, r.Context().Value("userId").(int64))
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.JSON(w, r, response.OK())
	}
}
