<?php

use Illuminate\Support\Facades\Route;
use Inertia\Inertia;

// ログイントップページ
Route::get('/', function () { return Inertia::render('login/LoginTop'); })->name('login.top');

// ログインフォーム
Route::get('/login', function () { return Inertia::render('login/Register'); })->name('login.form');

// 新規登録 Step1
Route::get('/register/step1', function () { return Inertia::render('login/RegisterStep1'); })->name('register.step1');
Route::get('/register/step2', function () { return Inertia::render('login/RegisterStep2'); })->name('register.step2');
Route::get('/register/step3', function () { return Inertia::render('login/RegisterStep3'); })->name('register.step3');