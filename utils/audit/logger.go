package audit

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Action represents type of action being performed
type Action string

const (
	// ActionCreate represents creation of a resource
	ActionCreate Action = "create"
	// ActionUpdate represents update of a resource
	ActionUpdate Action = "update"
	// ActionDelete represents deletion of a resource
	ActionDelete Action = "delete"
	// ActionRead represents reading of a resource
	ActionRead Action = "read"
	// ActionLogin represents user login
	ActionLogin Action = "login"
)

// Resource represents the type of resource being audited
type Resource string

const (
	// ResourceUser represents user resource
	ResourceUser Resource = "user"
	// ResourceStation represents station resource
	ResourceStation Resource = "station"
	// ResourceRoute represents route resource
	ResourceRoute Resource = "route"
	// ResourceStaff represents staff resource
	ResourceStaff Resource = "staff"
	// ResourceVehicle represents vehicle resource
	ResourceVehicle Resource = "vehicle"
	// ResourceSchedule represents schedule resource
	ResourceSchedule Resource = "schedule"
	// ResourceFavorite represents favorite resource
	ResourceFavorite Resource = "favorite"
)

// LogEntry represents an audit log entry
type LogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	UserID      int       `json:"user_id"`
	UserEmail   string    `json:"user_email"`
	UserRole    string    `json:"user_role"`
	Action      Action    `json:"action"`
	Resource    Resource  `json:"resource"`
	ResourceID  string    `json:"resource_id,omitempty"`
	IP          string    `json:"ip"`
	Description string    `json:"description,omitempty"`
	Details     string    `json:"details,omitempty"`
}

// Log records an audit event
func Log(c *fiber.Ctx, action Action, resource Resource, resourceID string, description string, details interface{}) {
	// Get user info from context
	userID, _ := c.Locals("userID").(int)
	userEmail, _ := c.Locals("userEmail").(string)
	userRole, _ := c.Locals("userRole").(string)

	// Convert details to JSON string if provided
	var detailsStr string
	if details != nil {
		detailsBytes, err := json.Marshal(details)
		if err != nil {
			log.Printf("Error marshaling audit details: %v", err)
		} else {
			detailsStr = string(detailsBytes)
		}
	}

	// Create log entry
	entry := LogEntry{
		Timestamp:   time.Now(),
		UserID:      userID,
		UserEmail:   userEmail,
		UserRole:    userRole,
		Action:      action,
		Resource:    resource,
		ResourceID:  resourceID,
		IP:          c.IP(),
		Description: description,
		Details:     detailsStr,
	}

	// Log to console (in production this would go to a database or dedicated logging service)
	entryJSON, _ := json.Marshal(entry)
	log.Printf("[AUDIT] %s", string(entryJSON))

	// TODO: In production, store this in database or send to logging service
}

// LogWithContext records an audit event using context instead of fiber.Ctx
func LogWithContext(ctx context.Context, userID int, userEmail, userRole string, 
	action Action, resource Resource, resourceID string, 
	description string, details interface{}, ip string) {

	// Convert details to JSON string if provided
	var detailsStr string
	if details != nil {
		detailsBytes, err := json.Marshal(details)
		if err != nil {
			log.Printf("Error marshaling audit details: %v", err)
		} else {
			detailsStr = string(detailsBytes)
		}
	}

	// Create log entry
	entry := LogEntry{
		Timestamp:   time.Now(),
		UserID:      userID,
		UserEmail:   userEmail,
		UserRole:    userRole,
		Action:      action,
		Resource:    resource,
		ResourceID:  resourceID,
		IP:          ip,
		Description: description,
		Details:     detailsStr,
	}

	// Log to console (in production this would go to a database or dedicated logging service)
	entryJSON, _ := json.Marshal(entry)
	log.Printf("[AUDIT] %s", string(entryJSON))

	// TODO: In production, store this in database or send to logging service
}
