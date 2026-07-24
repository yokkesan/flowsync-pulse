import {
    useMemo,
    useState,
} from 'react';
import {
    Navigate,
    useNavigate,
} from 'react-router-dom';

import { AppHeader } from '../../components/layout/AppHeader';
import { AppSidebar } from '../../components/layout/AppSidebar';
import {
    TaskFilters,
    type TaskPriorityFilter,
    type TaskStatusFilter,
} from '../../components/tasks/TaskFilters';
import { TaskPagination } from '../../components/tasks/TaskPagination';
import { TaskTable } from '../../components/tasks/TaskTable';
import { useAuth } from '../../contexts/AuthContext';
import { useAllTasks } from '../../hooks/useAllTasks';
import type { ProjectMember } from '../../types/project';

const TASKS_PER_PAGE = 10;

export function AllTaskListPage() {
    const navigate = useNavigate();

    const { status, user } = useAuth();

    const {
        tasks,
        isLoading,
        errorMessage,
    } = useAllTasks();

    const [
        searchKeyword,
        setSearchKeyword,
    ] = useState('');

    const [
        statusFilter,
        setStatusFilter,
    ] = useState<TaskStatusFilter>('all');

    const [
        assigneeFilter,
        setAssigneeFilter,
    ] = useState('all');

    const [
        projectFilter,
        setProjectFilter,
    ] = useState('all');

    const [
        priorityFilter,
        setPriorityFilter,
    ] = useState<TaskPriorityFilter>('all');

    const [
        currentPage,
        setCurrentPage,
    ] = useState(1);

    const assignees = useMemo<ProjectMember[]>(
        () => {
            const assigneeMap = new Map<
                number,
                ProjectMember
            >();

            tasks.forEach((task) => {
                if (
                    assigneeMap.has(
                        task.assignee_user_id,
                    )
                ) {
                    return;
                }

                assigneeMap.set(
                    task.assignee_user_id,
                    {
                        user_id:
                            task.assignee_user_id,
                        display_name:
                            task.assignee_name,
                        role: '',
                        status: 'active',
                    },
                );
            });

            return Array.from(
                assigneeMap.values(),
            ).sort((first, second) =>
                first.display_name.localeCompare(
                    second.display_name,
                    'ja',
                ),
            );
        },
        [tasks],
    );

    const projects = useMemo(
        () => {
            const projectMap = new Map<
                number,
                {
                    project_id: number;
                    project_name: string;
                }
            >();

            tasks.forEach((task) => {
                if (
                    projectMap.has(
                        task.project_id,
                    )
                ) {
                    return;
                }

                projectMap.set(
                    task.project_id,
                    {
                        project_id:
                            task.project_id,
                        project_name:
                            task.project_name,
                    },
                );
            });

            return Array.from(
                projectMap.values(),
            ).sort((first, second) =>
                first.project_name.localeCompare(
                    second.project_name,
                    'ja',
                ),
            );
        },
        [tasks],
    );

    const filteredTasks = useMemo(() => {
        const normalizedKeyword =
            searchKeyword
                .trim()
                .toLowerCase();

        return tasks.filter((task) => {
            const matchesKeyword =
                normalizedKeyword === '' ||
                task.name
                    .toLowerCase()
                    .includes(
                        normalizedKeyword,
                    ) ||
                task.project_name
                    .toLowerCase()
                    .includes(
                        normalizedKeyword,
                    ) ||
                task.branch_name
                    .toLowerCase()
                    .includes(
                        normalizedKeyword,
                    );

            const matchesProject =
                projectFilter === 'all' ||
                task.project_id ===
                    Number(projectFilter);

            const matchesStatus =
                statusFilter === 'all' ||
                task.status ===
                    statusFilter;

            const matchesAssignee =
                assigneeFilter === 'all' ||
                task.assignee_user_id ===
                    Number(assigneeFilter);

            const matchesPriority =
                priorityFilter === 'all' ||
                task.priority ===
                    priorityFilter;

            return (
                matchesKeyword &&
                matchesProject &&
                matchesStatus &&
                matchesAssignee &&
                matchesPriority
            );
        });
    }, [
        assigneeFilter,
        priorityFilter,
        projectFilter,
        searchKeyword,
        statusFilter,
        tasks,
    ]);

    const totalPages = Math.max(
        1,
        Math.ceil(
            filteredTasks.length /
                TASKS_PER_PAGE,
        ),
    );

    const normalizedCurrentPage = Math.min(
        currentPage,
        totalPages,
    );

    const displayedTasks = useMemo(() => {
        const startIndex =
            (normalizedCurrentPage - 1) *
            TASKS_PER_PAGE;

        return filteredTasks.slice(
            startIndex,
            startIndex +
                TASKS_PER_PAGE,
        );
    }, [
        filteredTasks,
        normalizedCurrentPage,
    ]);

    function resetCurrentPage(): void {
        setCurrentPage(1);
    }

    function handleSearchKeywordChange(
        value: string,
    ): void {
        setSearchKeyword(value);
        resetCurrentPage();
    }

    function handleProjectFilterChange(
        value: string,
    ): void {
        setProjectFilter(value);
        resetCurrentPage();
    }

    function handleStatusFilterChange(
        value: TaskStatusFilter,
    ): void {
        setStatusFilter(value);
        resetCurrentPage();
    }

    function handleAssigneeFilterChange(
        value: string,
    ): void {
        setAssigneeFilter(value);
        resetCurrentPage();
    }

    function handlePriorityFilterChange(
        value: TaskPriorityFilter,
    ): void {
        setPriorityFilter(value);
        resetCurrentPage();
    }

    function handleCreateTask(): void {
        navigate('/tasks/new');
    }

    function handleSelectTask(
        taskId: number,
    ): void {
        const selectedTask = tasks.find(
            (task) =>
                task.task_id === taskId,
        );

        if (!selectedTask) {
            return;
        }

        navigate(
            `/projects/${selectedTask.project_id}/tasks/${selectedTask.task_id}`,
        );
    }

    if (status === 'checking') {
        return (
            <div role="status">
                ログイン状態を確認しています。
            </div>
        );
    }

    if (
        status !== 'authenticated' ||
        !user
    ) {
        return (
            <Navigate
                to="/login"
                replace
            />
        );
    }

    return (
        <div className="task-list-page">
            <AppSidebar
                displayName={
                    user.display_name
                }
                role={user.role}
            />

            <div className="task-list-page__body">
                <AppHeader
                    title="タスク管理"
                    displayName={
                        user.display_name
                    }
                />

                <main className="task-list-page__main">
                    <header className="task-list-page__heading">
                        <div className="task-list-page__heading-text">
                            <h1 className="task-list-page__title">
                                全タスク一覧
                            </h1>

                            <p className="task-list-page__description">
                                会社内の全プロジェクトに登録されているタスクを確認できます。
                            </p>
                        </div>

                        <button
                            type="button"
                            className="task-list-page__create-button"
                            onClick={handleCreateTask}
                        >
                            新規登録
                        </button>
                    </header>

                    {isLoading ? (
                        <section
                            className="task-list-page__state"
                            role="status"
                        >
                            タスク一覧を読み込んでいます。
                        </section>
                    ) : errorMessage ? (
                        <section
                            className="task-list-page__state task-list-page__state--error"
                            role="alert"
                        >
                            {errorMessage}
                        </section>
                    ) : (
                        <>
                            <TaskFilters
                                searchKeyword={
                                    searchKeyword
                                }
                                projectFilter={
                                    projectFilter
                                }
                                statusFilter={
                                    statusFilter
                                }
                                assigneeFilter={
                                    assigneeFilter
                                }
                                priorityFilter={
                                    priorityFilter
                                }
                                projects={
                                    projects
                                }
                                members={
                                    assignees
                                }
                                onSearchKeywordChange={
                                    handleSearchKeywordChange
                                }
                                onProjectFilterChange={
                                    handleProjectFilterChange
                                }
                                onStatusFilterChange={
                                    handleStatusFilterChange
                                }
                                onAssigneeFilterChange={
                                    handleAssigneeFilterChange
                                }
                                onPriorityFilterChange={
                                    handlePriorityFilterChange
                                }
                            />

                            <section className="task-list-page__table-card">
                                <TaskTable
                                    tasks={
                                        displayedTasks
                                    }
                                    onSelectTask={
                                        handleSelectTask
                                    }
                                />

                                <TaskPagination
                                    currentPage={
                                        normalizedCurrentPage
                                    }
                                    totalPages={
                                        totalPages
                                    }
                                    totalItems={
                                        filteredTasks.length
                                    }
                                    itemsPerPage={
                                        TASKS_PER_PAGE
                                    }
                                    onPageChange={
                                        setCurrentPage
                                    }
                                />
                            </section>
                        </>
                    )}
                </main>
            </div>
        </div>
    );
}