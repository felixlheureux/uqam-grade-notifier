package auth

import (
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

// RequireAuth est un middleware qui vérifie l'authentification
func (h *Handler) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Non authentifié")
		}

		email, err := h.tokenStore.ValidateSessionToken(cookie.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Session invalide")
		}

		c.Set("email", email)
		return next(c)
	}
}

// SaveCourses sauvegarde les cours d'un utilisateur pour un semestre donné
func (h *Handler) SaveCourses(c echo.Context) error {
	email := c.Get("email").(string)
	semester := c.Param("semester")

	var req requestCourses
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Format de requête invalide")
	}

	changes, err := h.courseStore.SaveCourses(email, semester, req.Courses)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, changes)
}

// GetCourses récupère les cours d'un utilisateur pour un semestre donné
func (h *Handler) GetCourses(c echo.Context) error {
	email := c.Get("email").(string)
	semester := c.Param("semester")

	courses, err := h.courseStore.GetCourses(email, semester)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, courses)
}

// DeleteCourse supprime un cours spécifique
func (h *Handler) DeleteCourse(c echo.Context) error {
	email := c.Get("email").(string)
	semester := c.Param("semester")
	course := c.Param("course")

	if err := h.courseStore.DeleteCourse(email, semester, course); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateGrade met à jour la note d'un cours
func (h *Handler) UpdateGrade(c echo.Context) error {
	email := c.Get("email").(string)
	semester := c.Param("semester")
	course := c.Param("course")

	var req struct {
		Grade string `json:"grade"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Format de requête invalide")
	}

	change, err := h.courseStore.UpdateGrade(email, semester, course, req.Grade)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if change == nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, change)
}

// DeleteCourses supprime tous les cours d'un semestre
func (h *Handler) DeleteCourses(c echo.Context) error {
	email := c.Get("email").(string)
	semester := c.Param("semester")

	if err := h.courseStore.DeleteCourses(email, semester); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// isValidGrade vérifie si la note est dans un format valide (nombre entre 0 et 100 avec 2 décimales)
func isValidGrade(grade string) bool {
	gradeFloat, err := strconv.ParseFloat(grade, 64)
	if err != nil {
		return false
	}
	return gradeFloat >= 0 && gradeFloat <= 100
}
