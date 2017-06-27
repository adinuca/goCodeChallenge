-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create TABLE expense (
  id   INT AUTO_INCREMENT NOT NULL,
  date VARCHAR(10)        NOT NULL,
  amount FLOAT(10,2) NOT NULL ,
  reason VARCHAR(250) NOT NULL ,
  PRIMARY KEY (id)
);
