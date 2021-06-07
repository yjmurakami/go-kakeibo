# go-kakeibo

家計簿のREST APIです。
家計簿のCRUDを行います。

<br>

## URL

http://go-kakeibo-lb-868206596.ap-northeast-1.elb.amazonaws.com

サンプル  
http://go-kakeibo-lb-868206596.ap-northeast-1.elb.amazonaws.com/api/v1/transactions?from=2021-06-01&to=2021-06-30&page=1&page_size=3&sort=-date

<br>

## JSONフォーマット

### 家計簿 transactions

```
{
  "id": 2,                    # ID
  "date": "2021-06-06",       # 日付
  "type": 2,                  # 費目分類 1=収入 2=支出
  "type_name": "支出",        # 費目分類名称
  "category_id": 5,           # 費目
  "category_name": "電気代",  # 費目名称
  "amount": 1234,             # 金額
  "note": "test"              # メモ
}
```

### 費目 categories

```
{
  "id": 1,              # ID
  "type": 1,            # 費目分類 1=収入 2=支出
  "type_name": "収入",  # 費目分類名称
  "name": "給料"        # 費目
},
```

<br>

## HTTPリクエスト サンプル

### 家計簿 参照

```
# リクエスト
curl 'http://go-kakeibo-lb-868206596.ap-northeast-1.elb.amazonaws.com/api/v1/transactions?from=2021-06-01&to=2021-06-30&page=1&page_size=3&sort=-date'

# レスポンス
{
  "metadata": {
    "current_page": 1,
    "page_size": 3,
    "first_page": 1,
    "last_page": 1,
    "total_records": 2
  },
  "data": [
    {
      "id": 2,
      "date": "2021-06-06",
      "type": 2,
      "type_name": "支出",
      "category_id": 5,
      "category_name": "電気代",
      "amount": 1234,
      "note": "test"
    },
    {
      "id": 8,
      "date": "2021-06-03",
      "type": 1,
      "type_name": "収入",
      "category_id": 1,
      "category_name": "給料",
      "amount": 9876,
      "note": "test changed"
    }
  ]
}
```

```
# リクエスト 絞り込みなし
curl 'http://go-kakeibo-lb-868206596.ap-northeast-1.elb.amazonaws.com/api/v1/transactions'
```

```
# リクエスト ID検索
# {id}を指定する
curl 'http://go-kakeibo-lb-868206596.ap-northeast-1.elb.amazonaws.com/api/v1/transactions/{id}' 
```

<br>

### 家計簿 登録
```
# リクエスト
curl -X POST -d '{"date": "2021-06-06", "category_id": 5, "amount": 1234, "note": "test"}' 'http://go-kakeibo-lb-868206596.ap-northeast-1.elb.amazonaws.com/api/v1/transactions' 
```

<br>

### 家計簿 更新

```
# リクエスト
# {id}を指定する
curl -X PATCH -d '{"date": "2021-06-03", "category_id": 1, "amount": 9876, "note": "test changed"}' 'http://go-kakeibo-lb-868206596.ap-northeast-1.elb.amazonaws.com/api/v1/transactions/{id}' 
```

<br>

### 家計簿 削除

```
# リクエスト
# {id}を指定する
curl -X DELETE 'http://go-kakeibo-lb-868206596.ap-northeast-1.elb.amazonaws.com/api/v1/transactions/{id}' 
```

<br>

### 費目 参照
```
# リクエスト
curl 'http://go-kakeibo-lb-868206596.ap-northeast-1.elb.amazonaws.com/api/v1/categories' 
```

<br>

## TODO

* ユーザー認証
* ユニットテスト
* インテグレーションテスト

