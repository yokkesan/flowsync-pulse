import type {
    CurrentUserResponse,
    LoginRequest,
    LoginResponse,
} from '../types/auth';

const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL ??
    'http://localhost:8081/api';

type ApiErrorResponse = {
    message?: string;
};

export class ApiError extends Error {
    status: number;

    constructor(
        message: string,
        status: number,
    ) {
        super(message);
        this.name = 'ApiError';
        this.status = status;
    }
}

async function parseErrorResponse(
    response: Response,
): Promise<ApiError> {
    let message = '通信処理に失敗しました。';

    try {
        const data =
            (await response.json()) as ApiErrorResponse;

        if (data.message) {
            message = data.message;
        }
    } catch {
        // JSON形式でないエラーレスポンスの場合は
        // 共通メッセージを使用します。
    }

    return new ApiError(
        message,
        response.status,
    );
}

export async function login(
    request: LoginRequest,
): Promise<LoginResponse> {
    const response = await fetch(
        `${API_BASE_URL}/auth/login`,
        {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Accept: 'application/json',
            },
            body: JSON.stringify(request),
        },
    );

    if (!response.ok) {
        throw await parseErrorResponse(response);
    }

    return (await response.json()) as LoginResponse;
}

export async function getCurrentUser(
    accessToken: string,
): Promise<CurrentUserResponse> {
    const response = await fetch(
        `${API_BASE_URL}/me`,
        {
            method: 'GET',
            headers: {
                Accept: 'application/json',
                Authorization: `Bearer ${accessToken}`,
            },
        },
    );

    if (!response.ok) {
        throw await parseErrorResponse(response);
    }

    return (await response.json()) as CurrentUserResponse;
}