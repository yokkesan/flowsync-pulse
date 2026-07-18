import type { OfficeRoom, OfficeRoomId } from '../../types/office';

type OfficeSwitcherProps = {
    rooms: OfficeRoom[];
    selectedRoomId: OfficeRoomId;
    onRoomChange: (roomId: OfficeRoomId) => void;
};

export function OfficeSwitcher({
    rooms,
    selectedRoomId,
    onRoomChange,
}: OfficeSwitcherProps) {
    return (
        <nav
            className="office-switcher"
            aria-label="オフィス切り替え"
        >
            {rooms.map((room) => {
                const isActive = room.id === selectedRoomId;

                return (
                    <button
                        key={room.id}
                        type="button"
                        className={
                            isActive
                                ? 'office-switcher__button office-switcher__button--active'
                                : 'office-switcher__button'
                        }
                        aria-pressed={isActive}
                        onClick={() => onRoomChange(room.id)}
                    >
                        {room.name}
                    </button>
                );
            })}
        </nav>
    );
}