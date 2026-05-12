package handlers

import (
    "net/http"
    //"strconv"
    "forum-dark-jurassic/internal/services"
    "forum-dark-jurassic/internal/utils"
)

type CategoryHandler struct {
    Categories *services.CategoryService
}

func NewCategoryHandler(categories *services.CategoryService) *CategoryHandler {
    return &CategoryHandler{Categories: categories}
}

func (handler *CategoryHandler) ShowCategories(w http.ResponseWriter, r *http.Request) {
    categories, err := handler.Categories.GetAllCategories()
    if err != nil {
        http.Error(w, "Erreur interne.", http.StatusInternalServerError)
        return
    }
    utils.Render(w, "./internal/templates/categories.html", categories)
}

/*func (handler *CategoryHandler) ShowTopics(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    categoryID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Catégorie incorrecte.", http.StatusBadRequest)
        return
    }
    topics, err := handler.Categories.Topics.GetTopicsByCategory(categoryID)
    if err != nil {
        http.Error(w, "Erreur interne.", http.StatusInternalServerError)
        return
    }
    tmpl := template.Must(template.ParseFiles("./internal/templates/topics.html"))
    tmpl.Execute(w, topics)
}*/