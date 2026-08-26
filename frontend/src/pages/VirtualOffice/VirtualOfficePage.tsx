import {
    useEffect,
    useState,
} from 'react';
import { Navigate } from 'react-router-dom';

import { AppSidebar } from '../../components/layout/AppSidebar';
import { OfficeAvatar } from '../../components/office/OfficeAvatar';
import { OfficeControls } from '../../components/office/OfficeControls';
import { OfficeSwitcher } from '../../components/office/OfficeSwitcher';
import { officeRooms } from '../../constants/officeRooms';
import { useAuth } from '../../contexts/AuthContext';
import { getAccessToken } from '../../services/authStorage';
import {
    getCurrentWorkContext,
} from '../../services/realtimeApi';
import { connectRealtimeSocket } from '../../services/realtimeSocket';

import type { OfficeRoomId } from '../../types/office';
import type { WorkContextChangedPayload } from '../../types/realtime';

export function VirtualOfficePage() {
    const { status, user } = useAuth();

    const [selectedRoomId, setSelectedRoomId] =
        useState<OfficeRoomId>('main-office');

    const [workContext, setWorkContext] =
        useState<WorkContextChangedPayload | null>(
            null,
        );

    useEffect(() => {
        if (
            status !== 'authenticated' ||
            !user
        ) {
            return;
        }

        const accessToken = getAccessToken();

        if (!accessToken) {
            return;
        }

        let socket: WebSocket | null = null;
        let disposed = false;

        const connect = async () => {
            try {
                const currentWorkContext =
                    await getCurrentWorkContext(
                        accessToken,
                    );

                if (disposed) {
                    return;
                }

                if (
                    currentWorkContext &&
                    currentWorkContext.user_id ===
                        user.id
                ) {
                    setWorkContext(
                        currentWorkContext,
                    );
                }

                const connectedSocket =
                    await connectRealtimeSocket({
                        accessToken,
                        onEvent: (event) => {
                            if (
                                event.type !==
                                'work_context.changed'
                            ) {
                                return;
                            }

                            if (
                                event.payload.user_id !==
                                user.id
                            ) {
                                return;
                            }

                            setWorkContext(
                                event.payload,
                            );
                        },
                    });

                if (disposed) {
                    connectedSocket.close();
                    return;
                }

                socket = connectedSocket;
            } catch {
                // 接続失敗時の表示・再接続処理は
                // 後続で実装します。
            }
        };

        void connect();

        return () => {
            disposed = true;
            socket?.close();
        };
    }, [
        status,
        user,
    ]);

    if (status === 'checking') {
        return (
            <div role="status">
                ログイン状態を確認しています。
            </div>
        );
    }

    if (
        status !== 'authenticated' ||
        !user
    ) {
        return (
            <Navigate
                to="/login"
                replace
            />
        );
    }

    const selectedRoom = officeRooms.find(
        (room) => room.id === selectedRoomId,
    );

    if (!selectedRoom) {
        throw new Error(
            '選択されたオフィスが見つかりません。',
        );
    }

    return (
        <div className="virtual-office-layout">
            <AppSidebar
                displayName={user.display_name}
                role={user.role}
            />

            <main className="virtual-office-page">
                <header className="virtual-office-page__header">
                    <div className="virtual-office-page__heading">
                        <p className="virtual-office-page__service-name">
                            FlowSync Pulse
                        </p>

                        <h1 className="virtual-office-page__title">
                            {selectedRoom.name}
                        </h1>
                    </div>

                    <div className="virtual-office-page__connection">
                        <span
                            className="virtual-office-page__connection-indicator"
                            aria-hidden="true"
                        />

                        <span>接続中</span>
                    </div>
                </header>

                <OfficeSwitcher
                    rooms={officeRooms}
                    selectedRoomId={selectedRoomId}
                    onRoomChange={setSelectedRoomId}
                />

                <section
                    className={[
                        'virtual-office-page__scene',
                        `virtual-office-page__scene--${selectedRoom.id}`,
                    ].join(' ')}
                    aria-label={`${selectedRoom.name}のバーチャルオフィス`}
                >
                    <img
                        className="virtual-office-page__background"
                        src={selectedRoom.image}
                        alt=""
                        draggable={false}
                    />

                    <OfficeAvatar
                        displayName={user.display_name}
                        workContext={workContext}
                    />
                </section>

                <OfficeControls />
            </main>
        </div>
    );
}