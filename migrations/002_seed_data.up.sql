-- Insert default categories
INSERT INTO categories (name, prefix, priority, color_code, description, icon) VALUES
    ('General', 'A', 1, '#3B82F6', 'General inquiries and services', 'users'),
    ('Priority', 'P', 5, '#EF4444', 'Priority services for VIP customers', 'star'),
    ('Billing', 'B', 2, '#10B981', 'Billing and payment services', 'credit-card'),
    ('Support', 'S', 3, '#F59E0B', 'Technical support', 'wrench'),
    ('Consultation', 'C', 2, '#8B5CF6', 'Consultation services', 'chat');

-- Insert default counters
INSERT INTO counters (number, name, location, status, is_active) VALUES
    ('1', 'Counter 1', 'Main Hall', 'active', true),
    ('2', 'Counter 2', 'Main Hall', 'active', true),
    ('3', 'Counter 3', 'Side Hall', 'inactive', true),
    ('4', 'Counter 4', 'Side Hall', 'inactive', true);

-- Assign categories to counters
INSERT INTO counter_categories (counter_id, category_id) VALUES
    (1, 1), (1, 2), -- Counter 1: General, Priority
    (2, 1), (2, 3), -- Counter 2: General, Billing
    (3, 4),         -- Counter 3: Support
    (4, 3), (4, 5); -- Counter 4: Billing, Consultation

-- Insert default admin user (password: admin123)
INSERT INTO users (username, password, full_name, email, role, is_active) VALUES
    ('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'System Administrator', 'admin@queuesystem.com', 'admin', true);

-- Insert sample staff users (password: staff123)
INSERT INTO users (username, password, full_name, email, role, counter_id, is_active) VALUES
    ('staff1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'John Smith', 'john@queuesystem.com', 'staff', 1, true),
    ('staff2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Jane Doe', 'jane@queuesystem.com', 'staff', 2, true);

-- Update counters with current staff
UPDATE counters SET current_staff_id = 2 WHERE id = 1;
UPDATE counters SET current_staff_id = 3 WHERE id = 2;
