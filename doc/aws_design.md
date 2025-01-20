# kifupのAWSリソース情報

## 概要

### URLとリソースの紐づけ

|ドメイン|URL|CloudFront|制限|メソッド|ALB|ECS|S3|
|---|---|---|---|---|---|---|---|
|prd|`/api/*`|`cf-kifup`|なし|GET/HEAD/PUT/POST/DELETE|`alb-main`|`ecssvc-kifup-api`|-|
|prd|`/file/private/*`|`cf-kifup`|署名・OAC|GET/HEAD|-|-|`s3-kifup/private/*`|
|prd|`/file/private-upload/*`|`cf-kifup`|署名・OAC|GET/HEAD/PUT/POST/DELETE|-|-|`s3-kifup/private/*`|
|prd|`/file/public/*`|`cf-kifup`|OAC|GET/HEAD|-|-|`s3-kifup/public/*`|
|prd|`/file/public-upload/*`|`cf-kifup`|署名・OAC|GET/HEAD/PUT/POST/DELETE|-|-|`s3-kifup/public/*`|
|prd|その他|`cf-kifup`|なし|GET/HEAD|-|-|`s3-kifup/web/*`|
|stg|`/api/*`|`cf-kifup-stg`|なし|GET/HEAD/PUT/POST/DELETE|`alb-main`|`ecssvc-kifup-stg-api`|-|
|stg|`/swagger/*`|`cf-kifup-stg`|BASIC認証|GET/HEAD|`alb-main`|`ecssvc-kifup-stg-api`|-|
|stg|`/file/private/*`|`cf-kifup-stg`|BASIC認証・署名付きURL・OAC|GET/HEAD|-|-|`s3-kifup-stg/private/*`|
|stg|`/file/private-upload/*`|`cf-kifup-stg`|BASIC認証・署名付きURL・OAC|GET/HEAD/PUT/POST/DELETE|-|-|`s3-kifup-stg/private/*`|
|stg|`/file/public/*`|`cf-kifup-stg`|BASIC認証・OAC|GET/HEAD|-|-|`s3-kifup-stg/public/*`|
|stg|`/file/public-upload/*`|`cf-kifup-stg`|BASIC認証・署名付きURL・OAC|GET/HEAD/PUT/POST/DELETE|-|-|`s3-kifup-stg/public/*`|
|stg|その他|`cf-kifup-stg`|BASIC認証|GET/HEAD|-|-|`s3-kifup-stg/web/*`|

### 現在の状況

ステージング環境のみ構築済み。

## IAM

（省略）

## Route53

### <マイドメイン>

- 用途：事業で管理するWebコンテンツ全般
- レコード
  - kifup-stg.<マイドメイン>
    - タイプ：A
    - ルーティングポリシー：シンプル
    - エイリアス：はい
    - ルーティング先：cf-kifup-stgへのエイリアス

## CloudFront

### CF関数：cff-pre-release

- 用途：index.htmlの付与、ベーシック認証の設定
- 対象ビヘイビア：STG環境、API以外
- 処理内容1：swaggerのURIを除き、ビューワーリクエストのURIがスラッシュで終わる場合、index.htmlを付与する
- 処理内容2：ビューワーリクエストに関して、ベーシック認証を検証し、一致しなければ401エラーを返す

### CFディストリビューション：cf-kifup-stg

- 用途：kifup（ステージング環境）
- カスタムSSL証明書：設定済み
- デフォルトルートオブジェクト：index.html
- オリジン
  - origin-kifup-stg-api
    - オリジンドメイン：alb-mainのドメイン
    - プロトコル：HTTPのみ
    - オリジンパス：設定なし
    - カスタムヘッダー：`CloudFrontAuthentication`を指定
  - origin-kifup-stg-file-private
    - オリジンドメイン：s3-kifup-stgのドメイン
    - オリジンパス：`/private`
    - オリジンアクセス：OAC設定（`oac-kifup-stg`）
  - origin-kifup-stg-file-public
    - オリジンドメイン：s3-kifup-stgのドメイン
    - オリジンパス：`/public`
    - オリジンアクセス：OAC設定（`oac-kifup-stg`）
  - origin-kifup-stg-frontend
    - オリジンドメイン：s3-kifup-stgのドメイン
    - オリジンパス：`/frontend`
    - オリジンアクセス：OAC設定（`oac-kifup-stg`）
