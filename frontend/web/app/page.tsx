import { redirect } from 'next/navigation';

export default function RootPage() {
  // TODO: 認証実装後は、ログイン状態に応じて /login または /books にリダイレクト
  redirect('/books');
}
