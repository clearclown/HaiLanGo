# 実装指示書: Flutter モバイルアプリ - 初期セットアップ

## 概要
HaiLanGo の Flutter モバイルアプリケーションの初期セットアップと基本構造の構築。iOS・Android 両対応のクロスプラットフォームアプリを作成する。

## 担当範囲
- **ディレクトリ**: `frontend/mobile/`
- **プラットフォーム**: iOS, Android
- **状態管理**: Riverpod
- **HTTP クライアント**: Dio

## 前提条件
- Flutter SDK 3.0+ がインストール済み
- Android Studio または Xcode がインストール済み
- Dart の基本知識

## 実装ステップ

### Step 1: Flutter プロジェクトの作成

```bash
cd frontend/
flutter create mobile
cd mobile
```

### Step 2: 依存関係の追加

**ファイル**: `frontend/mobile/pubspec.yaml`

```yaml
name: hailango_mobile
description: AI-Powered Language Learning Mobile App
publish_to: 'none'
version: 1.0.0+1

environment:
  sdk: '>=3.0.0 <4.0.0'

dependencies:
  flutter:
    sdk: flutter

  # State Management
  flutter_riverpod: ^2.4.0
  riverpod_annotation: ^2.2.0

  # HTTP Client
  dio: ^5.3.0

  # Storage
  shared_preferences: ^2.2.0
  flutter_secure_storage: ^9.0.0

  # UI Components
  google_fonts: ^6.1.0

  # Routing
  go_router: ^12.0.0

  # Utils
  uuid: ^4.1.0
  intl: ^0.18.1

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^3.0.0
  riverpod_generator: ^2.3.0
  build_runner: ^2.4.6

flutter:
  uses-material-design: true
  assets:
    - assets/images/
    - assets/icons/
```

### Step 3: ディレクトリ構造の作成

```bash
mkdir -p lib/core/{providers,models,services,utils}
mkdir -p lib/features/{auth,home,books,learning,review,settings}
mkdir -p lib/shared/{widgets,constants}
mkdir -p assets/{images,icons}
```

**最終的なディレクトリ構造**:

```
mobile/
├── lib/
│   ├── main.dart                # エントリーポイント
│   ├── app.dart                 # アプリルート
│   ├── core/
│   │   ├── providers/           # Riverpod プロバイダー
│   │   ├── models/              # データモデル
│   │   ├── services/            # API サービス
│   │   └── utils/               # ユーティリティ
│   ├── features/
│   │   ├── auth/                # 認証機能
│   │   ├── home/                # ホーム画面
│   │   ├── books/               # 書籍管理
│   │   ├── learning/            # 学習機能
│   │   ├── review/              # 復習機能
│   │   └── settings/            # 設定
│   └── shared/
│       ├── widgets/             # 共有ウィジェット
│       └── constants/           # 定数
└── assets/
    ├── images/
    └── icons/
```

### Step 4: 定数の定義

**ファイル**: `lib/shared/constants/api_constants.dart`

```dart
class ApiConstants {
  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );

  static const String apiVersion = 'v1';
  static const String apiPrefix = '/api/$apiVersion';

  // Endpoints
  static const String authRegister = '$apiPrefix/auth/register';
  static const String authLogin = '$apiPrefix/auth/login';
  static const String authRefresh = '$apiPrefix/auth/refresh';
  static const String authLogout = '$apiPrefix/auth/logout';

  static const String books = '$apiPrefix/books';
  static const String review = '$apiPrefix/review';
  static const String stats = '$apiPrefix/stats';
}
```

**ファイル**: `lib/shared/constants/colors.dart`

