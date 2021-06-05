-- Options
-1
-- Type Name
Transaction1
-- Type Comment

-- Func Name
SelectTransaction
-- Func Comment

-- SQL
SELECT
  t.id
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
  t.id = %%id int%%
