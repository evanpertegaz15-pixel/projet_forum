package handlers
// notifications -> voir et marquer comme lues les notifs
import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum-dark-jurassic/internal/services"
)

type NotificationsHandler struct {
	Notifications *services.NotificationService
}

func NewNotificationsHandler(notificationService *services.NotificationService) *NotificationsHandler {
	return &NotificationsHandler{
		Notifications: notificationService,
	}
}

// Voir toutes les notifications de l'utilisateur connecté
func (h *NotificationsHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	notifications, err := h.Notifications.GetUserNotifications(userID.(int))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// Marquer une notification comme lue
func (h *NotificationsHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	notificationID, err := strconv.Atoi(r.FormValue("notification_id"))
	if err != nil {
		http.Error(w, "notification_id invalide", http.StatusBadRequest)
		return
	}

	err = h.Notifications.MarkAsRead(notificationID, userID.(int))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Notification marquée comme lue",
	})
}

// Marquer toutes les notifications comme lues
func (h *NotificationsHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	err := h.Notifications.MarkAllAsRead(userID.(int))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Toutes les notifications ont été marquées comme lues",
	})
}

// Supprimer une notification (optionnel UX)
func (h *NotificationsHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	notificationID, err := strconv.Atoi(r.FormValue("notification_id"))
	if err != nil {
		http.Error(w, "notification_id invalide", http.StatusBadRequest)
		return
	}

	err = h.Notifications.DeleteNotification(notificationID, userID.(int))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Notification supprimée",
	})
}