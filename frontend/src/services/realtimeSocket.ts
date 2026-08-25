import {
    buildRealtimeWebSocketUrl,
    createRealtimeConnectionTicket,
} from './realtimeApi';

import type {
    RealtimeEvent,
} from '../types/realtime';

type RealtimeSocketOptions = {
    accessToken: string;
    onEvent: (event: RealtimeEvent) => void;
    onOpen?: () => void;
    onClose?: () => void;
    onError?: () => void;
};

export async function connectRealtimeSocket(
    options: RealtimeSocketOptions,
): Promise<WebSocket> {
    const ticketResponse =
        await createRealtimeConnectionTicket(
            options.accessToken,
        );

    const socketUrl =
        buildRealtimeWebSocketUrl(
            ticketResponse.ticket,
        );

    const socket =
        new WebSocket(socketUrl);

    socket.addEventListener(
        'open',
        () => {
            options.onOpen?.();
        },
    );

    socket.addEventListener(
        'message',
        (messageEvent) => {
            if (
                typeof messageEvent.data !==
                'string'
            ) {
                return;
            }

            try {
                const event =
                    JSON.parse(
                        messageEvent.data,
                    ) as RealtimeEvent;

                options.onEvent(event);
            } catch {
                // 不正なJSONは無視します。
            }
        },
    );

    socket.addEventListener(
        'close',
        () => {
            options.onClose?.();
        },
    );

    socket.addEventListener(
        'error',
        () => {
            options.onError?.();
        },
    );

    return socket;
}