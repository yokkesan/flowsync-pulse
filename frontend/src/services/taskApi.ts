import type {
    CreateTaskRequest,
    Task,
    TaskListResponse,
    UpdateTaskRequest,
} from '../types/task';

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

function validateTaskId(
    taskId: number,
): void {
    if (
        !Number.isSafeInteger(taskId) ||
        taskId <= 0
    ) {
        throw new Error(
            'タスクIDが正しくありません。',
        );
    }
}

export async function getTasks(
    accessToken: string,
    projectId: number,
): Promise<TaskListResponse> {
    validateProjectId(projectId);

    const response = await fetch(
        `${API_BASE_URL}/projects/${projectId}/tasks`,
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
            'タスク一覧の取得に失敗しました。',
        );
    }

    return (await response.json()) as TaskListResponse;
}

export async function getTask(
    accessToken: string,
    projectId: number,
    taskId: number,
): Promise<Task> {
    validateProjectId(projectId);
    validateTaskId(taskId);

    const response = await fetch(
        `${API_BASE_URL}/projects/${projectId}/tasks/${taskId}`,
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
            'タスク詳細の取得に失敗しました。',
        );
    }

    return (await response.json()) as Task;
}

export async function createTask(
    accessToken: string,
    projectId: number,
    request: CreateTaskRequest,
): Promise<Task> {
    validateProjectId(projectId);

    const response = await fetch(
        `${API_BASE_URL}/projects/${projectId}/tasks`,
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
            'タスク登録に失敗しました。',
        );
    }

    return (await response.json()) as Task;
}

export async function updateTask(
    accessToken: string,
    projectId: number,
    taskId: number,
    request: UpdateTaskRequest,
): Promise<Task> {
    validateProjectId(projectId);
    validateTaskId(taskId);

    const response = await fetch(
        `${API_BASE_URL}/projects/${projectId}/tasks/${taskId}`,
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
            'タスク編集に失敗しました。',
        );
    }

    return (await response.json()) as Task;
}

export async function deleteTask(
    accessToken: string,
    projectId: number,
    taskId: number,
): Promise<void> {
    validateProjectId(projectId);
    validateTaskId(taskId);

    const response = await fetch(
        `${API_BASE_URL}/projects/${projectId}/tasks/${taskId}`,
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
            'タスク削除に失敗しました。',
        );
    }
}