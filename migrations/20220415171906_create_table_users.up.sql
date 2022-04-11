CREATE TABLE IF NOT EXISTS "users" ("balance" decimal(10,2) check(balance>=0) NOT NULL DEFAULT 0.000000,"id" serial NOT NULL UNIQUE,PRIMARY KEY ("id"));
