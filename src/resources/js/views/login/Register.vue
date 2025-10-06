<template>
    <div class="login-form">
        <h1 class="login-form__title">ログイン</h1>

        <form @submit.prevent="submit">
            <div class="login-form__field">
                <label for="email">メールアドレス</label>
                <input id="email" type="email" v-model="form.email" required />
                <p v-if="errors.email" class="login-form__error">{{ errors.email[0] }}</p>
            </div>

            <div class="login-form__field">
                <label for="password">パスワード</label>
                <input id="password" type="password" v-model="form.password" required />
                <p v-if="errors.password" class="login-form__error">{{ errors.password[0] }}</p>
            </div>

            <button type="submit" class="login-form__btn" :disabled="loading">
                {{ loading ? 'ログイン中...' : 'ログイン' }}
            </button>
        </form>

        <p class="login-form__back">
            <a href="/" @click.prevent="backToTop">← 戻る</a>
        </p>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { router } from '@inertiajs/vue3'

const form = ref({ email: '', password: '' })
const errors = ref({})
const loading = ref(false)

const submit = async () => {
    errors.value = {}
    loading.value = true

    try {
        const res = await axios.post('/login', form.value)
        if (res.data?.redirect) {
            router.visit(res.data.redirect)
        }
    } catch (e) {
        if (e.response?.status === 422) {
            errors.value = e.response.data.errors
        } else {
            console.error(e)
        }
    } finally {
        loading.value = false
    }
}

const backToTop = () => {
    router.visit('/')
}

import '@/../css/loginForm.scss'
</script>