\c balance_db
CREATE TABLE IF NOT EXISTS balance
(
    user_id bigserial not null primary key,
    amount  numeric    not null default 0 check ( amount >= 0 )
);

CREATE TABLE IF NOT EXISTS transactions
(
    id bigserial not null primary key,
    from_id bigint references balance(user_id),
    to_id bigint references balance(user_id),
    amount numeric not null
)