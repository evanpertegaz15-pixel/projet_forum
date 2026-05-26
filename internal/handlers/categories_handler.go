package handlers

import (
    "net/http"
    "strconv"
    "forum-dark-jurassic/internal/models"
    "forum-dark-jurassic/internal/services"
    "forum-dark-jurassic/internal/utils"
)

type CategoryHandler struct {
    Categories *services.CategoryService
    Posts      *services.PostService
    Topics     *services.TopicService
    Auth       *services.AuthService
}

func NewCategoryHandler(categories *services.CategoryService, posts *services.PostService, topics *services.TopicService, auth *services.AuthService) *CategoryHandler {
    return &CategoryHandler{Categories: categories, Posts: posts, Topics: topics, Auth: auth}
}

func (handler *CategoryHandler) ShowCategories(w http.ResponseWriter, r *http.Request) {
    categories, err := handler.Categories.GetAllCategories()
    if err != nil {
        utils.ErrorInternal(w, "Erreur interne.")
        return
    }
    showForm := r.URL.Query().Get("show") == "1"
    selectedCategoryID := 0
    if idStr := r.URL.Query().Get("id"); idStr != "" {
        selectedCategoryID, _ = strconv.Atoi(idStr)
    }
    mine := r.URL.Query().Get("mine") == "1"
    liked := r.URL.Query().Get("liked") == "1"
    var user *models.User
    if sessionID := utils.GetCookie(r, "session_id"); sessionID != "" {
        user, _ = handler.Auth.GetUserFromSession(sessionID)
    }
    authorID := 0
    likedByUserID := 0
    filterMessage := ""
    skipPosts := false
    if mine {
        if user != nil {
            authorID = user.ID
        } else {
            filterMessage = "Connectez-vous pour voir vos posts."
            skipPosts = true
        }
    }
    if liked {
        if user != nil {
            likedByUserID = user.ID
        } else {
            if filterMessage != "" {
                filterMessage += " "
            }
            filterMessage += "Connectez-vous pour voir vos posts aimés."
            skipPosts = true
        }
    }
    showPosts := (r.URL.Query().Get("id") != "" || mine || liked)
    var posts []models.Post
    if showPosts && !skipPosts {
        posts, err = handler.Posts.GetFilteredPosts(selectedCategoryID, authorID, likedByUserID)
        if err != nil {
            utils.ErrorInternal(w, "Impossible de charger les posts filtrés.")
            return
        }
        for i := range posts {
            posts[i].CreatedAtAgo = utils.TimeAgo(posts[i].CreatedAt)
        }
    }
    var topics []models.Topic
    if showPosts && !skipPosts && liked && likedByUserID > 0 {
        topics, err = handler.Topics.GetLikedTopicsByUser(likedByUserID, selectedCategoryID)
        if err != nil {
            utils.ErrorInternal(w, "Impossible de charger les topics aimés.")
            return
        }
        for i := range topics {
            topics[i].CreatedAtAgo = utils.TimeAgo(topics[i].CreatedAt)
        }
    }
    data := struct {
        Categories         []models.Category
        Posts              []models.Post
        Topics             []models.Topic
        SelectedCategoryID int
        Mine               bool
        Liked              bool
        User               *models.User
        FilterMessage      string
        ShowForm           bool
        ShowPosts          bool
    }{
        Categories:         categories,
        Posts:              posts,
        Topics:             topics,
        SelectedCategoryID: selectedCategoryID,
        Mine:               mine,
        Liked:              liked,
        User:               user,
        FilterMessage:      filterMessage,
        ShowForm:           showForm,
        ShowPosts:          showPosts && !skipPosts,
    }
    utils.Render(w, "./internal/templates/categories.html", data)

}