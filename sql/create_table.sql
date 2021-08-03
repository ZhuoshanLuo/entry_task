//用户表
CREATE TABLE IF NOT EXISTS user_tab(
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    passwd VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    Avatar VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL,
    created_at INT UNSIGNED NOT NULL,
    PRIMARY KEY(id),
    UNIQUE KEY uniq_name(name),
    UNIQUE KEY uniq_email(email)
);

//活动类型表
CREATE TABLE IF NOT EXISTS activities_type_tab(
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    PRIMARY KEY(id),
    UNIQUE KEY uniq_name(name)
);

//活动表
CREATE TABLE IF NOT EXISTS activities_tab(
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    type_id INT UNSIGNED NOT NULL,
    title VARCHAR(50) NOT NULL,
    content VARCHAR(255) NOT NULL,
    location VARCHAR(50) NOT NULL,
    start_time INT UNSIGNED NOT NULL,
    end_time INT UNSIGNED NOT NULL,
    PRIMARY KEY(id)
);

//报名表
CREATE TABLE IF NOT EXISTS form_tab(
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    activity_id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    joined_at INT UNSIGNED NOT NULL,
    PRIMARY KEY(id)
);

//评论表
CREATE TABLE IF NOT EXISTS comments_tab(
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id INT UNSIGNED NOT NULL,
    activity_id INT UNSIGNED NOT NULL,
    content VARCHAR(50) NOT NULL,
    created_at INT UNSIGNED NOT NULL,
    PRIMARY KEY(id)
);

//session表
CREATE TABLE IF NOT EXISTS session_tab(
    id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    PRIMARY KEY(id)
);
