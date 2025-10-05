<template>
    <div class="register-step2">
        <h1 class="register-step2__title">新規登録（Step 2）</h1>
        <p class="register-step2__desc">
            ご登録のメールアドレスに送信された認証コードを入力してください。
        </p>

        <form @submit.prevent="goNext">
            <div class="register-step2__field">
                <label for="code">認証コード</label>
                <input id="code" v-model="form.code" type="text" required maxlength="6" />
            </div>

            <button type="submit" class="register-step2__btn">次へ</button>
        </form>

        <p class="register-step2__back">
            <a href="/register/step1" @click.prevent="backToStep1">← 戻る</a>
        </p>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { router } from '@inertiajs/vue3'

const form = ref({
    code: ''
})

const goNext = () => {
    console.log('認証コード:', form.value.code)
    router.visit('/register/step3')
}

const backToStep1 = () => {
    router.visit('/register/step1')
}

// Sass 読み込み
import '@/../css/registerStep2.scss'
</script>