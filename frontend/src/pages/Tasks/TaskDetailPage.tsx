import {
    useState,
} from 'react';
import {
    Navigate,
    useNavigate,
    useParams,
} from 'react-router-dom';

import { AppHeader } from '../../components/layout/AppHeader';
import { AppSidebar } from '../../components/layout/AppSidebar';
import { TaskPriorityBadge } from '../../components/tasks/TaskPriorityBadge';
import { TaskStatusBadge } from '../../components/tasks/TaskStatusBadge';
import { useAuth } from '../../contexts/AuthContext';
import { useTask } from '../../hooks/useTask';
import { ApiError } from '../../services/authApi';
import { getAccessToken } from '../../services/authStorage';
import { deleteTask } from '../../services/taskApi';
import { formatTaskDate } from '../../utils/taskFormatters';

function parseId(
    value: string | undefined,
): number | null {
    if (!value) {
        return null;
    }

    const id = Number(value);

    if (
        !Number.isSafeInteger(id) ||
        id <= 0
    ) {
        return null;
    }

    return id;
}

function formatDateTime(
    value: string | null,
): string {
    if (!value) {
        return '未設定';
    }

    const date = new Date(value);

    if (Number.isNaN(date.getTime())) {
        return value;
    }

    return new Intl.DateTimeFormat(
        'ja-JP',
        {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
        },
    ).format(date);
}

