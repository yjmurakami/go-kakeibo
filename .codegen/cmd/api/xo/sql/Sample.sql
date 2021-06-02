-- Options
-1
-- Type Name
Sample
-- Type Comment

-- Func Name
SelectSample
-- Func Comment

-- SQL
SELECT
  u.id
  , u.login_id
FROM
  kakeibo.users AS u 
  INNER JOIN kakeibo.transactions AS t
    ON  u.id = t.user_id
WHERE
  u.login_id = %%loginId int%%
  AND t.category_id = %%categoryId int%%
ORDER BY
  u.id
