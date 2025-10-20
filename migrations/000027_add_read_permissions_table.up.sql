INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE code =  'users:read', 'role:read', 'facilitator_rating:read', 'course:read', 'course_posting:read', 'session:read', 'user_session:read', 'attendance:read')
FROM users;
