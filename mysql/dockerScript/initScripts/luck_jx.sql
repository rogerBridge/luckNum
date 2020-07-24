USE db_play;

CREATE TABLE IF NOT EXISTS jx_luck (
    id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
    specific_num INT NOT NULL ,
    leave_value INT NOT NULL ,
    stop_probability FLOAT NOT NULL ,
    hope_income FLOAT NOT NULL
)