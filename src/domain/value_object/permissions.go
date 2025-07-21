package value_object

import "encoding/json"

type Permissions []string

func ParsePermissions(data json.RawMessage) (Permissions, error) {
	var permissions Permissions

	err := json.Unmarshal(data, &permissions)
	return permissions, err
}

func AllPermissions() Permissions {
	var permissions = Permissions{
		// User
		"view_user",
		"create_user",
		"update_user",
		"delete_user",

		// Role
		"view_role",
		"create_role",
		"update_role",
		"delete_role",

		// Client
		"view_client",
		"create_client",
		"update_client",
		"delete_client",

		// Driver
		"view_driver",
		"create_driver",
		"update_driver",
		"delete_driver",

		// Location
		"view_location",
		"create_location",
		"update_location",
		"delete_location",

		// Location
		"view_vehicle",
		"create_vehicle",
		"update_vehicle",
		"delete_vehicle",

		// Settings
		"view_settings",
		"create_settings",
		"update_settings",
		"delete_settings",

		// Tariff
		"view_tariff",
		"create_tariff",
		"update_tariff",
		"delete_tariff",

		// Inquiry
		"view_inquiry",
		"create_inquiry",
		"update_inquiry",
		"delete_inquiry",
	}

	return permissions
}
