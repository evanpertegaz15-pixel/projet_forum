package handlers

import (
	"net/http"
	"forum-dark-jurassic/internal/models"
	"forum-dark-jurassic/internal/utils"
)

type HomeHandler struct {
	Home *models.HomeModel
}

func NewHomeHandler(home *models.HomeModel) *HomeHandler {
	return &HomeHandler{Home: home}
}

func (handler *HomeHandler) ShowHome(w http.ResponseWriter, r *http.Request) {
	data, err := handler.Home.GetHomePageData(20)
	if err != nil {
		utils.ErrorInternal(w, "Erreur interne.")
		return
	}

	for i := range data.Topics {
		data.Topics[i].CreatedAtAgo = utils.TimeAgo(data.Topics[i].CreatedAt)
	}

	utils.Render(w, "./internal/templates/home.html", map[string]any{
		"Users":       data.Users,
		"Topics":      data.Topics,
		"CategoryMap": data.CategoryMap,
	})
}
