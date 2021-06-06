-- Options

-- Type Name
SelectTransactions
-- Type Comment

-- Func Name
SelectTransactions
-- Func Comment

-- SQL
SELECT
  COUNT(*) OVER() total_records
  , t.id
  , t.user_id
  , t.date
  , t.amount
  , t.note
  , t.created_at
  , t.modified_at
  , t.version
  , t.category_id
  , c.type category_type
  , c.name category_name
FROM
  kakeibo.transactions t
  INNER JOIN kakeibo.categories c 
    ON  t.category_id = c.id
WHERE
  (t.date >= %%from time.Time%% OR %%from time.Time%% IS NULL)
  AND (t.date <= %%to time.Time%% OR %%to time.Time%% IS NULL)
-- TODO 手動で実装
-- ORDER BY
--   %s %s, t.id
-- LIMIT ? OFFSET ?