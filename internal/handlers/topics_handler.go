package handlers

import (
    "log"
    "net/http"
    "strconv"
    "forum-dark-jurassic/internal/models"
    "forum-dark-jurassic/internal/services"
    "forum-dark-jurassic/internal/utils"
)

type TopicHandler struct {
    Topics *services.TopicService
    Posts  *services.PostService
    Categories *services.CategoryService
    Auth *services.AuthService
}

func NewTopicHandler(topics *services.TopicService, posts *services.PostService, categories *services.CategoryService, auth *services.AuthService) *TopicHandler {
    return &TopicHandler{
        Topics: topics,
        Posts:  posts,
        Categories: categories,
        Auth:   auth,
    }
}

func (handler *TopicHandler) ShowTopics(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    categoryID, err := strconv.Atoi(idStr)
    if err != nil {
        utils.ErrorBadRequest(w, "Catégorie incorrecte.")
        return
    }
    topics, err := handler.Topics.GetTopicsByCategory(categoryID)
    if err != nil {
        utils.ErrorInternal(w, "Erreur interne.")
        return
    }
    for i := range topics {
        topics[i].CreatedAtAgo = utils.TimeAgo(topics[i].CreatedAt)
    }
    utils.Render(w,"./internal/templates/topics.html", topics)
}

func (handler *TopicHandler) ShowTopic(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    topicID, err := strconv.Atoi(idStr)
    if err != nil {
        utils.ErrorBadRequest(w, "Topic incorrect.")
        return
    }
    topic, err := handler.Topics.GetTopicByID(topicID)
    if err != nil {
        utils.ErrorNotFound(w, "Topic introuvable.")
        return
    }
    postsWithReplies, err := handler.Posts.GetPostsWithRepliesByTopic(topicID)
    if err != nil {
        utils.ErrorInternal(w, "Erreur interne.")
        return
    }
    topic.CreatedAtAgo = utils.TimeAgo(topic.CreatedAt)
    for i := range postsWithReplies {
        postsWithReplies[i].Post.CreatedAtAgo = utils.TimeAgo(postsWithReplies[i].Post.CreatedAt)
        for j := range postsWithReplies[i].Replies {
            postsWithReplies[i].Replies[j].CreatedAtAgo = utils.TimeAgo(postsWithReplies[i].Replies[j].CreatedAt)
        }
    }
    data := struct {
        Topic            models.Topic
        PostsWithReplies []models.PostWithReplies
    }{
        Topic:            topic,
        PostsWithReplies: postsWithReplies,
    }
    utils.Render(w,"./internal/templates/topic.html", data)
}

func (handler *TopicHandler) ShowNewTopic(w http.ResponseWriter, r *http.Request) {
    user, ok := RequireAuth(w, r, handler.Auth)
    if !ok { return }
    if user.HasRole("blocked") {
        utils.ErrorForbidden(w, "Votre compte ne peut plus créer de topics.")
        return
    }
    categories, err := handler.Categories.GetAllCategories()
    if err != nil {
        log.Printf("ShowNewTopic: GetAllCategories error: %v", err)
        utils.ErrorInternal(w, "Impossible de charger les catégories.")
        return
    }
    utils.Render(w, "./internal/templates/new_topic.html", map[string]any{
        "User":       user,
        "Categories": categories,
    })
}

func (handler *TopicHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        utils.ErrorMethodNotAllowed(w, "Méthode non autorisée")
        return
    }
    user, ok := RequireAuth(w, r, handler.Auth)
    if !ok { return }
    categoryID, _ := strconv.Atoi(r.FormValue("category_id"))
    title := r.FormValue("title")
    content := r.FormValue("content")
    if categoryID <= 0 || content == "" {
        utils.ErrorBadRequest(w, "Champs incorrects.")
        return
    }
    if user.HasRole("blocked") {
        utils.ErrorForbidden(w, "Votre compte ne peut plus créer de topics.")
        return
    }
    if err := utils.ValidatePostTitle(title); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    topicID, err := handler.Topics.CreateTopic(categoryID, user.ID, title)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    _, err = handler.Posts.CreatePost(topicID, user.ID, content)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, "/topic?id="+strconv.Itoa(topicID), http.StatusSeeOther)
}