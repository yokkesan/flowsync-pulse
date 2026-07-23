import type {
    TaskPriority,
    TaskStatus,
} from '../types/task';

export function getTaskStatusLabel(
    status: TaskStatus,
): string {
    switch (status) {
        case 'not_started':
            return '未着手';

        case 'in_progress':
            return '進行中';

        case 'completed':
            return '完了';

        case 'suspended':
            return '保留';
    }
}

export function getTaskPriorityLabel(
    priority: TaskPriority,
): string {
    switch (priority) {
        case 'high':
            return '高';

        case 'medium':
            return '中';

        case 'low':
            return '低';
    }
}

export function formatTaskDate(
    value: string | null,
): string {
    if (!value) {
        return '未設定';
    }

    const date = new Date(
        `${value}T00:00:00`,
    );

    if (Number.isNaN(date.getTime())) {
        return value;
    }

    return new Intl.DateTimeFormat(
        'ja-JP',
        {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
        },
    ).format(date);
}

export function getAssigneeInitial(
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