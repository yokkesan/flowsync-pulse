import type {
    CreateProjectRequest,
    Project,
    ProjectListResponse,
    UpdateProjectRequest,
} from '../types/project';

import { ApiError } from './authApi';

const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL ??
    'http://localhost:8081/api';

type ApiErrorResponse = {
    message?: string;
};

async function parseErrorResponse(
    response: Response,
    defaultMessage: string,
): Promise<ApiError> {
    let message = defaultMessage;

    try {
        const data =
            (await response.json()) as ApiErrorResponse;

        if (data.message) {
            message = data.message;
        }
    } catch {
        // JSON形式でない場合は既定メッセージを使用します。
    }

    return new ApiError(
        message,
        response.status,
    );
}

function createAuthorizationHeaders(
    accessToken: string,
): HeadersInit {
    return {
        Accept: 'application/json',
        Authorization: `Bearer ${accessToken}`,
    };
}

function validateProjectId(
    projectId: number,
): void {
    if (
        !Number.isSafeInteger(projectId) ||
        projectId <= 0
    ) {
        throw new Error(
            'プロジェクトIDが正しくありません。',
        );
    }
}

export async function getProjects(
    accessToken: string,
): Promise<ProjectListResponse> {
    const response = await fetch(
        `${API_BASE_URL}/projects`,
        {
            method: 'GET',
            headers:
                createAuthorizationHeaders(
                    accessToken,
                ),
        },
    );

    if (!response.ok) {
        throw await parseErrorResponse(
            response,
            'プロジェクト一覧の取得に失敗しました。',
        );
    }

    return (await response.json()) as ProjectListResponse;
}

export async function getProject(
    accessToken: string,
    projectId: number,
): Promise<Project> {
    validateProjectId(projectId);

    const response = await fetch(
        `${API_BASE_URL}/projects/${projectId}`,
        {
            method: 'GET',
            headers:
                createAuthorizationHeaders(
                    accessToken,
                ),
        },
    );

    if (!response.ok) {
        throw await parseErrorResponse(
            response,
            'プロジェクト詳細の取得に失敗しました。',
        );
    }

    return (await response.json()) as Project;
}

export async function createProject(
    accessToken: string,
    request: CreateProjectRequest,
): Promise<Project> {
    const response = await fetch(
        `${API_BASE_URL}/projects`,
        {
            method: 'POST',
            headers: {
                ...createAuthorizationHeaders(
                    accessToken,
                ),
                'Content-Type':
                    'application/json',
            },
            body: JSON.stringify(request),
        },
    );

    if (!response.ok) {
        throw await parseErrorResponse(
            response,
            'プロジェクト登録に失敗しました。',
        );
    }

    return (await response.json()) as Project;
}

export async function updateProject(
    accessToken: string,
    projectId: number,
    request: UpdateProjectRequest,
): Promise<Project> {
    validateProjectId(projectId);

    const response = await fetch(
        `${API_BASE_URL}/projects/${projectId}`,
        {
            method: 'PUT',
            headers: {
                ...createAuthorizationHeaders(
                    accessToken,
                ),
                'Content-Type':
                    'application/json',
            },
            body: JSON.stringify(request),
        },
    );

    if (!response.ok) {
        throw await parseErrorResponse(
            response,
            'プロジェクト編集に失敗しました。',
        );
    }

    return (await response.json()) as Project;
}

export async function deleteProject(
    accessToken: string,
    projectId: number,
): Promise<void> {
    validateProjectId(projectId);

    const response = await fetch(
        `${API_BASE_URL}/projects/${projectId}`,
        {
            method: 'DELETE',
            headers:
                createAuthorizationHeaders(
                    accessToken,
                ),
        },
    );

    if (!response.ok) {
        throw await parseErrorResponse(
            response,
            'プロジェクト削除に失敗しました。',
        );
    }
}