export function TaskDetailPage() {
    const navigate = useNavigate();

    const {
        projectId: projectIdParam,
        taskId: taskIdParam,
    } = useParams();

    const { status, user } = useAuth();

    const projectId =
        parseId(projectIdParam);

    const taskId =
        parseId(taskIdParam);

    const {
        task,
        isLoading,
        errorMessage,
    } = useTask(
        projectId ?? 0,
        taskId ?? 0,
    );

    const [isDeleting, setIsDeleting] =
        useState(false);

    const [
        deleteErrorMessage,
        setDeleteErrorMessage,
    ] = useState('');

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

    if (!taskId) {
        return (
            <Navigate
                to={`/projects/${projectId}/tasks`}
                replace
            />
        );
    }

    async function handleDelete():
    Promise<void> {
        if (!task) {
            return;
        }

        const shouldDelete = window.confirm(
            `「${task.name}」を削除しますか？\nこの操作は取り消せません。`,
        );

        if (!shouldDelete) {
            return;
        }

        const accessToken =
            getAccessToken();

        if (!accessToken) {
            setDeleteErrorMessage(
                'ログイン情報が見つかりません。',
            );
            return;
        }

        setIsDeleting(true);
        setDeleteErrorMessage('');

        try {
            if (!projectId) {
                setDeleteErrorMessage(
                    'プロジェクトIDが正しくありません。',
                );
                return;
            }

            await deleteTask(
                accessToken,
                projectId,
                task.task_id,
            );

            navigate(
                `/projects/${projectId}/tasks`,
                {
                    replace: true,
                },
            );
        } catch (error) {
            if (error instanceof ApiError) {
                setDeleteErrorMessage(
                    error.message,
                );
                return;
            }

            setDeleteErrorMessage(
                'タスク削除中に予期しないエラーが発生しました。',
            );
        } finally {
            setIsDeleting(false);
        }
    }

    return (
        <div className="task-detail-page">
            <AppSidebar
                displayName={
                    user.display_name
                }
                role={user.role}
            />

            <div className="task-detail-page__body">
                <AppHeader
                    title="タスク詳細"
                    displayName={
                        user.display_name
                    }
                />

                <main className="task-detail-page__main">
                    <div className="task-detail-page__back">
                        <button
                            className="task-detail-page__back-button"
                            type="button"
                            onClick={() => {
                                navigate(
                                    `/projects/${projectId}/tasks`,
                                );
                            }}
                        >
                            ← タスク一覧へ
                        </button>
                    </div>

                    {isLoading ? (
                        <section
                            className="task-detail-page__state"
                            role="status"
                        >
                            タスク詳細を読み込んでいます。
                        </section>
                    ) : errorMessage ? (
                        <section
                            className="task-detail-page__state task-detail-page__state--error"
                            role="alert"
                        >
                            {errorMessage}
                        </section>
                    ) : !task ? (
                        <section className="task-detail-page__state">
                            タスクが見つかりません。
                        </section>
                    ) : (
                        <>
                            <header className="task-detail-page__heading">
                                <div className="task-detail-page__heading-text">
                                    <div className="task-detail-page__title-row">
                                        <h1 className="task-detail-page__title">
                                            {task.name}
                                        </h1>

                                        <TaskStatusBadge
                                            status={
                                                task.status
                                            }
                                        />

                                        <TaskPriorityBadge
                                            priority={
                                                task.priority
                                            }
                                        />
                                    </div>
                                </div>

                                <div className="task-detail-page__actions">
                                    <button
                                        className="task-detail-page__button task-detail-page__button--secondary"
                                        type="button"
                                        disabled={
                                            isDeleting
                                        }
                                        onClick={() => {
                                            navigate(
                                                `/projects/${projectId}/tasks/${task.task_id}/edit`,
                                            );
                                        }}
                                    >
                                        編集
                                    </button>

                                    <button
                                        className="task-detail-page__button task-detail-page__button--danger"
                                        type="button"
                                        disabled={
                                            isDeleting
                                        }
                                        onClick={() => {
                                            void handleDelete();
                                        }}
                                    >
                                        {isDeleting
                                            ? '削除中'
                                            : '削除'}
                                    </button>
                                </div>
                            </header>

                            {deleteErrorMessage && (
                                <div
                                    className="task-detail-page__delete-error"
                                    role="alert"
                                >
                                    {deleteErrorMessage}
                                </div>
                            )}

                            <section className="task-detail-page__content">
                                <div className="task-detail-page__card">
                                    <h2 className="task-detail-page__card-title">
                                        基本情報
                                    </h2>

                                    <dl className="task-detail-page__definition-list">
                                        <div className="task-detail-page__definition-item">
                                            <dt>
                                                説明
                                            </dt>

                                            <dd>
                                                {task.description ||
                                                    '未設定'}
                                            </dd>
                                        </div>

                                        <div className="task-detail-page__definition-item">
                                            <dt>
                                                担当者
                                            </dt>

                                            <dd>
                                                {
                                                    task.assignee_name
                                                }
                                            </dd>
                                        </div>

                                        <div className="task-detail-page__definition-item">
                                            <dt>
                                                開始日
                                            </dt>

                                            <dd>
                                                {formatTaskDate(
                                                    task.start_date,
                                                )}
                                            </dd>
                                        </div>

                                        <div className="task-detail-page__definition-item">
                                            <dt>
                                                期限
                                            </dt>

                                            <dd>
                                                {formatTaskDate(
                                                    task.due_date,
                                                )}
                                            </dd>
                                        </div>
                                    </dl>
                                </div>

                                <div className="task-detail-page__card">
                                    <h2 className="task-detail-page__card-title">
                                        管理情報
                                    </h2>

                                    <dl className="task-detail-page__definition-list">
                                        <div className="task-detail-page__definition-item">
                                            <dt>
                                                完了日時
                                            </dt>

                                            <dd>
                                                {formatDateTime(
                                                    task.completed_at,
                                                )}
                                            </dd>
                                        </div>

                                        <div className="task-detail-page__definition-item">
                                            <dt>
                                                登録日時
                                            </dt>

                                            <dd>
                                                {formatDateTime(
                                                    task.created_at,
                                                )}
                                            </dd>
                                        </div>

                                        <div className="task-detail-page__definition-item">
                                            <dt>
                                                更新日時
                                            </dt>

                                            <dd>
                                                {formatDateTime(
                                                    task.updated_at,
                                                )}
                                            </dd>
                                        </div>
                                    </dl>
                                </div>
                            </section>
                        </>
                    )}
                </main>
            </div>
        </div>
    );
}