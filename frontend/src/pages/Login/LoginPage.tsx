import {
    useState,
    type FormEvent,
} from 'react';
import {
    useLocation,
    useNavigate,
} from 'react-router-dom';

import {
    ApiError,
    login,
} from '../../services/authApi';
import { saveAccessToken } from '../../services/authStorage';

type LocationState = {
    from?: string;
};

export function LoginPage() {
    const navigate = useNavigate();
    const location = useLocation();

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [errorMessage, setErrorMessage] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);

    async function handleSubmit(
        event: FormEvent<HTMLFormElement>,
    ) {
        event.preventDefault();

        if (isSubmitting) {
            return;
        }

        setErrorMessage('');
        setIsSubmitting(true);

        try {
            const response = await login({
                email,
                password,
            });

            saveAccessToken(response.access_token);

            const state =
                location.state as LocationState | null;

            const destination =
                state?.from?.startsWith('/office')
                    ? state.from
                    : '/office';

            navigate(destination, {
                replace: true,
            });
        } catch (error) {
            if (error instanceof ApiError) {
                setErrorMessage(error.message);
            } else {
                setErrorMessage(
                    '通信に失敗しました。時間をおいて再度お試しください。',
                );
            }
        } finally {
            setIsSubmitting(false);
        }
    }

    const submitButtonClassName = [
        'login-page__submit',
        isSubmitting
            ? 'login-page__submit--submitting'
            : '',
    ]
        .filter(Boolean)
        .join(' ');

    return (
        <main className="login-page">
            <section
                className="login-page__card"
                aria-labelledby="login-title"
            >
                <header className="login-page__header">
                    <p className="login-page__service-name">
                        FlowSync Pulse
                    </p>

                    <h1
                        id="login-title"
                        className="login-page__title"
                    >
                        ログイン
                    </h1>

                    <p className="login-page__description">
                        登録済みのユーザー情報で
                        バーチャルオフィスへ入室します。
                    </p>
                </header>

                <form
                    className="login-page__form"
                    onSubmit={handleSubmit}
                >
                    <div className="login-page__field">
                        <label
                            className="login-page__label"
                            htmlFor="email"
                        >
                            メールアドレス
                        </label>

                        <input
                            id="email"
                            className="login-page__input"
                            name="email"
                            type="email"
                            value={email}
                            autoComplete="email"
                            placeholder="example@example.com"
                            maxLength={255}
                            required
                            disabled={isSubmitting}
                            onChange={(event) => {
                                setEmail(event.target.value);
                            }}
                        />
                    </div>

                    <div className="login-page__field">
                        <label
                            className="login-page__label"
                            htmlFor="password"
                        >
                            パスワード
                        </label>

                        <input
                            id="password"
                            className="login-page__input"
                            name="password"
                            type="password"
                            value={password}
                            autoComplete="current-password"
                            placeholder="パスワードを入力"
                            minLength={8}
                            maxLength={72}
                            required
                            disabled={isSubmitting}
                            onChange={(event) => {
                                setPassword(event.target.value);
                            }}
                        />
                    </div>

                    {errorMessage && (
                        <p
                            className="login-page__error"
                            role="alert"
                        >
                            {errorMessage}
                        </p>
                    )}

                    <button
                        className={submitButtonClassName}
                        type="submit"
                        disabled={isSubmitting}
                    >
                        {isSubmitting
                            ? 'ログイン中...'
                            : 'ログイン'}
                    </button>
                </form>
            </section>
        </main>
    );
}