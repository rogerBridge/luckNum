USE db_play;

CREATE TABLE IF NOT EXISTS forecast_jx (
                                       id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
                                       order_num VARCHAR(128) NOT NULL ,
                                       forecast_num INT NOT NULL ,
                                       forecast_result INT
)
