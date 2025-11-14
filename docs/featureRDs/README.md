# 🚨 CRITICAL - Feature Requirements Documents (RDs)

**PM**: あなた（PM）
**最終更新**: 2025-11-14 13:35
**ステータス**: 🔴 **CRITICAL FAILURES DETECTED**

---

## ⚠️ **重大な警告**

**既存のREADME.mdは虚偽情報を含んでいました。**

- ❌ 「すべての機能が100%実装完了」→ **嘘**
- ✅ 実際の実装率: **約40%**（フロントエンドのみ）
- ✅ バックエンドAPI: **ほぼ未実装**

**このような不正確な報告は二度と許されません。**

---

## 📂 ディレクトリ構造

```
featureRDs/
├── README.md                          # このファイル（PM用マスターガイド）
├── CRITICAL_01_Books_API.md           # 🔴 P0 - 即座実装必須
├── CRITICAL_02_Review_API.md          # 🔴 P0 - 即座実装必須
├── CRITICAL_03_Router_Integration.md  # 🔴 P0 - 即座実装必須
└── archives/                          # 過去のドキュメント（参考のみ）
```

---

## 🔴 CRITICAL ISSUES（即座に対応必須）

**現在の実装状況: UNACCEPTABLE**

### 実装済み
- ✅ フロントエンド: Books, Upload, Review, Settings ページ（95%）
- ✅ 認証API（Backend）
- ✅ アップロードAPI（Backend、部分的）

### **未実装（CRITICAL）**
- ❌ Books API（Backend）← **本の作成が失敗中**
- ❌ Review API（Backend）← **復習ページがエラー**
- ❌ Router Integration（Backend）← **404エラー多発**
- ❌ E2Eテストの35%が失敗

---

## 📋 P0 - CRITICAL Feature RDs（即座実装必須）

| No | Feature | ファイル | 見積 | ステータス |
|----|---------|----------|------|-----------|
| 1 | Books API | [`CRITICAL_01_Books_API.md`](CRITICAL_01_Books_API.md) | 4-6h | ❌ 未着手 |
| 2 | Review API (SRS) | [`CRITICAL_02_Review_API.md`](CRITICAL_02_Review_API.md) | 6-8h | ❌ 未着手 |
| 3 | Router Integration | [`CRITICAL_03_Router_Integration.md`](CRITICAL_03_Router_Integration.md) | 2-3h | ❌ 未着手 |

**期限**: 本日中（2025-11-14）に着手、3日以内（2025-11-17）に完了

---

## 📊 実装優先順位と役割分担

### 推奨チーム構成

#### 4人チーム:
- **Backend Engineer A** (Senior): Router Integration + Books API（6-9h）
- **Backend Engineer B** (Mid): Review API + SRS（6-8h）
- **Frontend Engineer**: E2Eテスト修正 + Stats Dashboard（8-13h）
- **Full-stack Engineer**: OCR API + Learning API（8-12h）

#### 2人チーム:
- **Backend Lead**: すべてのCRITICAL（16-23h、2-3日）
- **Frontend/Full-stack**: E2Eテスト + Stats + Learning API（12-19h、1.5-2.5日）

詳細は各RD内の「役割分担」セクションを参照。

---

## 🎯 完了基準（Project DoD）

### バックエンド
- [ ] すべてのCRITICAL APIが実装
- [ ] すべてのハンドラーがルーターに登録
- [ ] すべてのエンドポイントが動作
- [ ] ユニットテストカバレッジ80%以上

### フロントエンド
- [ ] すべてのページが正常動作
- [ ] E2Eテストが100%パス
- [ ] APIエラーがゼロ

### 統合
- [ ] フロントエンド→バックエンド通信が全成功
- [ ] 認証フローが正常動作
- [ ] ファイルアップロードが動作
- [ ] 復習機能が動作

---

## 📖 実装ガイドライン

### 必須事項

1. **各エンジニアは実装前に該当RDを完全に読むこと**
2. **仕様を勝手に変更しないこと**
3. **テストを必ず書くこと**（カバレッジ80%以上）
4. **エラーハンドリングを省略しないこと**
5. **実装完了後は動作確認を行うこと**

### 禁止事項

- ❌ 推測で実装する
- ❌ テストを書かない
- ❌ ドキュメントを読まない
- ❌ PMに確認せず仕様変更
- ❌ セキュリティチェックを省略

### コミットメッセージ規約
```
feat(books): implement Books API endpoints
fix(review): correct SRS algorithm calculation
docs(api): update API documentation
test(books): add unit tests for BookRepository
```

---

## 🚀 開発環境セットアップ

```bash
# プロジェクトルート
cd /home/ablaze/Projects/haiLanGo

# バックエンド起動
cd backend
go mod download
go run cmd/server/main.go

# フロントエンド起動（別ターミナル）
cd frontend/web
pnpm install
pnpm dev

# データベース起動（別ターミナル）
podman-compose up -d
```

### 動作確認
```bash
# Health Check
curl http://localhost:8080/health

# Books API
curl -X GET http://localhost:8080/api/v1/books -H "Authorization: Bearer {token}"

# Review API
curl -X GET http://localhost:8080/api/v1/review/stats -H "Authorization: Bearer {token}"
```

---

## 📞 報告とコミュニケーション

### 日次報告（必須）

毎日の終業時にPMに報告：

```
【日時】: 2025-11-14 18:00
【担当者】: Backend Engineer A
【実装した内容】: Router Integration 70%完了
【完了したタスク】: ✅ router.go書き直し
【ブロッカー】: なし
【明日の予定】: Books API実装
```

### 質問・不明点

**不明点があれば即座にPMに確認すること。推測で実装するな。**

---

## ⏰ タイムライン

### Day 1（本日 - 2025-11-14）
- 9:00-12:00: CRITICAL_03 Router Integration完了
- 13:00-18:00: CRITICAL_01 Books API完了
- 18:00-21:00: CRITICAL_02 Review API 50%完了

### Day 2（2025-11-15）
- 9:00-12:00: CRITICAL_02 Review API完了
- 13:00-15:00: E2Eテスト修正
- 15:00-18:00: 統合テスト・動作確認

### Day 3（2025-11-16）
- 9:00-12:00: バグ修正
- 13:00-15:00: ドキュメント更新
- 15:00-17:00: 最終確認・リリース準備

**期限**: 2025-11-17まで（3日以内）

---

## 📖 参考ドキュメント

- [要件定義書](../requirements_definition.md)
- [UI/UX設計書](../ui_ux_design_document.md)
- [モック構築戦略](../mocking_strategy.md)
- [API統合提案書](../api_integration_proposal.md)

---

## 🔄 更新履歴

| 日付 | 内容 | 更新者 |
|------|------|--------|
| 2025-11-14 | CRITICAL状況を反映、虚偽情報を修正 | PM (Claude Code) |
| 2025-11-13 | 初版作成（不正確） | Previous Session |

---

## 最後に

**このプロジェクトの成功はチーム全員の責任である。**

**言い訳は不要。結果を出せ。**

**期限**: 3日以内（2025-11-17まで）
**品質基準**: 妥協なし

---

**質問がある場合は即座にPMに連絡すること。**

**Good Luck.**
