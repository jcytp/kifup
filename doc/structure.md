# kifupのシステム構成

## 全体（kifup）

### アプリケーション概要

- 将棋の棋譜を管理するWebアプリケーション
- ユーザーは棋譜のアップロード・新規作成・コメント付与・公開などができる
- アカウント未登録のユーザーでも棋譜を閲覧可能

### 技術スタック

- DB
  - SQLite
  - 実行方法：バックエンドサーバーで作成/更新
  - バックアップ：AWS S3に保存
- バックエンド
  - 言語：Go
  - フレームワーク：Gin
  - デプロイ方法：AWS ECSにFargateでコンテナを起動
- フロントエンド
  - 言語：TypeScript
  - フレームワーク：SvelteKit
  - デプロイ方法：SSGで静的エクスポートし、S3上に配置
- API定義
  - OpenAPI 3.0 (Swagger)
  - 閲覧方法：開発環境のバックエンドサーバーでSwaggerUIを提供

### インフラ構成

- AWS S3
  - フロントエンド（静的エクスポート）
  - SQLiteデータベースファイル
  - 棋譜ファイル
- AWS ECS Fargate
  - APIサーバー

### ソースコードの構成

```text
kifup/
├── api/      # バックエンド（kifup-api）
├── dev/      # Docckerfileや開発サーバー起動用シェル
├── doc/      # 開発資料となるドキュメント類
└── frontend/ # フロントエンド（kifup-frontend）
```

## バックエンド（kifup-api）

### apiディレクトリの役割

```text
kifup/api/
├── common/         # 各APIエンドポイントで共通の処理
│   ├── auxi/       # 汎用処理
│   ├── aws/        # AWS関連の処理
│   ├── db/         # DB接続
│   ├── env/        # 環境変数の処理
│   ├── handler/    # ginのミドルウェア、ハンドラー、レスポンス生成用関数など
│   └── log/        # ロギング設定
├── service/        # ビジネスロジックとデータ構造
│   ├── api/        # 各APIエンドポイントに対応する関数群
│   │   └── parser/ # 棋譜データのパーサー
│   ├── dao/        # データベースへのアクセス
│   └── model/      # 構造体の定義
├── swagger/        # SwaggerUIの表示用コンテンツ
├── go.mod
├── go.sum
└── main.go         # エントリーポイント
```

### apiソースコードツリー

```text
kifup/api$ tree -F -I "tmp|build|.*" --dirsfirst --noreport

```

### api説明

- APIエンドポイントに対応する関数（API関数）は、service/apiパッケージに記述される。
- handler関数でラッピングすることにより、API関数はリクエストデータを引数にとり、レスポンスのデータ部分とエラーを戻り値とする。
- DBに対する基本的な操作はservice/daoパッケージに記述し、主要な構造体はservice/modelパッケージに記述する。

## フロントエンド（kifup-frontend）

### frontendソースコードツリー

```text
kifup/frontend$ tree -F -I "node_modules|build|.*" --dirsfirst --noreport

```

### frontend説明

- URLに対応するページはsrc/routes以下の+page.svelteで記述される。
- 独立性の高い部分や共通コンポーネントはsrc/lib/componentsに配置する。
- ssr=false;prerender=true;としており、ビルドすると静的コンテンツとしてbuildディレクトリに出力される。
- `npm run dev`で開発サーバー（.env.development）
- `npm run build:prd`で本番用ビルド（.env.production）
- `npm run build:stg`でステージング用ビルド（.env.staging）
