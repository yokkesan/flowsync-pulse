import { defineConfig } from 'vite';
import laravel from 'laravel-vite-plugin';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
    plugins: [
        laravel({
            input: ['resources/js/app.js'],
            refresh: true,
        }),
        vue(),
    ],
    server: {
        host: '0.0.0.0',  // コンテナ内からの接続を許可
        port: 5173,       // docker-compose で指定済みポート
        hmr: {
            host: 'localhost', // Laravelから見たホスト名
            port: 5173,        // 明示的に指定
        },
    },
});