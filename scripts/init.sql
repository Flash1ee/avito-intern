\c balance_db
drop table balance cascade;
drop table transactions;

create type transaction_type as
    enum ('write-off', 'refill', 'transfer');

CREATE TABLE IF NOT EXISTS balance
(
    user_id bigserial not null primary key,
    amount  numeric   not null default 0 check ( amount >= 0 )
);

CREATE TABLE IF NOT EXISTS transactions
(
    id          bigserial        not null primary key,
    type        transaction_type not null,
    sender_id   bigint references balance (user_id),
    receiver_id bigint references balance (user_id),
    amount      numeric          not null,
    created_at  timestamp        not null default now(),
    description text             not null
)