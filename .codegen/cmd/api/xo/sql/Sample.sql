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
  , u.user_name
FROM
  kakeibo.users AS u 
  INNER JOIN kakeibo.incomes AS i
    ON  u.id = i.user_id
WHERE
  u.user_id = %%userId int%%
  AND i.category_id = %%categoryId int%%
ORDER BY
  u.id
