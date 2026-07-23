import type { TaskStatus } from '../../types/task';
import { getTaskStatusLabel } from '../../utils/taskFormatters';

type TaskStatusBadgeProps = {
    status: TaskStatus;
};

export function TaskStatusBadge({
    status,
}: TaskStatusBadgeProps) {
    return (
        <span
            className={`task-status-badge task-status-badge--${status}`}
        >
            {getTaskStatusLabel(status)}
        </span>
    );
}