import {
    useRef,
    useState,
    type PointerEvent,
} from 'react';

type OfficeAvatarProps = {
    displayName: string;
};

type AvatarPosition = {
    x: number;
    y: number;
};

type DragState = {
    pointerId: number;
    offsetX: number;
    offsetY: number;
};

function getAvatarInitial(displayName: string): string {
    const normalizedDisplayName = displayName.trim();

    if (!normalizedDisplayName) {
        return '?';
    }

    return Array.from(normalizedDisplayName)[0];
}

export function OfficeAvatar({
    displayName,
}: OfficeAvatarProps) {
    const avatarLayerRef =
        useRef<HTMLDivElement | null>(null);

    const [avatarPosition, setAvatarPosition] =
        useState<AvatarPosition>({
            x: 50,
            y: 58,
        });

    const [dragState, setDragState] =
        useState<DragState | null>(null);

    const avatarInitial =
        getAvatarInitial(displayName);

    const handlePointerDown = (
        event: PointerEvent<HTMLDivElement>,
    ) => {
        const avatarLayer = avatarLayerRef.current;

        if (!avatarLayer) {
            return;
        }

        const layerRect =
            avatarLayer.getBoundingClientRect();

        const avatarX =
            layerRect.left +
            (avatarPosition.x / 100) *
                layerRect.width;

        const avatarY =
            layerRect.top +
            (avatarPosition.y / 100) *
                layerRect.height;

        event.currentTarget.setPointerCapture(
            event.pointerId,
        );

        setDragState({
            pointerId: event.pointerId,
            offsetX: event.clientX - avatarX,
            offsetY: event.clientY - avatarY,
        });
    };

    const handlePointerMove = (
        event: PointerEvent<HTMLDivElement>,
    ) => {
        if (
            !dragState ||
            dragState.pointerId !== event.pointerId
        ) {
            return;
        }

        const avatarLayer = avatarLayerRef.current;

        if (!avatarLayer) {
            return;
        }

        const layerRect =
            avatarLayer.getBoundingClientRect();

        const x =
            ((event.clientX -
                dragState.offsetX -
                layerRect.left) /
                layerRect.width) *
            100;

        const y =
            ((event.clientY -
                dragState.offsetY -
                layerRect.top) /
                layerRect.height) *
            100;

        setAvatarPosition({
            x: Math.min(96, Math.max(4, x)),
            y: Math.min(92, Math.max(8, y)),
        });
    };

    const handlePointerUp = (
        event: PointerEvent<HTMLDivElement>,
    ) => {
        if (
            !dragState ||
            dragState.pointerId !== event.pointerId
        ) {
            return;
        }

        if (
            event.currentTarget.hasPointerCapture(
                event.pointerId,
            )
        ) {
            event.currentTarget.releasePointerCapture(
                event.pointerId,
            );
        }

        setDragState(null);
    };

    return (
        <div
            ref={avatarLayerRef}
            className="virtual-office-page__avatar-layer"
        >
            <div
                className={[
                    'virtual-office-avatar',
                    dragState
                        ? 'virtual-office-avatar--dragging'
                        : '',
                ]
                    .filter(Boolean)
                    .join(' ')}
                style={{
                    left: `${avatarPosition.x}%`,
                    top: `${avatarPosition.y}%`,
                }}
                onPointerDown={handlePointerDown}
                onPointerMove={handlePointerMove}
                onPointerUp={handlePointerUp}
                onPointerCancel={handlePointerUp}
            >
                <div
                    className="virtual-office-avatar__icon"
                    aria-hidden="true"
                >
                    {avatarInitial}
                </div>

                <p className="virtual-office-avatar__name">
                    {displayName}
                </p>
            </div>
        </div>
    );
}