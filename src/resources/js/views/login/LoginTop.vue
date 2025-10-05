<template>
    <div class="login-top">
        <h1 class="title">FlowSync Pulse</h1>

        <div class="btn-area">
            <button @click="goLogin" class="btn btn-login">ログイン</button>
            <button @click="goRegister" class="btn btn-register">新規登録</button>
        </div>
    </div>
</template>

<script setup>
import { router as inertiaRouter } from '@inertiajs/vue3'

const goLogin = () => inertiaRouter.visit('/login')
const goRegister = () => inertiaRouter.visit('/register/step1')
</script>

<!-- スタイルは別ファイルで管理 -->