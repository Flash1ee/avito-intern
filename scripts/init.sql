CREATE DATABASE balance_db;
\c balance_db;
CREATE USER flashie WITH ENCRYPTED PASSWORD 'avito';
GRANT ALL PRIVILEGES ON DATABASE balance_db TO flashie;

CREATE TABLE IF NOT EXISTS balance
(
    user_id bigserial not null primary key,
    amount  bigint    not null default 0 check ( amount >= 0 )
);

CREATE TABLE IF NOT EXISTS transactions
(
    id bigserial not null primary key,
    from_id bigint references balance(user_id),
    to_id bigint references balance(user_id),
    amount bigint not null
)