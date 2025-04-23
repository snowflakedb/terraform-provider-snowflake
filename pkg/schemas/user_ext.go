package schemas

// This init is currently necessary to add sensitiveness to schemas which are generated (without such an underlying functionality yet).
func init() {
	// show output does not contain password or middle_name
	ShowUserSchema["login_name"].Sensitive = true
	ShowUserSchema["first_name"].Sensitive = true
	ShowUserSchema["last_name"].Sensitive = true
	ShowUserSchema["email"].Sensitive = true

	UserDescribeSchema["login_name"].Sensitive = true
	UserDescribeSchema["first_name"].Sensitive = true
	UserDescribeSchema["middle_name"].Sensitive = true
	UserDescribeSchema["last_name"].Sensitive = true
	UserDescribeSchema["email"].Sensitive = true
}
