# kifupの開発状況

## 進捗

### 完了したこと

1. kifupの開発を開始
2. kifup-frontendとkifup-apiの構成を設計
3. ページ・URLを設計
4. frontend用にSvelteKitのプロジェクトを新規作成し、Dockerで開発サーバを実行
5. api用にGinのプロジェクトを新規作成し、Dockerで開発サーバを実行
6. api側で、SQLiteのDB接続、ロギング、ミドルウェア、ハンドラー等を作成
7. frontend側で、基本的なUIを作成
8. api側でアカウント・セッション関係の処理を実装し、swaggerでAPI定義を共有
9. frontend側でアカウント作成・ログイン・セッション管理・ログアウトの機能を実装
10. api側で棋譜関連のmodelを作成
11. api側で棋譜関連のAPIエンドポイントを実装
12. frontend側で棋譜リストのデータ取得と表示を実装
13. frontendで棋譜の新規作成（局面から）を実装
14. frontendで棋譜情報の更新を実装
15. frontendで棋譜の指し手の更新を実装（一旦、分岐を含まない）
16. frontendで棋譜の閲覧を実装（一旦、分岐を含まない）
17. frontendで設定ページのアカウント情報更新を実装
18. アカウントの削除を実装（一旦、保有する棋譜は全削除）
19. AWS環境（ECS）でステージング環境をセットアップ
20. ステージング環境でkifupを起動する（Basic認証あり）
21. 動作モード（prd/stg/dev）を環境変数で設定する
22. いいね・コメント等のMOCを非表示
23. データベースをS3に定期バックアップと終了時バックアップをして永続化
24. APIサーバーでキャッシュコントロールを設定
25. 新規アカウント作成にメール認証を追加
26. 未実装の注意や開発メモなど、テストユーザー向けの情報を表示
27. テストユーザー募集の動画を公開
28. KIFフォーマットのデータから棋譜作成を実装（分岐なし）
29. CSAフォーマットのデータから棋譜作成を実装

### 今後の予定

1. モバイル向けのUI調整
2. 局面コメントの実装
3. いいね・感想コメントの設計
4. いいね・感想コメントの実装
5. 局面図ダウンロード機能の実装
6. 棋譜データの表示／ダウンロード
7. 分岐ありKIFデータの取り込み対応
8. 分岐ありKIFデータのダウンロード対応

## 資料

### 資料セット

- API定義.yml
  - api/swagger/api.yml
- api_ルーティングとミドルウェア.go
  - api/main.go
  - api/common/handler/middleware.go
  - api/common/handler/handler.go
  - api/common/handler/response.go
  - api/common/handler/pagination.go
- api_model_棋譜関連の型定義.go
  - api/service/model/Piece.go
  - api/service/model/BoardPosition.go
  - api/service/model/GameInfo.go
  - api/service/model/Kifu.go
  - api/service/model/KifuResponse.go
  - api/service/model/KifuMoveResponse.go
- frontend_棋譜関連の型定義.ts
  - frontend/src/lib/types/Kifu.ts
  - frontend/src/lib/types/Piece.ts
  - frontend/src/lib/types/BoardPosition.ts
- frontend_セッション管理.ts
  - frontend/src/lib/stores/session.ts
  - frontend/src/lib/apis/account.ts
  - frontend/src/lib/apis/session.ts

### 外部資料

- 棋譜フォーマット類
  - [KIF形式](http://kakinoki.o.oo7.jp/kif_format.html)
  - [CSA形式](http://www2.computer-shogi.org/protocol/record_v3.html)
  - [PSN形式](https://yaneuraou.yaneu.com/2021/06/22/what-is-the-shogi-game-format-psn/)
- 他サイト
  - [将棋MAP](https://shogimap.com/)
  - [将棋アイオー](https://shogi.io/)
  - [KENTO](https://www.kento-shogi.com/)
