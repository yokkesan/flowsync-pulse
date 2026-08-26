const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL ??
    'http://localhost:8081/api';

type ApiErrorResponse = {
    message?: string;
};

export type RealtimeConnectionTicketResponse = {
    ticket: string;
    expires_at: string;
};

export type CurrentWorkContextResponse = {
    user_id: number;
    project_id: number;
    project_name: string;
    repository_id: number;
    repository_name: string;
    task_id: number | null;
    task_key: string | null;
    task_name: string | null;
    branch_name: string;
    ticket_key: string | null;
    workspace_name: string | null;
    match_status: string;
    session_status: string;
    extension_active: boolean;
    started_at: string;
    last_heartbeat_at: string;
    ended_at: string | null;
    end_reason: string | null;
};

export class RealtimeApiError extends Error {
    status: number;

    constructor(
        message: string,
        status: number,
    ) {
        super(message);

        this.name = 'RealtimeApiError';
        this.status = status;
    }
}

async function parseErrorResponse(
    response: Response,
): Promise<RealtimeApiError> {
    let message =
        'リアルタイム接続の準備に失敗しました。';

    try {
        const data =
            (await response.json()) as ApiErrorResponse;

        if (
            typeof data.message === 'string' &&
            data.message.trim()
        ) {
            message = data.message.trim();
        }
    } catch {
        // JSON形式でない場合は
        // 共通メッセージを使用します。
    }

    return new RealtimeApiError(
        message,
        response.status,
    );
}

export async function createRealtimeConnectionTicket(
    accessToken: string,
): Promise<RealtimeConnectionTicketResponse> {
    const normalizedAccessToken =
        accessToken.trim();

    if (!normalizedAccessToken) {
        throw new RealtimeApiError(
            'ログイン情報がありません。',
            401,
        );
    }

    const response = await fetch(
        `${API_BASE_URL}/realtime/tickets`,
        {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                Authorization:
                    `Bearer ${normalizedAccessToken}`,
            },
        },
    );

    if (!response.ok) {
        throw await parseErrorResponse(response);
    }

    return (await response.json()) as
        RealtimeConnectionTicketResponse;
}

export async function getCurrentWorkContext(
    accessToken: string,
): Promise<CurrentWorkContextResponse | null> {
    const normalizedAccessToken =
        accessToken.trim();

    if (!normalizedAccessToken) {
        throw new RealtimeApiError(
            'ログイン情報がありません。',
            401,
        );
    }

    const response = await fetch(
        `${API_BASE_URL}/work-context/current`,
        {
            method: 'GET',
            headers: {
                Accept: 'application/json',
                Authorization:
                    `Bearer ${normalizedAccessToken}`,
            },
        },
    );

    if (response.status === 204) {
        return null;
    }

    if (!response.ok) {
        throw await parseErrorResponse(response);
    }

    return (await response.json()) as
        CurrentWorkContextResponse;
}

export function buildRealtimeWebSocketUrl(
    ticket: string,
): string {
    const normalizedTicket = ticket.trim();

    if (!normalizedTicket) {
        throw new RealtimeApiError(
            'リアルタイム接続チケットがありません。',
            400,
        );
    }

    const apiUrl =
        new URL(API_BASE_URL);

    apiUrl.protocol =
        apiUrl.protocol === 'https:'
            ? 'wss:'
            : 'ws:';

    apiUrl.pathname =
        `${apiUrl.pathname.replace(/\/+$/, '')}/realtime/ws`;

    apiUrl.search =
        new URLSearchParams({
            ticket: normalizedTicket,
        }).toString();

    return apiUrl.toString();
}