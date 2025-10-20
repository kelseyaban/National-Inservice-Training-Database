DELETE FROM permissions
WHERE code IN ('users:read', 'users:write', 'role:read', 'role:write', 'facilitator_rating:read', 'facilitator_rating:write', 'course:read', 'course:write', 'course_posting:read', 'course_posting:write', 'session:read', 'session:write', 'user_session:read', 'user_session:write', 'attendance:read', 'attendance:write');
