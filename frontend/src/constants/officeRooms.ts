import mainOfficeImage from '../assets/offices/main-office.png';
import meetingRoomImage from '../assets/offices/meeting.png';
import officerRoomImage from '../assets/offices/officer-room.png';

import type { OfficeRoom } from '../types/office';

export const officeRooms: OfficeRoom[] = [
    {
        id: 'main-office',
        name: 'メインオフィス',
        image: mainOfficeImage,
    },
    {
        id: 'meeting-room',
        name: '会議室',
        image: meetingRoomImage,
    },
    {
        id: 'officer-room',
        name: '役員室',
        image: officerRoomImage,
    },
];