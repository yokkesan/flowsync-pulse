import type {
    CompanyUserListResponse,
} from '../types/user';

import { ApiError } from './authApi';

const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL ??
    'http://localhost:8081/api';

type ApiErrorResponse = {
    message?: string;
};

async function parseErrorResponse(
    response: Response,
): Promise<ApiError> {
    let message =
        '所属ユーザー一覧の取得に失敗しました。';

    try {
        const data =
            (await response.json()) as ApiErrorResponse;

        if (data.message) {
            message = data.message;
        }
    } catch {
        // JSON形式でない場合は共通メッセージを使用します。
    }

    return new ApiError(
        message,
        response.status,
    );
}

export async function getCompanyUsers(
    accessToken: string,
): Promise<CompanyUserListResponse> {
    const response = await fetch(
        `${API_BASE_URL}/users`,
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

    return (await response.json()) as CompanyUserListResponse;
}