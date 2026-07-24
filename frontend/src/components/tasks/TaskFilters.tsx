import type {
    ChangeEvent,
} from 'react';

import type { ProjectMember } from '../../types/project';
import type {
    TaskPriority,
    TaskStatus,
} from '../../types/task';

export type TaskStatusFilter =
    | 'all'
    | TaskStatus;

export type TaskPriorityFilter =
    | 'all'
    | TaskPriority;

export type TaskProjectFilterOption = {
    project_id: number;
    project_name: string;
};

type TaskFiltersProps = {
    searchKeyword: string;
    statusFilter: TaskStatusFilter;
    assigneeFilter: string;
    priorityFilter: TaskPriorityFilter;
    members: ProjectMember[];
    projectFilter?: string;
    projects?: TaskProjectFilterOption[];
    onSearchKeywordChange: (
        value: string,
    ) => void;
    onStatusFilterChange: (
        value: TaskStatusFilter,
    ) => void;
    onAssigneeFilterChange: (
        value: string,
    ) => void;
    onPriorityFilterChange: (
        value: TaskPriorityFilter,
    ) => void;
    onProjectFilterChange?: (
        value: string,
    ) => void;
};

export function TaskFilters({
    searchKeyword,
    statusFilter,
    assigneeFilter,
    priorityFilter,
    members,
    projectFilter,
    projects,
    onSearchKeywordChange,
    onStatusFilterChange,
    onAssigneeFilterChange,
    onPriorityFilterChange,
    onProjectFilterChange,
}: TaskFiltersProps) {
    function handleSearchKeywordChange(
        event: ChangeEvent<HTMLInputElement>,
    ): void {
        onSearchKeywordChange(
            event.target.value,
        );
    }

    function handleStatusFilterChange(
        event: ChangeEvent<HTMLSelectElement>,
    ): void {
        onStatusFilterChange(
            event.target.value as TaskStatusFilter,
        );
    }

    function handleAssigneeFilterChange(
        event: ChangeEvent<HTMLSelectElement>,
    ): void {
        onAssigneeFilterChange(
            event.target.value,
        );
    }

    function handlePriorityFilterChange(
        event: ChangeEvent<HTMLSelectElement>,
    ): void {
        onPriorityFilterChange(
            event.target.value as TaskPriorityFilter,
        );
    }

    function handleProjectFilterChange(
        event: ChangeEvent<HTMLSelectElement>,
    ): void {
        onProjectFilterChange?.(
            event.target.value,
        );
    }

    const shouldShowProjectFilter =
        projectFilter !== undefined &&
        projects !== undefined &&
        onProjectFilterChange !== undefined;

    return (
        <section
            className="task-filters"
            aria-label="タスクの絞り込み"
        >
            <div className="task-filters__search">
                <span
                    className="task-filters__search-icon"
                    aria-hidden="true"
                >
                    ⌕
                </span>

                <input
                    className="task-filters__search-input"
                    type="search"
                    value={searchKeyword}
                    placeholder="タスク名・ブランチ名で検索"
                    aria-label="タスク名またはブランチ名で検索"
                    onChange={
                        handleSearchKeywordChange
                    }
                />
            </div>

            {shouldShowProjectFilter && (
                <select
                    className="task-filters__select"
                    value={projectFilter}
                    aria-label="プロジェクトで絞り込む"
                    onChange={
                        handleProjectFilterChange
                    }
                >
                    <option value="all">
                        すべてのプロジェクト
                    </option>

                    {projects.map((project) => (
                        <option
                            key={
                                project.project_id
                            }
                            value={
                                project.project_id
                            }
                        >
                            {
                                project.project_name
                            }
                        </option>
                    ))}
                </select>
            )}

            <select
                className="task-filters__select"
                value={statusFilter}
                aria-label="ステータスで絞り込む"
                onChange={
                    handleStatusFilterChange
                }
            >
                <option value="all">
                    すべてのステータス
                </option>

                <option value="not_started">
                    未着手
                </option>

                <option value="in_progress">
                    進行中
                </option>

                <option value="completed">
                    完了
                </option>

                <option value="suspended">
                    保留
                </option>
            </select>

            <select
                className="task-filters__select"
                value={assigneeFilter}
                aria-label="担当者で絞り込む"
                onChange={
                    handleAssigneeFilterChange
                }
            >
                <option value="all">
                    すべての担当者
                </option>

                {members.map((member) => (
                    <option
                        key={member.user_id}
                        value={member.user_id}
                    >
                        {member.display_name}
                    </option>
                ))}
            </select>

            <select
                className="task-filters__select"
                value={priorityFilter}
                aria-label="優先度で絞り込む"
                onChange={
                    handlePriorityFilterChange
                }
            >
                <option value="all">
                    すべての優先度
                </option>

                <option value="high">
                    高
                </option>

                <option value="medium">
                    中
                </option>

                <option value="low">
                    低
                </option>
            </select>
        </section>
    );
}