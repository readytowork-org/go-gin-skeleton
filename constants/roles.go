package constants

const (
	// List of roles
	Role              = "role"
	RoleAdmin         = "admin"
	RoleClient        = "client"
	RoleClientAdmin   = "client_admin"
	RoleClientGeneral = "client_general"
	RoleUser          = "user"
	RoleClientUser    = "client_user"

	// List of ID for different roles
	AdminID         = "admin_id"
	ClientID        = "client_id"
	UserID          = "user_id"
	ClientAdminID   = "client_admin_id"
	ClientGeneralID = "client_general_id"
	ClientUserID    = "client_user_id"
)

var RolePrivileged = []string{RoleAdmin, RoleClientAdmin}
var RoleAll = []string{RoleAdmin, RoleClient, RoleClientAdmin, RoleClientGeneral, RoleUser}
