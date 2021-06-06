-- Options

-- Type Name
SelectCategories
-- Type Comment

-- Func Name
SelectCategories
-- Func Comment

-- SQL
SELECT
  COUNT(*) OVER() total_records
  , id
  , type
  , name
  , created_at
  , modified_at
  , version
FROM
  kakeibo.categories
WHERE
  (type = %%categoryType int%% OR %%categoryType int%% = 0)
-- TODO 手動で実装
-- ORDER BY
--   %s %s, t.id
-- LIMIT ? OFFSET ?