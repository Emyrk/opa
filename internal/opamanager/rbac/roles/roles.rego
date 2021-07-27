package rbac.roles

role_permissions := {
    "system" : [
        {"actions": {"op":{"read", "write"}, "for":"all"}, "object": {"workspace"}},
        {"actions": {"op":{"read", "write"}, "for":"all"}, "object": {"images"}},
    ],
    "site-admin": [
        {"actions": {"op":{"read", "write"}, "for":"all"}, "object": {"workspace"}},
    ],
    "site-member": [
        {"actions": {"op":{"read", "write"}, "for":"own"}, "object": {"workspace"}},
    ],
}
