package workspace

import data.roles
import data.users

# forCan is a function that returns true if the "for" is "all"
# or if "own" is set and the user is the same as the owner
forCan(for) {
	for == "all"
}

forCan(for) {
	input.user == input.owner
    for == "own"
}

# Default to deny
default allow = false
allow {
    # Get the user's roles
    userRoles := users.user_roles[input.user]

    # For each role in the list
    r := userRoles[_]

	# Perms has actions and objects that it can perform those actions on
    perms := roles.role_permissions[r][i]
   	acts := perms.actions


    # Does the operation that we are trying to do exist in the permission?
    acts.op[input.op]; perms.object[input.object]
    # Is the "for" correct?
    forCan(acts.for)
}

allow_1 {
    # Get the user's roles
    userRoles := users.user_roles[input.user]

    # For each role in the list
    r := userRoles[_]

    # Perms has actions and objects that it can perform those actions on
    perms := roles.role_permissions[r][i]
    acts := perms.actions

    # Does the operation that we are trying to do exist in the permission?
    acts.op[input.op]; perms.object[input.object]
}

x {
    true
    false
}