- ビヘイビア
  - 0(behavior-kifup-stg-api)
    - パスパターン：`/api/*`
    - オリジン：origin-kifup-stg-api
    - ビューワープロトコル：HTTPSのみ
    - キャッシュポリシー：UseOriginCacheControlHeaders
    - オリジンリクエストポリシー：Managed-AllViewer
  - 1(behavior-kifup-stg-swagger)
    - パスパターン：`/swagger/*`
    - オリジン：origin-kifup-stg-api
    - ビューワープロトコル：HTTPをHTTPSにリダイレクト
    - キャッシュポリシー：UseOriginCacheControlHeaders
    - オリジンリクエストポリシー：Managed-AllViewer
  - 2(behavior-kifup-stg-private-get)
    - パスパターン：`/file/private/*`
    - オリジン：origin-kifup-stg-file-private
    - ビューワープロトコル：HTTPSのみ
    - 許可されたHTTPメソッド：GET,HEAD
    - ビューワーのアクセス制限：keygroup-cf-kifup
    - キャッシュポリシー：Managed-CachingOptimized
    - オリジンリクエストポリシー：Managed-CORS-S3Origin
  - 3(behavior-kifup-stg-private-manage)
    - パスパターン：`/file/private-manage/*`
    - オリジン：origin-kifup-stg-file-private
    - ビューワープロトコル：HTTPSのみ
    - 許可されたHTTPメソッド：GET,HEAD,OPTIONS,PUT,POST,PATCH,DELETE
    - ビューワーのアクセス制限：keygroup-cf-kifup
    - キャッシュポリシー：Managed-CachingOptimized
    - オリジンリクエストポリシー：Managed-CORS-S3Origin
  - 4(behavior-kifup-stg-public-get)
    - パスパターン：`/file/public/*`
    - オリジン：origin-kifup-stg-file-public
    - ビューワープロトコル：HTTPSのみ
    - 許可されたHTTPメソッド：GET,HEAD
    - キャッシュポリシー：Managed-CachingOptimized
    - オリジンリクエストポリシー：Managed-CORS-S3Origin
  - 5(behavior-kifup-stg-public-manage)
    - パスパターン：`/file/public-manage/*`
    - オリジン：origin-kifup-stg-file-public
    - ビューワープロトコル：HTTPSのみ
    - 許可されたHTTPメソッド：GET,HEAD,OPTIONS,PUT,POST,PATCH,DELETE
    - ビューワーのアクセス制限：keygroup-cf-kifup
    - キャッシュポリシー：Managed-CachingOptimized
    - オリジンリクエストポリシー：Managed-CORS-S3Origin
  - 6(behavior-kifup-stg-frontend)
    - パスパターン：`デフォルト(*)`
    - オリジン：origin-kifup-stg-frontend
    - ビューワープロトコル：HTTPをHTTPSにリダイレクト
    - 許可されたHTTPメソッド：GET,HEAD
    - キャッシュポリシー：Managed-CachingOptimized
    - オリジンリクエストポリシー：Managed-CORS-S3Origin

## VPC

（省略）

## セキュリティグループ

（省略）

## ALB

### ALB：alb-main

- 用途：サーバー（EC2／ECS）を要する各種Webコンテンツ
- VPC：vpc-main
- アベイラビリティゾーン：`apne1-az1`,`apne1-az4`
- セキュリティグループ：secgrp-main-alb
- リスナーとルール
  - HTTP:80
    - rule-kifup-api-stg
      - 用途：kifupのAPIサーバー（ステージング環境）
      - 条件1：パスパターンが`/api/*`または`/swagger/*`
      - 条件2：HTTPヘッダー`CloudFrontAuthentication`の値を指定
      - 条件3：HTTPホストヘッダーが`kifup-stg.<マイドメイン>`
      - アクション：ターゲットグループ`target-kifup-api-stg`へ転送
    - デフォルト
      - 用途：想定外のリクエスト
      - 条件：他のルールが適用されない場合
      - アクション：固定レスポンスを返す（503 Service Unavailable）
- ターゲットグループ
  - target-kifup-api-stg
    - ターゲットの種類：IP
    - プロトコルとポート：HTTP:80

## ECR

### ECRリポジトリ：ecr-kifup-api

- 用途：kifupのAPIサーバー
- タグのイミュータビリティ：Immutable
- 暗号化タイプ：KMS

## ECS

### ECSサービス：ecssvc-kifup-api-stg

- 用途：kifupのAPIサーバー（ステージング環境）
- タスク定義：ecstask-kifup-api-stgの最新リビジョン
- キャパシティプロバイダー：FARGATE_SPOT
- ロードバランサー
  - ロードバランサー：alb-main
  - リスナープロトコルとポート：HTTP:80
  - ターゲットグループとプロトコル：target-kifup-api-stg:HTTP
  - ヘルスチェックパス：`/api/status`
  - ヘルスチェックの猶予期間：60
- ネットワーク
  - VPC：vpc-main
  - サブネット：subnet-service
  - セキュリティグループ：secgrp-main-webapp
  - パブリックIP：ON

### ECSタスク定義：ecstask-kifup-api-stg

- 起動タイプ：AWS Fargate
- OS、アーキテクチャ：Linux/X86_64
- タスクサイズ：0.5vCPU,1GB
- タスクロール：ecsTaskRole
- タスク実行ロール：ecsTaskExecutionRole
- コンテナ
  - kifup-api-stg
    - イメージURI：kifup-apiの最新タグ
    - 必須コンテナ：はい
    - ポートマッピング
      - 80,TCP,HTTP
    - リソース割り当て制限：指定なし（デフォルト）
    - 環境変数
      - ENV："staging"
      - SECRET_KEY："＜シークレットキー＞"
      - FRONTEND_ORIGIN："https://kifup-stg.<マイドメイン>"
    - ログ収集：CloudWatchで収集
    - ストレージ
      - エフェメラルストレージ：指定なし（デフォルトで20GiB）

## S3

### S3バケット：s3-kifup-stg

- 用途：kifup（ステージング環境）
- リージョン：`ap-northeast-1`
- アクセス許可
  - パブリックアクセスをすべてブロック
  - バケットポリシーにてCloudFrontからのアクセスを許可
    - リソース：`/frontend/*`,`/private/*`,`/public/*`
    - アクション：`s3:GetObject`,`s3:PutObject`
