export type ProjectStatusView = {
    label: string;
    modifier: string;
};

export function getProjectStatusView(
    status: string,
): ProjectStatusView {
    switch (status) {
        case 'active':
            return {
                label: '進行中',
                modifier: 'active',
            };

        case 'completed':
            return {
                label: '完了',
                modifier: 'completed',
            };

        case 'paused':
        case 'stopped':
            return {
                label: '停止中',
                modifier: 'paused',
            };

        default:
            return {
                label: status || '未設定',
                modifier: 'unknown',
            };
    }
}

export function formatProjectDate(
    date: string | null,
): string {
    if (!date) {
        return '未設定';
    }

    const parsedDate = new Date(date);

    if (Number.isNaN(parsedDate.getTime())) {
        return date;
    }

    return new Intl.DateTimeFormat('ja-JP', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
    }).format(parsedDate);
}

export function getMemberInitial(
    displayName: string,
): string {
    const normalizedDisplayName =
        displayName.trim();

    if (!normalizedDisplayName) {
        return '?';
    }

    return Array.from(
        normalizedDisplayName,
    )[0];
}