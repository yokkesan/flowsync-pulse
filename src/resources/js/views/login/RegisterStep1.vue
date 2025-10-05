<template>
    <div class="register-step1">
        <h1 class="register-step1__title">新規登録（Step 1）</h1>

        <form @submit.prevent="goNext">
            <div class="register-step1__field">
                <label for="name">ユーザー名</label>
                <input id="name" v-model="form.name" type="text" required />
            </div>

            <div class="register-step1__field">
                <label for="email">メールアドレス</label>
                <input id="email" v-model="form.email" type="email" required />
            </div>

            <div class="register-step1__field">
                <label for="password">パスワード</label>
                <input id="password" v-model="form.password" type="password" required />
            </div>

            <button type="submit" class="register-step1__btn">次へ</button>
        </form>

        <p class="register-step1__back">
            <a href="/" @click.prevent="backToTop">← 戻る</a>
        </p>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { router } from '@inertiajs/vue3'

// フォームデータ
const form = ref({
    name: '',
    email: '',
    password: ''
})

// Step2へ進む
const goNext = () => {
    console.log('Step1 入力内容', form.value)
    router.visit('/register/step2')
}

// 戻る
const backToTop = () => {
    router.visit('/')
}

// Sass読み込み
import '@/../css/registerStep1.scss'
</script>