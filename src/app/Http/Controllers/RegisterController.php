<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Facades\Mail;
use Inertia\Inertia;
use Carbon\Carbon;
use App\Models\User;
use App\Models\EmailVerification;
use App\Mail\VerificationCodeMail;
use Illuminate\Support\Facades\Auth;

class RegisterController extends Controller
{
    /**
     * Step1 表示
     */
    public function showStep1()
    {
        return Inertia::render('login/RegisterStep1');
    }

    /**
     * Step1 登録処理（ユーザー作成＆メール送信）
     */
    public function storeStep1(Request $request)
    {
        $request->validate([
            'name'     => 'required|string|max:255',
            'email'    => 'required|email|unique:users',
            'password' => 'required|string|min:6',
        ]);

        // ユーザー作成
        $user = User::create([
            'name'     => $request->name,
            'email'    => $request->email,
            'password' => Hash::make($request->password),
        ]);

        // 認証コード生成
        $code = sprintf('%06d', mt_rand(0, 999999));

        EmailVerification::create([
            'user_id'    => $user->id,
            'email'      => $user->email,
            'code'       => $code,
            'expires_at' => Carbon::now()->addMinutes(10),
        ]);

        // メール送信（ローカル用ならMailtrapなど）
        Mail::to($user->email)->send(new VerificationCodeMail($code));

        // Step1完了をセッションに保存（直打ち防止）
        $request->session()->put('register_step1_done', true);
        $request->session()->put('register_user_id', $user->id);

        return response()->json([
            'nextStep' => '/register/step2',
            'message'  => '認証コードをメールに送信しました。',
        ]);
    }

    /**
     * Step2 表示（Step1完了チェック付き）
     */
    public function showStep2(Request $request)
    {
        if (!$request->session()->has('register_step1_done')) {
            return redirect()->route('register.step1');
        }

        return Inertia::render('login/RegisterStep2');
    }

    /**
     * Step2 認証コード確認
     */
    public function storeStep2(Request $request)
    {
        $request->validate(['code' => 'required|string|size:6']);

        $verification = EmailVerification::where('code', $request->code)
            ->where('expires_at', '>', Carbon::now())
            ->latest()
            ->first();

        if (!$verification) {
            return response()->json([
                'errors' => ['code' => ['認証コードが無効または期限切れです。']],
            ], 422);
        }

        $user = User::find($verification->user_id);
        $user->update(['email_verified_at' => Carbon::now()]);
        $verification->delete();

        // Step2完了をセッションに保存
        $request->session()->put('register_step2_done', true);

        return response()->json([
            'nextStep' => '/register/step3',
            'message'  => 'メール認証が完了しました。',
        ]);
    }

    /**
     * Step3 表示（Step2完了チェック付き）
     */
    public function showStep3(Request $request)
    {
        if (!$request->session()->has('register_step2_done')) {
            return redirect()->route('register.step2');
        }

        return Inertia::render('login/RegisterStep3');
    }

    /**
     * Step3 GitHub / Repository 保存
     */
    public function storeStep3(Request $request)
    {
        $request->validate([
            'github'      => 'nullable|string|max:255',
            'repository'  => 'nullable|string|max:255',
        ]);

        $userId = $request->session()->get('register_user_id');
        $user   = User::find($userId);

        if (!$user) {
            return response()->json(['error' => 'User not found'], 404);
        }

        $user->update([
            'github'     => $request->github,
            'repository' => $request->repository,
        ]);

        // セッションをクリア
        $request->session()->forget([
            'register_step1_done',
            'register_step2_done',
            'register_user_id',
        ]);

        return response()->json([
            'nextStep' => '/login',
            'message'  => '登録が完了しました。',
        ]);
    }

    /**
     * ログイン処理
     */
    public function login(Request $request)
    {
        $credentials = $request->validate([
            'email' => ['required', 'email'],
            'password' => ['required'],
        ]);

        if (Auth::attempt($credentials)) {
            $request->session()->regenerate();

            return response()->json([
                'redirect' => '/office',
                'message' => 'ログイン成功',
            ]);
        }

        return response()->json([
            'errors' => [
                'email' => ['メールアドレスまたはパスワードが正しくありません。'],
            ],
        ], 422);
    }

    /**
     * ログアウト処理
     */
    public function logout(Request $request)
    {
        Auth::logout();
        $request->session()->invalidate();
        $request->session()->regenerateToken();

        return redirect('/login');
    }
}