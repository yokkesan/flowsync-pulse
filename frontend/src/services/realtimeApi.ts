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

    const apiUrl = new URL(API_BASE_URL);

    apiUrl.protocol =
        apiUrl.protocol === 'https:'
            ? 'wss:'
            : 'ws:';

    apiUrl.pathname =
        `${apiUrl.pathname.replace(/\/+$/, '')}/realtime/ws`;

    apiUrl.search = new URLSearchParams({
        ticket: normalizedTicket,
    }).toString();

    return apiUrl.toString();
}