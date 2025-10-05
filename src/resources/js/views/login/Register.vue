<template>
    <div class="login-form">
        <h1 class="login-form__title">ログイン</h1>

        <form @submit.prevent="submit">
            <div class="login-form__field">
                <label for="email">メールアドレス</label>
                <input id="email" type="email" v-model="form.email" required />
            </div>

            <div class="login-form__field">
                <label for="password">パスワード</label>
                <input id="password" type="password" v-model="form.password" required />
            </div>

            <button type="submit" class="login-form__btn">ログイン</button>
        </form>

        <p class="login-form__back">
            <a href="/" @click.prevent="backToTop">← 戻る</a>
        </p>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { router } from '@inertiajs/vue3'

const form = ref({
  email: '',
  password: ''
})

const submit = () => {
  router.visit('/office')
}

const backToTop = () => {
  router.visit('/')
}

import '@/../css/loginForm.scss'
</script>