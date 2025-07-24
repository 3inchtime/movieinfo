-- 插入初始数据

USE movieinfo;

-- 插入电影分类初始数据
INSERT INTO categories (name, description, sort_order) VALUES
                                                           ('动作', '动作冒险类电影', 1),
                                                           ('喜剧', '喜剧搞笑类电影', 2),
                                                           ('剧情', '剧情类电影', 3),
                                                           ('科幻', '科幻类电影', 4),
                                                           ('恐怖', '恐怖惊悚类电影', 5),
                                                           ('爱情', '爱情浪漫类电影', 6),
                                                           ('动画', '动画类电影', 7),
                                                           ('纪录片', '纪录片类电影', 8),
                                                           ('战争', '战争类电影', 9),
                                                           ('犯罪', '犯罪类电影', 10);

-- 插入示例电影数据
INSERT INTO movies (title, original_title, description, director, actors, release_date, duration, country, language, category_id, status) VALUES
                                                                                                                                              ('肖申克的救赎', 'The Shawshank Redemption', '讲述银行家安迪因被误判为杀害妻子及其情人的罪名入狱后，他与囚犯瑞德建立友谊，并在监狱中逐步获得影响力的故事。', '弗兰克·德拉邦特', '["蒂姆·罗宾斯", "摩根·弗里曼"]', '1994-09-23', 142, '美国', '英语', 3, 1),
                                                                                                                                              ('阿甘正传', 'Forrest Gump', '阿甘是一个智商只有75的低能儿，但他善良单纯，通过自己的努力创造了一个又一个奇迹。', '罗伯特·泽米吉斯', '["汤姆·汉克斯", "罗宾·怀特"]', '1994-07-06', 142, '美国', '英语', 3, 1),
                                                                                                                                              ('泰坦尼克号', 'Titanic', '1912年4月14日，载着1316号乘客和891名船员的豪华巨轮泰坦尼克号与冰山相撞而沉没，这场海难被认为是20世纪人间十大灾难之一。', '詹姆斯·卡梅隆', '["莱昂纳多·迪卡普里奥", "凯特·温斯莱特"]', '1997-12-19', 194, '美国', '英语', 6, 1);

-- 插入示例用户数据（密码为 'password123' 的哈希值，实际使用时应该用更安全的密码）
INSERT INTO users (username, email, password_hash, nickname, status, email_verified) VALUES
                                                                                         ('admin', 'admin@movieinfo.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1VdLSnqpjLjMTYcYxZ8VQjLOqpOqrAu', '管理员', 1, true),
                                                                                         ('testuser', 'test@movieinfo.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1VdLSnqpjLjMTYcYxZ8VQjLOqpOqrAu', '测试用户', 1, true);

-- 插入示例评分数据
INSERT INTO user_ratings (user_id, movie_id, rating, comment) VALUES
                                                                  (1, 1, 10, '经典中的经典，值得反复观看'),
                                                                  (1, 2, 9, '非常感人的电影'),
                                                                  (2, 1, 9, '很好的电影'),
                                                                  (2, 3, 8, '经典爱情电影');

-- 更新电影的评分统计（触发器或应用程序中处理）
UPDATE movies SET
                  rating_average = (SELECT AVG(rating) FROM user_ratings WHERE movie_id = movies.id),
                  rating_count = (SELECT COUNT(*) FROM user_ratings WHERE movie_id = movies.id)
WHERE id IN (1, 2, 3);
