import {
    useEffect,
    useRef,
    useState,
    type PointerEvent as ReactPointerEvent,
} from 'react';
import { useNavigate } from 'react-router-dom';

import { useAuth } from '../../contexts/AuthContext';

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
    startClientX: number;
    startClientY: number;
};

const DRAG_THRESHOLD_PX = 5;

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
    const navigate = useNavigate();
    const { logout } = useAuth();

    const avatarLayerRef =
        useRef<HTMLDivElement | null>(null);

    const avatarRef =
        useRef<HTMLDivElement | null>(null);

    const menuRef =
        useRef<HTMLDivElement | null>(null);

    const hasDraggedRef = useRef(false);

    const [avatarPosition, setAvatarPosition] =
        useState<AvatarPosition>({
            x: 50,
            y: 58,
        });

    const [dragState, setDragState] =
        useState<DragState | null>(null);

    const [isMenuOpen, setIsMenuOpen] =
        useState(false);

    const avatarInitial =
        getAvatarInitial(displayName);

    useEffect(() => {
        if (!isMenuOpen) {
            return;
        }

        const handleDocumentPointerDown = (
            event: PointerEvent,
        ) => {
            const target = event.target;

            if (!(target instanceof Node)) {
                return;
            }

            if (
                avatarRef.current?.contains(target) ||
                menuRef.current?.contains(target)
            ) {
                return;
            }

            setIsMenuOpen(false);
        };

        const handleDocumentKeyDown = (
            event: KeyboardEvent,
        ) => {
            if (event.key === 'Escape') {
                setIsMenuOpen(false);
            }
        };

        document.addEventListener(
            'pointerdown',
            handleDocumentPointerDown,
        );

        document.addEventListener(
            'keydown',
            handleDocumentKeyDown,
        );

        return () => {
            document.removeEventListener(
                'pointerdown',
                handleDocumentPointerDown,
            );

            document.removeEventListener(
                'keydown',
                handleDocumentKeyDown,
            );
        };
    }, [isMenuOpen]);

    const handlePointerDown = (
        event: ReactPointerEvent<HTMLDivElement>,
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

        hasDraggedRef.current = false;
        setIsMenuOpen(false);

        setDragState({
            pointerId: event.pointerId,
            offsetX: event.clientX - avatarX,
            offsetY: event.clientY - avatarY,
            startClientX: event.clientX,
            startClientY: event.clientY,
        });
    };

    const handlePointerMove = (
        event: ReactPointerEvent<HTMLDivElement>,
    ) => {
        if (
            !dragState ||
            dragState.pointerId !== event.pointerId
        ) {
            return;
        }

        const movedX = Math.abs(
            event.clientX - dragState.startClientX,
        );

        const movedY = Math.abs(
            event.clientY - dragState.startClientY,
        );

        if (
            movedX >= DRAG_THRESHOLD_PX ||
            movedY >= DRAG_THRESHOLD_PX
        ) {
            hasDraggedRef.current = true;
        }

        if (!hasDraggedRef.current) {
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
        event: ReactPointerEvent<HTMLDivElement>,
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

        const shouldOpenMenu =
            !hasDraggedRef.current;

        setDragState(null);

        if (shouldOpenMenu) {
            setIsMenuOpen((currentValue) => {
                return !currentValue;
            });
        }
    };

    const handleLogout = () => {
        logout();

        navigate('/login', {
            replace: true,
        });
    };

    return (
        <div
            ref={avatarLayerRef}
            className="virtual-office-page__avatar-layer"
        >
            <div
                ref={avatarRef}
                className={[
                    'virtual-office-avatar',
                    dragState
                        ? 'virtual-office-avatar--dragging'
                        : '',
                    isMenuOpen
                        ? 'virtual-office-avatar--menu-open'
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

                {isMenuOpen && (
                    <div
                        ref={menuRef}
                        className="virtual-office-avatar__menu"
                        role="menu"
                        aria-label={`${displayName}のユーザーメニュー`}
                        onPointerDown={(event) => {
                            event.stopPropagation();
                        }}
                    >
                        <p className="virtual-office-avatar__menu-name">
                            {displayName}
                        </p>

                        <button
                            className="virtual-office-avatar__menu-button virtual-office-avatar__menu-button--logout"
                            type="button"
                            role="menuitem"
                            onClick={handleLogout}
                        >
                            ログアウト
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
}