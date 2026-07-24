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
import { TaskForm } from '../../components/tasks/TaskForm';
import { useAuth } from '../../contexts/AuthContext';
import { useProject } from '../../hooks/useProject';
import { ApiError } from '../../services/authApi';
import { getAccessToken } from '../../services/authStorage';
import { createTask } from '../../services/taskApi';
import type {
    TaskWriteRequest,
} from '../../types/task';

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

export function TaskCreatePage() {
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

    const [isSubmitting, setIsSubmitting] =
        useState(false);

    const [
        submitErrorMessage,
        setSubmitErrorMessage,
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

    async function handleSubmit(
        request: TaskWriteRequest,
    ): Promise<void> {
        if (!projectId) {
            setSubmitErrorMessage(
                'プロジェクトIDが正しくありません。',
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
            const task = await createTask(
                accessToken,
                projectId,
                request,
            );

            navigate(
                `/projects/${projectId}/tasks/${task.task_id}`,
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
                'タスク登録中に予期しないエラーが発生しました。',
            );
        } finally {
            setIsSubmitting(false);
        }
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
                    title="タスク登録"
                    displayName={
                        user.display_name
                    }
                />

                <main className="task-list-page__main">
                    <section className="task-list-page__heading">
                        <div>
                            <h1 className="task-list-page__title">
                                タスクの登録
                            </h1>

                            <p className="task-list-page__description">
                                プロジェクトに新しいタスクを追加します。
                            </p>
                        </div>
                    </section>

                    {isProjectLoading ? (
                        <section
                            className="task-list-page__state"
                            role="status"
                        >
                            プロジェクト情報を読み込んでいます。
                        </section>
                    ) : projectErrorMessage ? (
                        <section
                            className="task-list-page__state task-list-page__state--error"
                            role="alert"
                        >
                            {projectErrorMessage}
                        </section>
                    ) : !project ? (
                        <section className="task-list-page__state">
                            プロジェクトが見つかりません。
                        </section>
                    ) : (
                        <TaskForm
                            members={project.members}
                            isSubmitting={
                                isSubmitting
                            }
                            submitErrorMessage={
                                submitErrorMessage
                            }
                            submitLabel="登録する"
                            onSubmit={handleSubmit}
                            onCancel={() => {
                                navigate(
                                    `/projects/${projectId}/tasks`,
                                );
                            }}
                        />
                    )}
                </main>
            </div>
        </div>
    );
}