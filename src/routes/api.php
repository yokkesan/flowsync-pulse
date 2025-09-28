<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;
use App\Http\Controllers\BranchController;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
| ここにAPIエンドポイントを定義します。
| これらは自動で "api" ミドルウェアグループが適用されます。
| 例: http://localhost/api/branches
|--------------------------------------------------------------------------
*/

// ブランチ情報を保存（VSCode拡張からPOST）
Route::post('/branches', [BranchController::class, 'store']);

// ブランチ一覧を取得（Vue.jsダッシュボード用 GET）
Route::get('/branches', [BranchController::class, 'index']);