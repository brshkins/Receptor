INSERT INTO ingredients (name) VALUES
('egg'),
('milk'),
('salt'),
('chicken'),
('tomato'),
('cheese');

INSERT INTO recipes (title, description, image_url, cooking_time, category)
VALUES
('Omelette', 'Simple omelette', 'https://img.com/omelette.jpg', 10, 'breakfast'),
('Chicken Salad', 'Healthy salad', 'https://img.com/salad.jpg', 20, 'lunch');

INSERT INTO recipe_ingredients (recipe_id, ingredient_id, amount) VALUES
(1, 1, '2 pcs'),
(1, 2, '50 ml'),
(1, 3, 'pinch'),

(2, 4, '200g'),
(2, 5, '1 pc'),
(2, 3, 'pinch');