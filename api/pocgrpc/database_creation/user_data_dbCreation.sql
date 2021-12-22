CREATE DATABASE IF NOT EXISTS 'gRPC_Db';
USE 'gRPC_Db';

CREATE Table IF NOT EXISTS 'user_data'(
    'id' INT(11) NOT NULL AUTO_INCREMENT,
    'pwd_hash' VARCHAR(256) NOT NULL,
    'name' VARCHAR(100) NOT NULL,
    'age' INT(2) NOT NULL,
    'parent_id' INT(11) NULL,

    PRIMARY KEY('id'),
    CONSTRAINT 'FK_user_parent_user' FOREIGN KEY('parent_id')
        REFERENCES 'gRPC_Db'.'user_data'('id')
)