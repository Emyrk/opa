package rbac.resources.workspace

import data.rbac.roles


forCan(for) {
	for == "all"
}

forCan(for) {
	input.actor.user == input.object.owner
    for == "own"
}


forCan(for) {
	shared := {input.object.shared[_]}

	shared[input.actor.user]
}


default allow = false
allow {
	actor := input.actor
    obj := input.object

    # Get the user's roles
    roles := input.actor.roles

    # For each role in the list
    r := roles[_]

	# Perms has actions and objects that it can perform those actions on
    perms := data.rbac.roles.role_permissions[r][i]
   	acts := perms.actions


    acts.op[actor.op]; perms.object[obj.type]

    forCan(acts.for)
}
