USE db_play;

DROP TABLE IF EXISTS db_play.forecast_jx;

CREATE TABLE IF NOT EXISTS forecast_jx (
                                       id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
                                       order_num VARCHAR(128) NOT NULL ,
                                       forecast_num INT NOT NULL ,
                                       forecast_result INT
)
