import {
    useState,
} from 'react';
import {
    Navigate,
    useNavigate,
} from 'react-router-dom';

import { AppHeader } from '../../components/layout/AppHeader';
import { AppSidebar } from '../../components/layout/AppSidebar';
import { ProjectForm } from '../../components/projects/ProjectForm';
import { useAuth } from '../../contexts/AuthContext';
import { useCompanyUsers } from '../../hooks/useCompanyUsers';
import { ApiError } from '../../services/authApi';
import { getAccessToken } from '../../services/authStorage';
import { createProject } from '../../services/projectApi';

import type {
    ProjectWriteRequest,
} from '../../types/project';

export function ProjectCreatePage() {
    const navigate = useNavigate();

    const { status, user } = useAuth();

    const {
        users,
        isLoading: isUsersLoading,
        errorMessage: usersErrorMessage,
    } = useCompanyUsers();

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

    async function handleSubmit(
        request: ProjectWriteRequest,
    ): Promise<void> {
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
            const project =
                await createProject(
                    accessToken,
                    request,
                );

            navigate(
                `/projects/${project.project_id}`,
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
                'プロジェクト登録中に予期しないエラーが発生しました。',
            );
        } finally {
            setIsSubmitting(false);
        }
    }

    return (
        <div className="project-list-page">
            <AppSidebar
                displayName={
                    user.display_name
                }
                role={user.role}
            />

            <div className="project-list-page__body">
                <AppHeader
                    title="プロジェクト登録"
                    displayName={
                        user.display_name
                    }
                />

                <main className="project-list-page__main">
                    <section className="project-list-page__heading">
                        <div>
                            <h1 className="project-list-page__title">
                                プロジェクトの登録
                            </h1>

                            <p className="project-list-page__description">
                                新しいプロジェクトを作成します。
                            </p>
                        </div>
                    </section>

                    <ProjectForm
                        users={users}
                        isUsersLoading={
                            isUsersLoading
                        }
                        usersErrorMessage={
                            usersErrorMessage
                        }
                        isSubmitting={
                            isSubmitting
                        }
                        submitErrorMessage={
                            submitErrorMessage
                        }
                        submitLabel="登録する"
                        onSubmit={handleSubmit}
                        onCancel={() => {
                            navigate('/projects');
                        }}
                    />
                </main>
            </div>
        </div>
    );
}