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
    Likes *services.LikeService
    ImageService *services.ImageService
    ImageModel *models.ImageModel
}

func NewTopicHandler(topics *services.TopicService, posts *services.PostService, categories *services.CategoryService, auth *services.AuthService, likes *services.LikeService, imageService *services.ImageService, imageModel *models.ImageModel) *TopicHandler {
    return &TopicHandler{
        Topics: topics,
        Posts:  posts,
        Categories: categories,
        Auth:   auth,
        Likes:  likes,
        ImageService: imageService,
        ImageModel: imageModel,
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
    topicLikes, err := handler.Likes.CountTopicLikes(topicID)
    if err != nil {
        utils.ErrorInternal(w, "Erreur interne.")
        return
    }
    topic.LikesCount = topicLikes
    for i := range postsWithReplies {
        postsWithReplies[i].Post.CreatedAtAgo = utils.TimeAgo(postsWithReplies[i].Post.CreatedAt)
        count, err := handler.Likes.CountPostLikes(postsWithReplies[i].Post.ID)
        if err != nil {
            utils.ErrorInternal(w, "Erreur interne.")
            return
        }
        postsWithReplies[i].Post.LikesCount = count
        postImages, err := handler.ImageModel.GetImagesByPostID(postsWithReplies[i].Post.ID)
        if err != nil {
            utils.ErrorInternal(w, "Erreur interne.")
            return
        }
        postsWithReplies[i].Post.Images = postImages
        dcount, err := handler.Likes.CountPostDislikes(postsWithReplies[i].Post.ID)
        if err != nil {
            utils.ErrorInternal(w, "Erreur interne.")
            return
        }
        postsWithReplies[i].Post.DislikesCount = dcount
        for j := range postsWithReplies[i].Replies {
            postsWithReplies[i].Replies[j].CreatedAtAgo = utils.TimeAgo(postsWithReplies[i].Replies[j].CreatedAt)
            replyCount, err := handler.Likes.CountPostLikes(postsWithReplies[i].Replies[j].ID)
            if err != nil {
                utils.ErrorInternal(w, "Erreur interne.")
                return
            }
            postsWithReplies[i].Replies[j].LikesCount = replyCount
            replyImages, err := handler.ImageModel.GetImagesByPostID(postsWithReplies[i].Replies[j].ID)
            if err != nil {
                utils.ErrorInternal(w, "Erreur interne.")
                return
            }
            postsWithReplies[i].Replies[j].Images = replyImages
            replyDcount, err := handler.Likes.CountPostDislikes(postsWithReplies[i].Replies[j].ID)
            if err != nil {
                utils.ErrorInternal(w, "Erreur interne.")
                return
            }
            postsWithReplies[i].Replies[j].DislikesCount = replyDcount
        }
    }
    user, _ := RequireAuth(w, r, handler.Auth)
    reportStatus := ""
    switch r.URL.Query().Get("report_status") {
    case "success":
        reportStatus = "Signalement envoyé."
    case "error":
        reportStatus = "Erreur lors du signalement."
    }
    data := struct {
        Topic            models.Topic
        PostsWithReplies []models.PostWithReplies
        User             *models.User
        ReportStatus     string
        ReturnURL        string
    }{
        Topic:            topic,
        PostsWithReplies: postsWithReplies,
        User:             user,
        ReportStatus:     reportStatus,
        ReturnURL:        "/topic?id=" + strconv.Itoa(topicID),
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
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        utils.ErrorBadRequest(w, "Impossible de traiter le formulaire.")
        return
    }
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
    postID, err := handler.Posts.CreatePost(topicID, user.ID, content)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    file, header, err := r.FormFile("image")
    if err == nil {
        filename, uploadErr := handler.ImageService.UploadImage(file, header)
        if uploadErr != nil {
            log.Printf("CreateTopic: upload error for %q: %v", header.Filename, uploadErr)
            utils.ErrorBadRequest(w, "Impossible d’uploader l’image.")
            return
        }
        _, dbErr := handler.ImageModel.CreateImage(filename, user.ID, &postID)
        if dbErr != nil {
            log.Printf("CreateTopic: image db error for %q: %v", filename, dbErr)
            utils.ErrorInternal(w, "Impossible d’enregistrer l’image.")
            return
        }
    } else if err != http.ErrMissingFile {
        utils.ErrorBadRequest(w, "Impossible de traiter le fichier image.")
        return
    }
    http.Redirect(w, r, "/topic?id="+strconv.Itoa(topicID), http.StatusSeeOther)
}

func (handler *TopicHandler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        utils.ErrorMethodNotAllowed(w, "Méthode non autorisée.")
        return
    }
    user, ok := RequireAuth(w, r, handler.Auth)
    if !ok { return }
    topicID, err := strconv.Atoi(r.FormValue("topic_id"))
    if err != nil {
        utils.ErrorBadRequest(w, "ID de topic invalide.")
        return
    }
    err = handler.Topics.DeleteTopic(user, topicID)
    if err != nil {
        utils.ErrorForbidden(w, "Permission refusée.")
        return
    }
    categoryID := r.FormValue("category_id")
    if categoryID != "" {
        http.Redirect(w, r, "/topics?id="+categoryID, http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}