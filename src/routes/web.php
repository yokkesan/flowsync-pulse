<?php

use Illuminate\Support\Facades\Route;
use Inertia\Inertia;
use App\Http\Controllers\RegisterController;

// ログイン処理
Route::get('/', function () { return Inertia::render('login/LoginTop'); })->name('login.top');
Route::get('/login', fn() => Inertia::render('login/Register'))->name('login.form');
Route::post('/login', [RegisterController::class, 'login'])->name('login.post');
Route::post('/logout', [RegisterController::class, 'logout'])->name('logout');
Route::get('/office', fn() => Inertia::render('MainOffice'))
    ->middleware('auth')
    ->name('office');
// 初期登録 step1
Route::get('/register/step1', [RegisterController::class, 'showStep1'])->name('register.step1');
Route::post('/register/step1', [RegisterController::class, 'storeStep1'])->name('register.store.step1');

// Step2
Route::get('/register/step2', [RegisterController::class, 'showStep2'])->name('register.step2');
Route::post('/register/step2', [RegisterController::class, 'storeStep2'])->name('register.store.step2');

// Step3
Route::get('/register/step3', [RegisterController::class, 'showStep3'])->name('register.step3');
Route::post('/register/step3', [RegisterController::class, 'storeStep3'])->name('register.store.step3');