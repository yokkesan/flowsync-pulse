import { useState } from 'react';

import { OfficeControls } from '../../components/office/OfficeControls';
import { OfficeSwitcher } from '../../components/office/OfficeSwitcher';
import { officeRooms } from '../../constants/officeRooms';

import type { OfficeRoomId } from '../../types/office';

export function VirtualOfficePage() {
    const [selectedRoomId, setSelectedRoomId] =
        useState<OfficeRoomId>('main-office');

    const selectedRoom = officeRooms.find(
        (room) => room.id === selectedRoomId,
    );

    if (!selectedRoom) {
        throw new Error('選択されたオフィスが見つかりません。');
    }

    return (
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

                <div className="virtual-office-page__avatar-layer">
                    {/* 後ほどアバターを配置 */}
                </div>
            </section>

            <OfficeControls />
        </main>
    );
}