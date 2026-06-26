package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/jwt"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
)

// Context keys for values stashed by the auth middleware.
const (
	CtxUserID      = "userID"
	CtxUsername    = "username"
	CtxRole        = "role"
	CtxSubjectType = "subjectType"
)

// Auth returns middleware that requires a valid access token. On success it
// stashes the user id, username, and role in the Fiber context.
func Auth(jwtMgr *jwt.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := bearerToken(c)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "unauthorized", "missing bearer token")
		}
		claims, err := jwtMgr.VerifyAccess(token)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid or expired token")
		}
		c.Locals(CtxUserID, claims.UserID)
		c.Locals(CtxUsername, claims.Username)
		c.Locals(CtxRole, claims.Role)
		c.Locals(CtxSubjectType, string(claims.SubjectType))
		return c.Next()
	}
}

// Optional returns middleware that attaches identity when a valid token is
// present but never rejects the request.
func Optional(jwtMgr *jwt.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if token, ok := bearerToken(c); ok {
			if claims, err := jwtMgr.VerifyAccess(token); err == nil {
				c.Locals(CtxUserID, claims.UserID)
				c.Locals(CtxUsername, claims.Username)
				c.Locals(CtxRole, claims.Role)
				c.Locals(CtxSubjectType, string(claims.SubjectType))
			}
		}
		return c.Next()
	}
}

// RequireRole returns middleware that rejects requests where the authenticated
// subject does not hold one of the allowed roles.
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals(CtxRole).(string)
		if _, ok := allowed[role]; !ok {
			return response.Error(c, fiber.StatusForbidden, "forbidden", "insufficient role")
		}
		return c.Next()
	}
}

// SubjectType returns the subject type from context ("user" or "company").
func SubjectType(c *fiber.Ctx) string {
	if v, ok := c.Locals(CtxSubjectType).(string); ok {
		return v
	}
	return "user"
}

func bearerToken(c *fiber.Ctx) (string, bool) {
	h := c.Get("Authorization")
	if h == "" {
		return "", false
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

// UserID returns the authenticated user id from the context, or "" if absent.
func UserID(c *fiber.Ctx) string {
	if v, ok := c.Locals(CtxUserID).(string); ok {
		return v
	}
	return ""
}
