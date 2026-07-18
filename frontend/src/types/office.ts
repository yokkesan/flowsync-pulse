export type OfficeRoomId =
    | 'main-office'
    | 'meeting-room'
    | 'officer-room';

export type OfficeRoom = {
    id: OfficeRoomId;
    name: string;
    image: string;
};