```dart
import 'package:flutter/material.dart';

class AppColors {
  // Primary Colors
  static const Color primary = Color(0xFF4A90E2);
  static const Color secondary = Color(0xFF50C878);
  static const Color accent = Color(0xFFFF6B6B);

  // Background
  static const Color background = Color(0xFFFFFFFF);
  static const Color backgroundSecondary = Color(0xFFF5F7FA);

  // Text
  static const Color textPrimary = Color(0xFF2C3E50);
  static const Color textSecondary = Color(0xFF7F8C8D);

  // Status
  static const Color success = Color(0xFF27AE60);
  static const Color warning = Color(0xFFF39C12);
  static const Color error = Color(0xFFE74C3C);
  static const Color info = Color(0xFF3498DB);

  // Border
  static const Color border = Color(0xFFE0E6ED);
}
```

### Step 5: HTTP クライアントの設定

**ファイル**: `lib/core/services/api_client.dart`

```dart
import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../../shared/constants/api_constants.dart';

class ApiClient {
  late final Dio _dio;
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  ApiClient() {
    _dio = Dio(BaseOptions(
      baseURL: ApiConstants.baseUrl,
      connectTimeout: const Duration(seconds: 30),
      receiveTimeout: const Duration(seconds: 30),
      headers: {
        'Content-Type': 'application/json',
      },
    ));

    // Add interceptors
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        // Add access token
        final token = await _storage.read(key: 'access_token');
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        return handler.next(options);
      },
      onError: (error, handler) async {
        // Handle 401 Unauthorized - refresh token
        if (error.response?.statusCode == 401) {
          if (await _refreshToken()) {
            // Retry original request
            return handler.resolve(await _retry(error.requestOptions));
          }
        }
        return handler.next(error);
      },
    ));
  }

  Dio get dio => _dio;

  Future<bool> _refreshToken() async {
    try {
      final refreshToken = await _storage.read(key: 'refresh_token');
      if (refreshToken == null) return false;

      final response = await _dio.post(
        ApiConstants.authRefresh,
        data: {'refresh_token': refreshToken},
      );

      if (response.statusCode == 200) {
        await _storage.write(
          key: 'access_token',
          value: response.data['access_token'],
        );
        return true;
      }
      return false;
    } catch (e) {
      return false;
    }
  }

  Future<Response<dynamic>> _retry(RequestOptions requestOptions) async {
    final options = Options(
      method: requestOptions.method,
      headers: requestOptions.headers,
    );
    return _dio.request<dynamic>(
      requestOptions.path,
      data: requestOptions.data,
      queryParameters: requestOptions.queryParameters,
      options: options,
    );
  }
}
```

### Step 6: Riverpod プロバイダーの設定

**ファイル**: `lib/core/providers/api_client_provider.dart`

```dart
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../services/api_client.dart';

final apiClientProvider = Provider<ApiClient>((ref) {
  return ApiClient();
});
```

### Step 7: アプリルートの作成

**ファイル**: `lib/app.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import 'features/auth/login_screen.dart';
import 'features/home/home_screen.dart';
import 'features/books/books_screen.dart';
import 'features/learning/learning_screen.dart';
import 'features/review/review_screen.dart';
import 'features/settings/settings_screen.dart';
import 'shared/constants/colors.dart';

class HaiLanGoApp extends ConsumerWidget {
  const HaiLanGoApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = GoRouter(
      initialLocation: '/login',
      routes: [
        GoRoute(
          path: '/login',
          builder: (context, state) => const LoginScreen(),
        ),
        GoRoute(
          path: '/home',
          builder: (context, state) => const HomeScreen(),
        ),
        GoRoute(
          path: '/books',
          builder: (context, state) => const BooksScreen(),
        ),
        GoRoute(
          path: '/learning/:bookId/:pageNumber',
          builder: (context, state) {
            final bookId = state.pathParameters['bookId']!;
            final pageNumber = int.parse(state.pathParameters['pageNumber']!);
            return LearningScreen(bookId: bookId, pageNumber: pageNumber);
          },
        ),
        GoRoute(
          path: '/review',
          builder: (context, state) => const ReviewScreen(),
        ),
        GoRoute(
          path: '/settings',
          builder: (context, state) => const SettingsScreen(),
        ),
      ],
    );

    return MaterialApp.router(
      title: 'HaiLanGo',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        primaryColor: AppColors.primary,
        scaffoldBackgroundColor: AppColors.background,
        textTheme: GoogleFonts.notoSansTextTheme(),
        colorScheme: ColorScheme.fromSeed(
          seedColor: AppColors.primary,
          secondary: AppColors.secondary,
        ),
        appBarTheme: const AppBarTheme(
          backgroundColor: AppColors.primary,
          foregroundColor: Colors.white,
          elevation: 0,
        ),
        elevatedButtonTheme: ElevatedButtonThemeData(
          style: ElevatedButton.styleFrom(
            backgroundColor: AppColors.primary,
            foregroundColor: Colors.white,
            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8),
            ),
          ),
        ),
      ),
      routerConfig: router,
    );
  }
}
```

