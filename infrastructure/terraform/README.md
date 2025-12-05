# HaiLanGo Infrastructure - Terraform

Google Cloud Platform インフラストラクチャを管理するTerraform構成です。

## 前提条件

1. [Terraform](https://www.terraform.io/downloads) (>= 1.0)
2. [Google Cloud SDK](https://cloud.google.com/sdk/install)
3. GCPプロジェクトへのアクセス権限

## セットアップ

### 1. Google Cloud認証

```bash
# gcloud CLIでログイン
gcloud auth login

# Application Default Credentialsを設定
gcloud auth application-default login

# プロジェクトを設定
gcloud config set project hailango
```

### 2. 変数ファイルの作成

```bash
# サンプルをコピー
cp terraform.tfvars.example terraform.tfvars

# エディタで編集
vim terraform.tfvars
```

### 3. Billing Account IDの取得

```bash
# Billing Account一覧を表示
gcloud billing accounts list

# 出力例:
# ACCOUNT_ID            NAME                  OPEN  MASTER_ACCOUNT_ID
# 012345-6789AB-CDEF01  My Billing Account    True
```

### 4. Terraformの実行

```bash
# 初期化
terraform init

# 計画を確認
terraform plan

# 適用
terraform apply

# 出力の確認
terraform output
```

## 管理されるリソース

### APIs
- Cloud Vision API (OCR)
- Cloud Text-to-Speech API
- Cloud Speech-to-Text API
- Cloud Translation API
- Cloud Storage
- Secret Manager
- IAM

### Storage Buckets
- `hailango-user-uploads-{env}` - ユーザーアップロード（書籍、画像）
- `hailango-audio-cache-{env}` - TTS音声キャッシュ（7日で自動削除）
- `hailango-ocr-cache-{env}` - OCR結果キャッシュ（30日で自動削除）

### IAM
- サービスアカウント: `hailango-api-sa`
- カスタムロール: `hailangoTtsSttUser`
- APIキー（制限付き）

### Budget & Monitoring
- 月額予算アラート（50%, 80%, 90%, 100%）
- APIクォータ監視ダッシュボード
- クォータ使用率アラート

## コスト管理

### 予算設定

```hcl
# terraform.tfvars
monthly_budget_amount   = 50  # $50/月
budget_alert_thresholds = [0.5, 0.8, 0.9, 1.0]
```

### コスト最適化のヒント

1. **キャッシュの活用**
   - TTS結果は7日間キャッシュ
   - OCR結果は30日間キャッシュ

2. **ストレージ階層**
   - 30日後: NEARLINE
   - 90日後: COLDLINE

3. **レート制限**
   - アプリケーション側でAPI呼び出しを制限

## 環境変数の取得

Terraformの出力から.envファイル用の内容を取得：

```bash
# 環境変数の内容を表示
terraform output -raw env_file_content

# サービスアカウントキーを保存
terraform output -raw service_account_key_json | base64 -d > credentials.json
```

## トラブルシューティング

### APIが有効化されない

```bash
# 手動でAPIを有効化
gcloud services enable vision.googleapis.com
gcloud services enable texttospeech.googleapis.com
gcloud services enable speech.googleapis.com
```

### 権限エラー

```bash
# 必要な権限を確認
gcloud projects get-iam-policy hailango

# 自分のアカウントにオーナー権限を付与
gcloud projects add-iam-policy-binding hailango \
  --member="user:your-email@example.com" \
  --role="roles/owner"
```

### Budget APIエラー

予算APIにはBilling Account Admin権限が必要です：

```bash
gcloud billing accounts get-iam-policy BILLING_ACCOUNT_ID
```

## クリーンアップ

**注意**: 本番環境では実行しないでください。

```bash
# すべてのリソースを削除
terraform destroy
```

## 参考リンク

- [Google Cloud Terraform Provider](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- [Cloud Vision API料金](https://cloud.google.com/vision/pricing)
- [Cloud TTS料金](https://cloud.google.com/text-to-speech/pricing)
- [Cloud STT料金](https://cloud.google.com/speech-to-text/pricing)
