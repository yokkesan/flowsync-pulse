import {
    useMemo,
    useState,
} from 'react';
import {
    Navigate,
    useNavigate,
    useParams,
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
import { useProject } from '../../hooks/useProject';
import { useTasks } from '../../hooks/useTasks';

const TASKS_PER_PAGE = 10;

function parseProjectId(
    value: string | undefined,
): number | null {
    if (!value) {
        return null;
    }

    const projectId = Number(value);

    if (
        !Number.isSafeInteger(projectId) ||
        projectId <= 0
    ) {
        return null;
    }

    return projectId;
}

export function TaskListPage() {
    const navigate = useNavigate();

    const { projectId: projectIdParam } =
        useParams();

    const { status, user } = useAuth();

    const projectId =
        parseProjectId(projectIdParam);

    const {
        project,
        isLoading: isProjectLoading,
        errorMessage: projectErrorMessage,
    } = useProject(projectId ?? 0);

    const {
        tasks,
        isLoading: isTasksLoading,
        errorMessage: tasksErrorMessage,
    } = useTasks(projectId ?? 0);

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
        priorityFilter,
        setPriorityFilter,
    ] = useState<TaskPriorityFilter>('all');

    const [
        currentPage,
        setCurrentPage,
    ] = useState(1);

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
                task.branch_name
                    .toLowerCase()
                    .includes(
                        normalizedKeyword,
                    );

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
                matchesStatus &&
                matchesAssignee &&
                matchesPriority
            );
        });
    }, [
        assigneeFilter,
        priorityFilter,
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

    function handleSelectTask(
        taskId: number,
    ): void {
        if (!projectId) {
            return;
        }

        navigate(
            `/projects/${projectId}/tasks/${taskId}`,
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

    if (!projectId) {
        return (
            <Navigate
                to="/projects"
                replace
            />
        );
    }

    const isLoading =
        isProjectLoading ||
        isTasksLoading;

    const errorMessage =
        projectErrorMessage ||
        tasksErrorMessage;

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
                    <div className="task-list-page__back">
                        <button
                            className="task-list-page__back-button"
                            type="button"
                            onClick={() => {
                                navigate(
                                    `/projects/${projectId}`,
                                );
                            }}
                        >
                            ← プロジェクト詳細へ
                        </button>
                    </div>

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
                    ) : !project ? (
                        <section className="task-list-page__state">
                            プロジェクトが見つかりません。
                        </section>
                    ) : (
                        <>
                            <header className="task-list-page__heading">
                                <div className="task-list-page__heading-text">
                                    <p className="task-list-page__eyebrow">
                                        プロジェクト
                                    </p>

                                    <h1 className="task-list-page__title">
                                        {project.name}
                                    </h1>
                                </div>

                                <button
                                    className="task-list-page__create-button"
                                    type="button"
                                    onClick={() => {
                                        navigate(
                                            `/projects/${projectId}/tasks/new`,
                                        );
                                    }}
                                >
                                    ＋ タスクを追加
                                </button>
                            </header>

                            <nav
                                className="task-list-page__tabs"
                                aria-label="プロジェクト詳細メニュー"
                            >
                                <button
                                    className="task-list-page__tab"
                                    type="button"
                                    onClick={() => {
                                        navigate(
                                            `/projects/${projectId}`,
                                        );
                                    }}
                                >
                                    概要
                                </button>

                                <button
                                    className="task-list-page__tab task-list-page__tab--active"
                                    type="button"
                                    aria-current="page"
                                >
                                    タスク
                                </button>
                            </nav>

                            <TaskFilters
                                searchKeyword={
                                    searchKeyword
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
                                members={
                                    project.members
                                }
                                onSearchKeywordChange={
                                    handleSearchKeywordChange
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