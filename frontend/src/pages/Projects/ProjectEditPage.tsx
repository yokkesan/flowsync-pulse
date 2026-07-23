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
    ProjectForm,
    type ProjectFormValues,
} from '../../components/projects/ProjectForm';
import { useAuth } from '../../contexts/AuthContext';
import { useCompanyUsers } from '../../hooks/useCompanyUsers';
import { useProject } from '../../hooks/useProject';
import { ApiError } from '../../services/authApi';
import { getAccessToken } from '../../services/authStorage';
import { updateProject } from '../../services/projectApi';

import type {
    ProjectWriteRequest,
} from '../../types/project';

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

export function ProjectEditPage() {
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

    const initialValues =
        useMemo<ProjectFormValues | undefined>(
            () => {
                if (!project) {
                    return undefined;
                }

                return {
                    name: project.name,
                    slug: project.slug,
                    description:
                        project.description ?? '',
                    status: project.status,
                    startDate:
                        project.start_date ?? '',
                    endDate:
                        project.end_date ?? '',
                    memberIds:
                        project.members.map(
                            (member) =>
                                member.user_id,
                        ),
                };
            },
            [project],
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

    async function handleSubmit(
        request: ProjectWriteRequest,
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
            const updatedProject =
                await updateProject(
                    accessToken,
                    projectId,
                    request,
                );

            navigate(
                `/projects/${updatedProject.project_id}`,
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
                'プロジェクト編集中に予期しないエラーが発生しました。',
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
                    title="プロジェクト編集"
                    displayName={
                        user.display_name
                    }
                />

                <main className="project-list-page__main">
                    <section className="project-list-page__heading">
                        <div>
                            <h1 className="project-list-page__title">
                                プロジェクトの編集
                            </h1>

                            <p className="project-list-page__description">
                                プロジェクト情報を変更します。
                            </p>
                        </div>
                    </section>

                    {isProjectLoading ? (
                        <section
                            className="project-detail-page__state"
                            role="status"
                        >
                            プロジェクト情報を読み込んでいます。
                        </section>
                    ) : projectErrorMessage ? (
                        <section
                            className="project-detail-page__state project-detail-page__state--error"
                            role="alert"
                        >
                            {projectErrorMessage}
                        </section>
                    ) : !initialValues ? (
                        <section className="project-detail-page__state">
                            プロジェクトが見つかりません。
                        </section>
                    ) : (
                        <ProjectForm
                            initialValues={
                                initialValues
                            }
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
                            submitLabel="変更を保存"
                            onSubmit={
                                handleSubmit
                            }
                            onCancel={() => {
                                navigate(
                                    `/projects/${projectId}`,
                                );
                            }}
                        />
                    )}
                </main>
            </div>
        </div>
    );
}