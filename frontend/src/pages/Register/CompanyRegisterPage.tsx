import {
    useState,
    type ChangeEvent,
    type FormEvent,
} from 'react';
import { useNavigate } from 'react-router-dom';

type CompanyRegisterForm = {
    companyName: string;
    companySlug: string;
};

type CompanyRegisterResponse = {
    company_id: number;
    name: string;
    slug: string;
};

type ApiErrorResponse = {
    message?: string;
};

const initialForm: CompanyRegisterForm = {
    companyName: '',
    companySlug: '',
};

const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL ??
    'http://localhost:8081';

export function CompanyRegisterPage() {
    const navigate = useNavigate();

    const [form, setForm] =
        useState<CompanyRegisterForm>(initialForm);

    const [isSubmitting, setIsSubmitting] =
        useState(false);

    const [errorMessage, setErrorMessage] =
        useState('');

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
        if (!form.companyName.trim()) {
            return '会社名を入力してください。';
        }

        if (!form.companySlug.trim()) {
            return '会社スラッグを入力してください。';
        }

        const companySlugPattern =
            /^[a-z0-9]+(?:-[a-z0-9]+)*$/;

        if (
            !companySlugPattern.test(
                form.companySlug.trim().toLowerCase(),
            )
        ) {
            return '会社スラッグは半角英小文字・数字・ハイフンで入力してください。';
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
                `${API_BASE_URL}/api/companies`,
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        name: form.companyName.trim(),
                        slug: form.companySlug
                            .trim()
                            .toLowerCase(),
                    }),
                },
            );

            const responseBody = (await response.json()) as
                | CompanyRegisterResponse
                | ApiErrorResponse;

            if (!response.ok) {
                const apiError =
                    responseBody as ApiErrorResponse;

                throw new Error(
                    apiError.message ??
                    '会社登録に失敗しました。',
                );
            }

            const company =
                responseBody as CompanyRegisterResponse;

            navigate(
                `/register/company/${company.company_id}/user`,
            );
        } catch (error) {
            setErrorMessage(
                error instanceof Error
                    ? error.message
                    : '会社登録に失敗しました。',
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
                        STEP 1 / 2
                    </p>

                    <h1 className="register-page__title">
                        会社登録
                    </h1>

                    <p className="register-page__description">
                        FlowSync Pulseを利用する会社の情報を登録します。
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
                            htmlFor="companyName"
                        >
                            会社名
                        </label>

                        <input
                            id="companyName"
                            name="companyName"
                            type="text"
                            className="register-form__input"
                            value={form.companyName}
                            onChange={handleChange}
                            autoComplete="organization"
                            maxLength={150}
                            disabled={isSubmitting}
                            required
                        />
                    </div>

                    <div className="register-form__field">
                        <label
                            className="register-form__label"
                            htmlFor="companySlug"
                        >
                            会社スラッグ
                        </label>

                        <input
                            id="companySlug"
                            name="companySlug"
                            type="text"
                            className="register-form__input"
                            value={form.companySlug}
                            onChange={handleChange}
                            placeholder="flowsync"
                            autoCapitalize="none"
                            autoCorrect="off"
                            maxLength={100}
                            disabled={isSubmitting}
                            required
                        />

                        <p className="register-form__help">
                            半角英小文字・数字・ハイフンを使用できます。
                        </p>
                    </div>

                    <button
                        type="submit"
                        className="register-form__button register-form__button--primary"
                        disabled={isSubmitting}
                    >
                        {isSubmitting
                            ? '会社を登録しています...'
                            : '会社を登録して次へ'}
                    </button>
                </form>
            </section>
        </main>
    );
}