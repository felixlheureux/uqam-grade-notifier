package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/alert"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/store"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	tokenStore   *store.TokenStore
	tokenManager *TokenManager
	courseStore  *store.CourseStore
	emailPass    string
	emailFrom    string
	baseURL      string
}

func NewHandler(tokenStore *store.TokenStore, tokenManager *TokenManager, courseStore *store.CourseStore, emailPass, emailFrom, baseURL string) *Handler {
	return &Handler{
		tokenStore:   tokenStore,
		tokenManager: tokenManager,
		courseStore:  courseStore,
		emailPass:    emailPass,
		emailFrom:    emailFrom,
		baseURL:      baseURL,
	}
}

type requestLogin struct {
	Email string `json:"email"`
}

type requestCourses struct {
	Semester string            `json:"semester"`
	Courses  map[string]string `json:"courses"`
}

type requestGradeUpdate struct {
	Grade string `json:"grade"`
}

func (h *Handler) RequestLogin(c echo.Context) error {
	var req requestLogin
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := ValidateEmail(req.Email); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	token, err := h.tokenManager.GenerateToken(req.Email, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	if err := h.tokenStore.SaveToken(req.Email, token); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save token")
	}

	loginURL := h.baseURL + "/auth/login?token=" + token
	subject := "Connexion à votre compte"
	body := "Cliquez sur le lien suivant pour vous connecter : " + loginURL

	if err := alert.SendEmail(h.emailPass, h.emailFrom, req.Email, subject, body); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send email")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Email sent successfully",
	})
}

func (h *Handler) Login(c echo.Context) error {
	tokenString := c.QueryParam("token")
	if tokenString == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Token is required")
	}

	claims, err := h.tokenManager.ValidateToken(tokenString)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
	}

	// Vérifier que le token est bien celui stocké pour cet email
	storedToken, err := h.tokenStore.GetToken(claims.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get token")
	}

	if storedToken != tokenString {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	if err := h.tokenStore.DeleteToken(claims.Email); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete token")
	}

	// Générer le token de session
	sessionToken, err := h.tokenManager.GenerateSessionToken(claims.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate session token")
	}

	// Créer le cookie de session
	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // 7 jours
		HttpOnly: true,
		Secure:   true, // Pour HTTPS
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login successful",
	})
}

func (h *Handler) SaveCourses(c echo.Context) error {
	// Vérifier la session
	sessionCookie, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Session invalide")
	}

	claims, err := h.tokenManager.ValidateSessionToken(sessionCookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Session expirée")
	}

	var req requestCourses
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Format de requête invalide")
	}

	// Valider le format des notes
	for course, grade := range req.Courses {
		if !isValidGrade(grade) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Format de note invalide pour le cours %s", course))
		}
	}

	changes, err := h.courseStore.SaveCourses(claims.Email, req.Semester, req.Courses)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Erreur lors de la sauvegarde des cours")
	}

	// Envoyer des notifications pour les changements de notes
	for _, change := range changes {
		subject := "Modification de note"
		body := fmt.Sprintf("La note du cours %s pour le semestre %s a été modifiée de %s à %s.",
			change.Course, change.Semester, change.OldGrade, change.NewGrade)

		if err := alert.SendEmail(h.emailPass, h.emailFrom, claims.Email, subject, body); err != nil {
			// Log l'erreur mais ne pas bloquer la réponse
			log.Printf("Erreur lors de l'envoi de la notification: %v", err)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cours sauvegardés avec succès"})
}

func (h *Handler) GetCourses(c echo.Context) error {
	// Vérifier la session
	sessionCookie, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Session invalide")
	}

	claims, err := h.tokenManager.ValidateSessionToken(sessionCookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Session expirée")
	}

	semester := c.QueryParam("semester")
	if semester == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Le paramètre semester est requis")
	}

	courses, err := h.courseStore.GetCourses(claims.Email, semester)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Erreur lors de la récupération des cours")
	}

	return c.JSON(http.StatusOK, courses)
}

func (h *Handler) UpdateGrade(c echo.Context) error {
	// Vérifier la session
	sessionCookie, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Session invalide")
	}

	claims, err := h.tokenManager.ValidateSessionToken(sessionCookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Session expirée")
	}

	semester := c.Param("semester")
	course := c.Param("course")

	var req requestGradeUpdate
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Format de requête invalide")
	}

	// Valider le format de la note
	if !isValidGrade(req.Grade) {
		return echo.NewHTTPError(http.StatusBadRequest, "Format de note invalide")
	}

	change, err := h.courseStore.UpdateGrade(claims.Email, semester, course, req.Grade)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if change != nil {
		// Envoyer une notification
		subject := "Modification de note"
		body := fmt.Sprintf("La note du cours %s pour le semestre %s a été modifiée de %s à %s.",
			change.Course, change.Semester, change.OldGrade, change.NewGrade)

		if err := alert.SendEmail(h.emailPass, h.emailFrom, claims.Email, subject, body); err != nil {
			// Log l'erreur mais ne pas bloquer la réponse
			log.Printf("Erreur lors de l'envoi de la notification: %v", err)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Note mise à jour avec succès"})
}

func (h *Handler) DeleteCourse(c echo.Context) error {
	// Vérifier la session
	sessionCookie, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Session invalide")
	}

	claims, err := h.tokenManager.ValidateSessionToken(sessionCookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Session expirée")
	}

	semester := c.Param("semester")
	course := c.Param("course")

	if err := h.courseStore.DeleteCourse(claims.Email, semester, course); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Erreur lors de la suppression du cours")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cours supprimé avec succès"})
}

// isValidGrade vérifie si la note est dans un format valide (nombre entre 0 et 100 avec 2 décimales)
func isValidGrade(grade string) bool {
	gradeFloat, err := strconv.ParseFloat(grade, 64)
	if err != nil {
		return false
	}
	return gradeFloat >= 0 && gradeFloat <= 100
}
