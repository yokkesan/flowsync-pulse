import type {
    KeyboardEvent,
} from 'react';

import type { Task } from '../../types/task';
import {
    formatTaskDate,
    getAssigneeInitial,
} from '../../utils/taskFormatters';

import { TaskPriorityBadge } from './TaskPriorityBadge';
import { TaskStatusBadge } from './TaskStatusBadge';

type TaskTableProps = {
    tasks: Task[];
    onSelectTask: (
        taskId: number,
    ) => void;
};

export function TaskTable({
    tasks,
    onSelectTask,
}: TaskTableProps) {
    function handleRowKeyDown(
        event: KeyboardEvent<HTMLTableRowElement>,
        taskId: number,
    ): void {
        if (
            event.key !== 'Enter' &&
            event.key !== ' '
        ) {
            return;
        }

        event.preventDefault();
        onSelectTask(taskId);
    }

    return (
        <div className="task-table">
            <div className="task-table__scroll">
                <table className="task-table__table">
                    <thead>
                        <tr>
                            <th scope="col">
                                タスク
                            </th>

                            <th scope="col">
                                ステータス
                            </th>

                            <th scope="col">
                                担当者
                            </th>

                            <th scope="col">
                                優先度
                            </th>

                            <th scope="col">
                                ブランチ
                            </th>

                            <th scope="col">
                                期限
                            </th>
                        </tr>
                    </thead>

                    <tbody>
                        {tasks.length === 0 ? (
                            <tr>
                                <td
                                    className="task-table__empty"
                                    colSpan={6}
                                >
                                    条件に一致するタスクはありません。
                                </td>
                            </tr>
                        ) : (
                            tasks.map((task) => (
                                <tr
                                    className="task-table__row"
                                    key={task.task_id}
                                    tabIndex={0}
                                    onClick={() => {
                                        onSelectTask(
                                            task.task_id,
                                        );
                                    }}
                                    onKeyDown={(event) => {
                                        handleRowKeyDown(
                                            event,
                                            task.task_id,
                                        );
                                    }}
                                >
                                    <td>
                                        <div className="task-table__task">
                                            <strong className="task-table__task-name">
                                                {
                                                    task.name
                                                }
                                            </strong>

                                            {task.description && (
                                                <span className="task-table__task-description">
                                                    {
                                                        task.description
                                                    }
                                                </span>
                                            )}
                                        </div>
                                    </td>

                                    <td>
                                        <TaskStatusBadge
                                            status={
                                                task.status
                                            }
                                        />
                                    </td>

                                    <td>
                                        <div className="task-table__assignee">
                                            <span
                                                className="task-table__avatar"
                                                aria-hidden="true"
                                            >
                                                {getAssigneeInitial(
                                                    task.assignee_name,
                                                )}
                                            </span>

                                            <span className="task-table__assignee-name">
                                                {
                                                    task.assignee_name
                                                }
                                            </span>
                                        </div>
                                    </td>

                                    <td>
                                        <TaskPriorityBadge
                                            priority={
                                                task.priority
                                            }
                                        />
                                    </td>

                                    <td>
                                        <code className="task-table__branch">
                                            {
                                                task.branch_name
                                            }
                                        </code>
                                    </td>

                                    <td>
                                        {formatTaskDate(
                                            task.due_date,
                                        )}
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}