package rbac.resources.workspace

import data.rbac.roles

test_workspace_allowed {
    # User owns resource
    allow with input as
    {
        "actor": {
            "user": "steven",
            "op": "read",
            "roles": ["site-member"]
        },
        "object": {
            "type": "workspace",
            "owner": "steven",
            "id": "1234",
            "shared": []
        }
    } with data.rbac.roles as data.rbac.roles
}