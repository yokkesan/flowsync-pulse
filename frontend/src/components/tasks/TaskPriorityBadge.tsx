import type { TaskPriority } from '../../types/task';
import { getTaskPriorityLabel } from '../../utils/taskFormatters';

type TaskPriorityBadgeProps = {
    priority: TaskPriority;
};

export function TaskPriorityBadge({
    priority,
}: TaskPriorityBadgeProps) {
    return (
        <span
            className={`task-priority-badge task-priority-badge--${priority}`}
        >
            {getTaskPriorityLabel(priority)}
        </span>
    );
}