import './bootstrap';
import { createApp, h } from 'vue';
import { createInertiaApp } from '@inertiajs/vue3';
import { resolvePageComponent } from 'laravel-vite-plugin/inertia-helpers';
import '../css/app.css';
import '../css/mainOffice.scss';
import '../css/sideDrawer.scss';
import '../css/loginTop.scss'; // ← LoginTop用（存在する場合）
import '../css/loginForm.scss'
import '../css/registerStep1.scss'
import '../css/registerStep2.scss'
import '../css/registerStep3.scss'

createInertiaApp({
    resolve: (name) =>
        resolvePageComponent(
            `./views/${name}.vue`,
            import.meta.glob('./views/**/*.vue')
        ),
    setup({ el, App, props, plugin }) {
        createApp({ render: () => h(App, props) })
            .use(plugin)
            .mount(el);
    },
});