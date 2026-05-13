package handlers

import (
    "net/http"
    "strconv"
    "forum-dark-jurassic/internal/models"
    "forum-dark-jurassic/internal/services"
    "forum-dark-jurassic/internal/utils"
)

type PostHandler struct {
    Posts *services.PostService
    Auth *services.AuthService
}

func NewPostHandler(posts *services.PostService, auth *services.AuthService) *PostHandler {
    return &PostHandler{Posts: posts, Auth: auth,}
}

func (handler *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        utils.ErrorMethodNotAllowed(w, "Méthode non autorisée.")
        return
    }
    user, ok := RequireAuth(w, r, handler.Auth)
    if !ok { return }
    if user.HasRole("blocked") {
        utils.ErrorForbidden(w, "Votre compte ne peut plus publier de messages.")
        return
    }
    topicID, err := strconv.Atoi(r.FormValue("topic_id"))
    if err != nil || topicID <= 0 {
        utils.ErrorBadRequest(w, "Topic invalide.")
        return
    }
    content := r.FormValue("content")
    if content == "" {
        utils.ErrorBadRequest(w, "Le contenu ne peut pas être vide.")
        return
    }
    _, err = handler.Posts.CreatePost(topicID, user.ID, content)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, "/topic?id="+strconv.Itoa(topicID), http.StatusSeeOther)
}

func (handler *PostHandler) ShowCreatePostForm(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        utils.ErrorMethodNotAllowed(w, "Méthode non autorisée.")
        return
    }
    user, ok := RequireAuth(w, r, handler.Auth)
    if !ok { return }
    if user.HasRole("blocked") {
        utils.ErrorForbidden(w, "Votre compte ne peut plus publier de messages.")
        return
    }
    topicID, _ := strconv.Atoi(r.URL.Query().Get("topic_id"))
    data := struct {
        User    *models.User
        TopicID int
    }{
        User:   user,
        TopicID:    topicID,
    }
    utils.Render(w, "./internal/templates/post_create.html", data)
}

func (handler *PostHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        utils.ErrorMethodNotAllowed(w, "Méthode non autorisée.")
        return
    }
    postID, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || postID <= 0 {
        utils.ErrorBadRequest(w, "Post invalide.")
        return
    }
    post, err := handler.Posts.Posts.GetPostByID(postID)
    if err != nil {
        utils.ErrorNotFound(w, "Post introuvable.")
        return
    }
    replies, err := handler.Posts.GetReplies(postID)
    if err != nil {
        utils.ErrorInternal(w, "Erreur lors du chargement des réponses.")
        return
    }
    post.CreatedAtAgo = utils.TimeAgo(post.CreatedAt)
    for i := range replies {
        replies[i].CreatedAtAgo = utils.TimeAgo(replies[i].CreatedAt)
    }
    data := struct {
        Post    models.Post
        Replies []models.Post
    }{
        Post:    post,
        Replies: replies,
    }
    utils.Render(w, "./internal/templates/post.html", data)
}

func (handler *PostHandler) CreateReply(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        utils.ErrorMethodNotAllowed(w, "Méthode non autorisée.")
        return
    }
    user, ok := RequireAuth(w, r, handler.Auth)
    if !ok { return }
    if user.HasRole("blocked") {
        utils.ErrorForbidden(w, "Votre compte ne peut plus publier de réponses.")
        return
    }
    topicID, _ := strconv.Atoi(r.FormValue("topic_id"))
    parentID, _ := strconv.Atoi(r.FormValue("parent_id"))
    content := r.FormValue("content")
    parentPost, err := handler.Posts.Posts.GetPostByID(parentID)
    if err != nil {
        utils.ErrorInternal(w, "Erreur lors de la vérification du post parent.")
        return
    }
    if parentPost.ParentID != nil {
        utils.ErrorBadRequest(w, "Impossible de répondre à une réponse.")
        return
    }
    _, err = handler.Posts.CreateReply(topicID, user.ID, parentID, content)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    http.Redirect(w, r, "/topic?id="+strconv.Itoa(topicID), http.StatusSeeOther)
}