### Step 8: エントリーポイントの作成

**ファイル**: `lib/main.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'app.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();

  runApp(
    const ProviderScope(
      child: HaiLanGoApp(),
    ),
  );
}
```

### Step 9: プレースホルダー画面の作成

各機能画面のプレースホルダーを作成します。

**ファイル**: `lib/features/auth/login_screen.dart`

```dart
import 'package:flutter/material.dart';

class LoginScreen extends StatelessWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text(
              'HaiLanGo',
              style: TextStyle(fontSize: 32, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 16),
            const Text('ログイン画面 - 実装予定'),
            const SizedBox(height: 32),
            ElevatedButton(
              onPressed: () {
                // TODO: 認証実装後に遷移
                // context.go('/home');
              },
              child: const Text('ログイン'),
            ),
          ],
        ),
      ),
    );
  }
}
```

**他の画面も同様に作成**:
- `lib/features/home/home_screen.dart`
- `lib/features/books/books_screen.dart`
- `lib/features/learning/learning_screen.dart`
- `lib/features/review/review_screen.dart`
- `lib/features/settings/settings_screen.dart`

### Step 10: ビルドと実行

```bash
# 依存関係のインストール
flutter pub get

# コード生成
flutter pub run build_runner build

# iOS で実行
flutter run -d ios

# Android で実行
flutter run -d android

# または特定のデバイスを指定
flutter devices
flutter run -d <device-id>
```

## テスト方法

1. **ビルドの確認**:
   ```bash
   flutter analyze
   flutter test
   ```

2. **実機またはシミュレータで起動**:
   ```bash
   flutter run
   ```

3. **確認項目**:
   - [ ] アプリが起動する
   - [ ] ログイン画面が表示される
   - [ ] エラーがない

## 完了条件

- [ ] Flutter プロジェクトが作成されている
- [ ] 依存関係がすべてインストールされている
- [ ] ディレクトリ構造が正しく作成されている
- [ ] API クライアントが設定されている
- [ ] Riverpod プロバイダーが設定されている
- [ ] ルーティングが設定されている
- [ ] すべての画面のプレースホルダーが作成されている
- [ ] iOS・Android 両方でビルドできる

## 次のステップ

1. **認証画面の実装**: `03_Flutter_Authentication.md` を参照
2. **ホーム画面の実装**: `04_Flutter_Home_Screen.md` を参照
3. **書籍画面の実装**: `05_Flutter_Books_Screen.md` を参照

## トラブルシューティング

### Flutter SDK のエラー
```bash
flutter doctor
```
で環境を確認。不足しているものをインストール。

### ビルドエラー
```bash
flutter clean
flutter pub get
flutter pub run build_runner build --delete-conflicting-outputs
```

### iOS ビルドエラー
```bash
cd ios
pod install
cd ..
```

## 参考資料

- [Flutter公式ドキュメント](https://flutter.dev/docs)
- [Riverpod公式ドキュメント](https://riverpod.dev/)
- [Go Router公式ドキュメント](https://pub.dev/packages/go_router)
