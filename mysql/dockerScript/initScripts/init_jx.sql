USE db_play;

CREATE TABLE IF NOT EXISTS jx11x5 (
                                      id INT AUTO_INCREMENT PRIMARY KEY NOT NULL , /* 必须和主键一起使用*/
                                      order_number VARCHAR(128) , /* */
                                      order_time VARCHAR(128) ,
                                      one INT NOT NULL ,
                                      two INT NOT NULL ,
                                      three INT NOT NULL ,
                                      four INT NOT NULL ,
                                      five INT NOT NULL ,
                                      six INT NOT NULL ,
                                      seven INT NOT NULL ,
                                      eight INT NOT NULL ,
                                      nine INT NOT NULL ,
                                      ten INT NOT NULL ,
                                      eleven INT NOT NULL ,
                                      result TEXT ,
                                      time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
