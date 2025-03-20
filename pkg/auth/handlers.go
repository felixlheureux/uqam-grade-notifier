package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/alert"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/httperror"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/service"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	authService   *service.AuthService
	courseService *service.CourseService
	emailPass     string
	emailFrom     string
	baseURL       string
}

func NewHandler(authService *service.AuthService, courseService *service.CourseService, emailPass, emailFrom, baseURL string) *Handler {
	return &Handler{
		authService:   authService,
		courseService: courseService,
		emailPass:     emailPass,
		emailFrom:     emailFrom,
		baseURL:       baseURL,
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
		return httperror.CoreRequestBindingFailed(err)
	}

	if err := ValidateEmail(req.Email); err != nil {
		return httperror.CoreRequestValidationFailed(err)
	}

	token, err := h.authService.GenerateToken(req.Email, 24*time.Hour)
	if err != nil {
		return httperror.CoreUnknownError(err)
	}

	if err := h.authService.SaveToken(req.Email, token); err != nil {
		return httperror.CoreUnknownError(err)
	}

	loginURL := h.baseURL + "/auth/login?token=" + token
	subject := "Login to your account"
	body := "Click on the following link to log in: " + loginURL

	if err := alert.SendEmail(h.emailPass, h.emailFrom, req.Email, subject, body); err != nil {
		return httperror.CoreUnknownError(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Email sent successfully",
	})
}

func (h *Handler) Login(c echo.Context) error {
	tokenString := c.QueryParam("token")
	if tokenString == "" {
		return httperror.CoreRequestValidationFailed(nil)
	}

	claims, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		return httperror.CoreUnauthorized(err)
	}

	storedToken, err := h.authService.GetToken(claims.Email)
	if err != nil {
		return httperror.CoreUnknownError(err)
	}

	if storedToken != tokenString {
		return httperror.CoreUnauthorized(nil)
	}

	if err := h.authService.DeleteToken(claims.Email); err != nil {
		return httperror.CoreUnknownError(err)
	}

	sessionToken, err := h.authService.GenerateSessionToken(claims.Email)
	if err != nil {
		return httperror.CoreUnknownError(err)
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login successful",
	})
}

// RequireAuth is a middleware that checks authentication
func (h *Handler) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session")
		if err != nil {
			return httperror.CoreUnauthorized(err)
		}

		email, err := h.authService.ValidateSessionToken(cookie.Value)
		if err != nil {
			return httperror.CoreUnauthorized(err)
		}

		c.Set("email", email)
		return next(c)
	}
}

// SaveCourses saves a user's courses for a given semester
func (h *Handler) SaveCourses(c echo.Context) error {
	email := c.Get("email").(string)
	semester := c.Param("semester")

	var req requestCourses
	if err := c.Bind(&req); err != nil {
		return httperror.CoreRequestBindingFailed(err)
	}

	ctx := domain.NewContext(c.Request(), c.Response())
	input := domain.CourseCreateInput{
		Email:    email,
		Semester: semester,
		Courses:  req.Courses,
	}

	course, err := h.courseService.Create(ctx, input)
	if err != nil {
		return httperror.CoreUnknownError(err)
	}

	return c.JSON(http.StatusOK, course)
}

// GetCourses retrieves a user's courses for a given semester
func (h *Handler) GetCourses(c echo.Context) error {
	ctx := domain.NewContext(c.Request(), c.Response())
	courses, err := h.courseService.Get(ctx)
	if err != nil {
		return httperror.CoreUnknownError(err)
	}

	return c.JSON(http.StatusOK, courses)
}

// DeleteCourse deletes a specific course
func (h *Handler) DeleteCourse(c echo.Context) error {
	course := c.Param("course")

	ctx := domain.NewContext(c.Request(), c.Response())
	if err := h.courseService.Delete(ctx, course); err != nil {
		return httperror.CoreUnknownError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateGrade updates a course grade
func (h *Handler) UpdateGrade(c echo.Context) error {
	course := c.Param("course")

	var req requestGradeUpdate
	if err := c.Bind(&req); err != nil {
		return httperror.CoreRequestBindingFailed(err)
	}

	if !isValidGrade(req.Grade) {
		return httperror.CoreRequestValidationFailed(nil)
	}

	ctx := domain.NewContext(c.Request(), c.Response())
	input := domain.CourseUpdateInput{
		Grade: req.Grade,
	}

	updatedCourse, err := h.courseService.Update(ctx, course, input)
	if err != nil {
		return httperror.CoreUnknownError(err)
	}

	return c.JSON(http.StatusOK, updatedCourse)
}

// DeleteCourses deletes all courses for a semester
func (h *Handler) DeleteCourses(c echo.Context) error {
	semester := c.Param("semester")

	ctx := domain.NewContext(c.Request(), c.Response())
	if err := h.courseService.Destroy(ctx, semester); err != nil {
		return httperror.CoreUnknownError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// isValidGrade checks if the grade is in a valid format (number between 0 and 100 with 2 decimals)
func isValidGrade(grade string) bool {
	gradeFloat, err := strconv.ParseFloat(grade, 64)
	if err != nil {
		return false
	}
	return gradeFloat >= 0 && gradeFloat <= 100
}
