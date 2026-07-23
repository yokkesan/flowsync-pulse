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
    TaskForm,
    type TaskFormValues,
} from '../../components/tasks/TaskForm';
import { useAuth } from '../../contexts/AuthContext';
import { useProject } from '../../hooks/useProject';
import { useTask } from '../../hooks/useTask';
import { ApiError } from '../../services/authApi';
import { getAccessToken } from '../../services/authStorage';
import { updateTask } from '../../services/taskApi';
import type {
    TaskWriteRequest,
} from '../../types/task';

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

export function TaskEditPage() {
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
        project,
        isLoading: isProjectLoading,
        errorMessage: projectErrorMessage,
    } = useProject(projectId ?? 0);

    const {
        task,
        isLoading: isTaskLoading,
        errorMessage: taskErrorMessage,
    } = useTask(
        projectId ?? 0,
        taskId ?? 0,
    );

    const [isSubmitting, setIsSubmitting] =
        useState(false);

    const [
        submitErrorMessage,
        setSubmitErrorMessage,
    ] = useState('');

    const initialValues =
        useMemo<TaskFormValues | undefined>(
            () => {
                if (!task) {
                    return undefined;
                }

                return {
                    name: task.name,
                    description:
                        task.description ?? '',
                    assigneeUserId: String(
                        task.assignee_user_id,
                    ),
                    branchName:
                        task.branch_name,
                    status: task.status,
                    priority: task.priority,
                    startDate:
                        task.start_date ?? '',
                    dueDate:
                        task.due_date ?? '',
                };
            },
            [task],
        );

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

    async function handleSubmit(
        request: TaskWriteRequest,
    ): Promise<void> {
        if (!projectId || !taskId) {
            setSubmitErrorMessage(
                'プロジェクトIDまたはタスクIDが正しくありません。',
            );
            return;
        }

        const accessToken =
            getAccessToken();

        if (!accessToken) {
            setSubmitErrorMessage(
                'ログイン情報が見つかりません。',
            );
            return;
        }

        setIsSubmitting(true);
        setSubmitErrorMessage('');

        try {
            const updatedTask =
                await updateTask(
                    accessToken,
                    projectId,
                    taskId,
                    request,
                );

            navigate(
                `/projects/${projectId}/tasks/${updatedTask.task_id}`,
                {
                    replace: true,
                },
            );
        } catch (error) {
            if (error instanceof ApiError) {
                setSubmitErrorMessage(
                    error.message,
                );
                return;
            }

            setSubmitErrorMessage(
                'タスク編集中に予期しないエラーが発生しました。',
            );
        } finally {
            setIsSubmitting(false);
        }
    }

    const isLoading =
        isProjectLoading ||
        isTaskLoading;

    const errorMessage =
        projectErrorMessage ||
        taskErrorMessage;

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
                    title="タスク編集"
                    displayName={
                        user.display_name
                    }
                />

                <main className="task-list-page__main">
                    <section className="task-list-page__heading">
                        <div>
                            <h1 className="task-list-page__title">
                                タスクの編集
                            </h1>

                            <p className="task-list-page__description">
                                タスク情報を変更します。
                            </p>
                        </div>
                    </section>

                    {isLoading ? (
                        <section
                            className="task-list-page__state"
                            role="status"
                        >
                            タスク情報を読み込んでいます。
                        </section>
                    ) : errorMessage ? (
                        <section
                            className="task-list-page__state task-list-page__state--error"
                            role="alert"
                        >
                            {errorMessage}
                        </section>
                    ) : !project ||
                      !initialValues ? (
                        <section className="task-list-page__state">
                            タスクが見つかりません。
                        </section>
                    ) : (
                        <TaskForm
                            initialValues={
                                initialValues
                            }
                            members={
                                project.members
                            }
                            isSubmitting={
                                isSubmitting
                            }
                            submitErrorMessage={
                                submitErrorMessage
                            }
                            submitLabel="変更を保存"
                            onSubmit={
                                handleSubmit
                            }
                            onCancel={() => {
                                navigate(
                                    `/projects/${projectId}/tasks/${taskId}`,
                                );
                            }}
                        />
                    )}
                </main>
            </div>
        </div>
    );
}