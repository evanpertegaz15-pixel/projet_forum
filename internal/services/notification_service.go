// like, commentaire, réponse, tag, follow, créer une notif, la marquer comme lue, les lister

package services

import (
	//"errors"
	"forum-dark-jurassic/internal/models"
)

type NotificationService struct {
	Notifications *models.NotificationModel
	Users         *models.UserModel
	Posts         *models.PostModel
}

func NewNotificationService(
	notifs *models.NotificationModel,
	users *models.UserModel,
	posts *models.PostModel,
) *NotificationService {
	return &NotificationService{
		Notifications: notifs,
		Users:         users,
		Posts:         posts,
	}
}

/*
// Types de notifications
const (
	TypeLike     = "like"
	TypeComment  = "comment"
	TypeReply    = "reply"
	TypeTag      = "tag"
	TypeFollow   = "follow"
)

// Créer une notification générique
func (s *NotificationService) createNotification(
	userID int,
	fromUserID int,
	notifType string,
	postID *int,
	commentID *int,
) error {

	if userID == fromUserID {
		return nil // éviter notif sur soi-même
	}

	return s.Notifications.Create(userID, fromUserID, notifType, postID, commentID)
}

// Notification like
func (s *NotificationService) NotifyLike(fromUserID, toUserID, postID int) error {
	return s.createNotification(toUserID, fromUserID, TypeLike, &postID, nil)
}

// Notification réponse
func (s *NotificationService) NotifyReply(fromUserID, toUserID, commentID int) error {
	return s.createNotification(toUserID, fromUserID, TypeReply, nil, &commentID)
}

// Notification tag
func (s *NotificationService) NotifyTag(fromUserID, toUserID, postID int) error {
	return s.createNotification(toUserID, fromUserID, TypeTag, &postID, nil)
}

// Notification follow
func (s *NotificationService) NotifyFollow(fromUserID, toUserID int) error {
	return s.createNotification(toUserID, fromUserID, TypeFollow, nil, nil)
}

// Lister les notifications d’un utilisateur
func (s *NotificationService) GetUserNotifications(userID int) ([]*models.NotificationModel, error) {
	return s.Notifications.GetByUser(userID)
}

// Marquer comme lue
func (s *NotificationService) MarkAsRead(notificationID int) error {
	return s.Notifications.MarkAsRead(notificationID)
}

// Tout marquer comme lu
func (s *NotificationService) MarkAllAsRead(userID int) error {
	return s.Notifications.MarkAllAsRead(userID)
}

// Compter non lues
func (s *NotificationService) CountUnread(userID int) (int, error) {
	return s.Notifications.CountUnread(userID)
}
*/