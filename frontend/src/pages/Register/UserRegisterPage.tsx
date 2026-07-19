import {
    useState,
    type ChangeEvent,
    type FormEvent,
} from 'react';
import {
    Navigate,
    useNavigate,
    useParams,
} from 'react-router-dom';

type UserRegisterForm = {
    displayName: string;
    email: string;
    password: string;
    passwordConfirmation: string;
};

type UserRegisterResponse = {
    user_id: number;
    company_id: number;
    display_name: string;
    email: string;
    role: 'owner';
};

type ApiErrorResponse = {
    message?: string;
};

type CurrentUser = {
    userId: number;
    companyId: number;
    displayName: string;
    email: string;
    role: 'owner';
};

const initialForm: UserRegisterForm = {
    displayName: '',
    email: '',
    password: '',
    passwordConfirmation: '',
};

const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL ??
    'http://localhost:8081';

export function UserRegisterPage() {
    const navigate = useNavigate();

    const { companyId } = useParams<{
        companyId: string;
    }>();

    const parsedCompanyId = Number(companyId);

    const [form, setForm] =
        useState<UserRegisterForm>(initialForm);

    const [isSubmitting, setIsSubmitting] =
        useState(false);

    const [errorMessage, setErrorMessage] =
        useState('');

    if (
        !Number.isInteger(parsedCompanyId) ||
        parsedCompanyId <= 0
    ) {
        return (
            <Navigate
                to="/register/company"
                replace
            />
        );
    }

    const handleChange = (
        event: ChangeEvent<HTMLInputElement>,
    ) => {
        const { name, value } = event.target;

        setForm((currentForm) => ({
            ...currentForm,
            [name]: value,
        }));
    };

    const validateForm = (): string | null => {
        if (!form.displayName.trim()) {
            return '表示名を入力してください。';
        }

        if (!form.email.trim()) {
            return 'メールアドレスを入力してください。';
        }

        if (!form.password) {
            return 'パスワードを入力してください。';
        }

        if (form.password.length < 8) {
            return 'パスワードは8文字以上で入力してください。';
        }

        if (!form.passwordConfirmation) {
            return 'パスワード確認を入力してください。';
        }

        if (
            form.password !==
            form.passwordConfirmation
        ) {
            return 'パスワード確認が一致しません。';
        }

        return null;
    };

    const handleSubmit = async (
        event: FormEvent<HTMLFormElement>,
    ) => {
        event.preventDefault();
        setErrorMessage('');

        const validationError = validateForm();

        if (validationError) {
            setErrorMessage(validationError);
            return;
        }

        setIsSubmitting(true);

        try {
            const response = await fetch(
                `${API_BASE_URL}/api/companies/${parsedCompanyId}/users`,
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        display_name:
                            form.displayName.trim(),
                        email:
                            form.email
                                .trim()
                                .toLowerCase(),
                        password: form.password,
                        password_confirmation:
                            form.passwordConfirmation,
                    }),
                },
            );

            const responseBody = (await response.json()) as
                | UserRegisterResponse
                | ApiErrorResponse;

            if (!response.ok) {
                const apiError =
                    responseBody as ApiErrorResponse;

                throw new Error(
                    apiError.message ??
                        'ユーザー登録に失敗しました。',
                );
            }

            const registeredUser =
                responseBody as UserRegisterResponse;

            const currentUser: CurrentUser = {
                userId: registeredUser.user_id,
                companyId: registeredUser.company_id,
                displayName:
                    registeredUser.display_name,
                email: registeredUser.email,
                role: registeredUser.role,
            };

            sessionStorage.setItem(
                'currentUser',
                JSON.stringify(currentUser),
            );

            navigate('/office');
        } catch (error) {
            setErrorMessage(
                error instanceof Error
                    ? error.message
                    : 'ユーザー登録に失敗しました。',
            );
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <main className="register-page">
            <section className="register-page__card">
                <header className="register-page__header">
                    <p className="register-page__service-name">
                        FlowSync Pulse
                    </p>

                    <p className="register-page__step">
                        STEP 2 / 2
                    </p>

                    <h1 className="register-page__title">
                        ユーザー登録
                    </h1>

                    <p className="register-page__description">
                        最初の管理ユーザーを登録します。
                    </p>
                </header>

                <form
                    className="register-form"
                    onSubmit={handleSubmit}
                    noValidate
                >
                    {errorMessage && (
                        <p
                            className="register-form__message register-form__message--error"
                            role="alert"
                        >
                            {errorMessage}
                        </p>
                    )}

                    <div className="register-form__field">
                        <label
                            className="register-form__label"
                            htmlFor="displayName"
                        >
                            表示名
                        </label>

                        <input
                            id="displayName"
                            name="displayName"
                            type="text"
                            className="register-form__input"
                            value={form.displayName}
                            onChange={handleChange}
                            autoComplete="name"
                            maxLength={100}
                            disabled={isSubmitting}
                            required
                        />
                    </div>

                    <div className="register-form__field">
                        <label
                            className="register-form__label"
                            htmlFor="email"
                        >
                            メールアドレス
                        </label>

                        <input
                            id="email"
                            name="email"
                            type="email"
                            className="register-form__input"
                            value={form.email}
                            onChange={handleChange}
                            autoComplete="email"
                            maxLength={255}
                            disabled={isSubmitting}
                            required
                        />
                    </div>

                    <div className="register-form__field">
                        <label
                            className="register-form__label"
                            htmlFor="password"
                        >
                            パスワード
                        </label>

                        <input
                            id="password"
                            name="password"
                            type="password"
                            className="register-form__input"
                            value={form.password}
                            onChange={handleChange}
                            autoComplete="new-password"
                            minLength={8}
                            maxLength={72}
                            disabled={isSubmitting}
                            required
                        />
                    </div>

                    <div className="register-form__field">
                        <label
                            className="register-form__label"
                            htmlFor="passwordConfirmation"
                        >
                            パスワード確認
                        </label>

                        <input
                            id="passwordConfirmation"
                            name="passwordConfirmation"
                            type="password"
                            className="register-form__input"
                            value={form.passwordConfirmation}
                            onChange={handleChange}
                            autoComplete="new-password"
                            minLength={8}
                            maxLength={72}
                            disabled={isSubmitting}
                            required
                        />
                    </div>

                    <div className="register-form__actions">
                        <button
                            type="button"
                            className="register-form__button register-form__button--secondary"
                            onClick={() =>
                                navigate(
                                    '/register/company',
                                )
                            }
                            disabled={isSubmitting}
                        >
                            戻る
                        </button>

                        <button
                            type="submit"
                            className="register-form__button register-form__button--primary"
                            disabled={isSubmitting}
                        >
                            {isSubmitting
                                ? 'ユーザーを登録しています...'
                                : 'ユーザーを登録する'}
                        </button>
                    </div>
                </form>
            </section>
        </main>
    );
}