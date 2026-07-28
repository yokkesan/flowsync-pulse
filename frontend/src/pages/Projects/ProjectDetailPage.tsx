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
import { ProjectMemberAvatars } from '../../components/projects/ProjectMemberAvatars';
import { ProjectStatusBadge } from '../../components/projects/ProjectStatusBadge';
import { useAuth } from '../../contexts/AuthContext';
import { useProject } from '../../hooks/useProject';
import { ApiError } from '../../services/authApi';
import { getAccessToken } from '../../services/authStorage';
import { deleteProject } from '../../services/projectApi';
import { formatProjectDate } from '../../utils/projectFormatters';

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

export function ProjectDetailPage() {
    const navigate = useNavigate();

    const { projectId: projectIdParam } =
        useParams();

    const { status, user } = useAuth();

    const projectId =
        parseProjectId(projectIdParam);

    const {
        project,
        isLoading,
        errorMessage,
    } = useProject(projectId ?? 0);

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

    async function handleDelete():
    Promise<void> {
        if (!project) {
            return;
        }

        const shouldDelete = window.confirm(
            `「${project.name}」を削除しますか？\nこの操作は取り消せません。`,
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
            await deleteProject(
                accessToken,
                project.project_id,
            );

            navigate('/projects', {
                replace: true,
            });
        } catch (error) {
            if (error instanceof ApiError) {
                setDeleteErrorMessage(
                    error.message,
                );
                return;
            }

            setDeleteErrorMessage(
                'プロジェクト削除中に予期しないエラーが発生しました。',
            );
        } finally {
            setIsDeleting(false);
        }
    }

    return (
        <div className="project-detail-page">
            <AppSidebar
                displayName={
                    user.display_name
                }
                role={user.role}
            />

            <div className="project-detail-page__body">
                <AppHeader
                    title="プロジェクト詳細"
                    displayName={
                        user.display_name
                    }
                />

                <main className="project-detail-page__main">
                    <div className="project-detail-page__back">
                        <button
                            className="project-detail-page__back-button"
                            type="button"
                            onClick={() => {
                                navigate(
                                    '/projects',
                                );
                            }}
                        >
                            ← プロジェクト一覧へ
                        </button>
                    </div>

                    {isLoading ? (
                        <section
                            className="project-detail-page__state"
                            role="status"
                        >
                            プロジェクト詳細を読み込んでいます。
                        </section>
                    ) : errorMessage ? (
                        <section
                            className="project-detail-page__state project-detail-page__state--error"
                            role="alert"
                        >
                            {errorMessage}
                        </section>
                    ) : !project ? (
                        <section className="project-detail-page__state">
                            プロジェクトが見つかりません。
                        </section>
                    ) : (
                        <>
                            <header className="project-detail-page__heading">
                                <div>
                                    <div className="project-detail-page__title-row">
                                        <h1 className="project-detail-page__title">
                                            {
                                                project.name
                                            }
                                        </h1>

                                        <ProjectStatusBadge
                                            status={
                                                project.status
                                            }
                                        />
                                    </div>

                                    <p className="project-detail-page__slug">
                                        {
                                            project.slug
                                        }
                                    </p>
                                </div>

                                <div className="project-detail-page__actions">
                                    <button
                                        className="project-detail-page__button project-detail-page__button--secondary"
                                        type="button"
                                        disabled={
                                            isDeleting
                                        }
                                        onClick={() => {
                                            navigate(
                                                `/projects/${project.project_id}/edit`,
                                            );
                                        }}
                                    >
                                        編集
                                    </button>

                                    <button
                                        className="project-detail-page__button project-detail-page__button--danger"
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

                            <nav
                                className="project-detail-page__tabs"
                                aria-label="プロジェクト詳細メニュー"
                            >
                                <button
                                    className="project-detail-page__tab project-detail-page__tab--active"
                                    type="button"
                                >
                                    概要
                                </button>

                                <button
                                    className="project-detail-page__tab"
                                    type="button"
                                    onClick={() => {
                                        navigate(
                                            `/projects/${project.project_id}/tasks`,
                                        );
                                    }}
                                >
                                    タスク
                                </button>
                            </nav>

                            {deleteErrorMessage && (
                                <div
                                    className="project-detail-page__delete-error"
                                    role="alert"
                                >
                                    {
                                        deleteErrorMessage
                                    }
                                </div>
                            )}

                            <section className="project-detail-page__content">
                                <div className="project-detail-page__card">
                                    <h2 className="project-detail-page__card-title">
                                        基本情報
                                    </h2>

                                    <dl className="project-detail-page__definition-list">
                                        <div className="project-detail-page__definition-item">
                                            <dt>
                                                説明
                                            </dt>

                                            <dd>
                                                {project.description ||
                                                    '未設定'}
                                            </dd>
                                        </div>

                                        <div className="project-detail-page__definition-item">
                                            <dt>
                                                リポジトリURL
                                            </dt>

                                            <dd>
                                                {project.repository_url ? (
                                                    <a
                                                        href={
                                                            project.repository_url
                                                        }
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                    >
                                                        {
                                                            project.repository_url
                                                        }
                                                    </a>
                                                ) : (
                                                    '未設定'
                                                )}
                                            </dd>
                                        </div>

                                        <div className="project-detail-page__definition-item">
                                            <dt>
                                                開始日
                                            </dt>

                                            <dd>
                                                {formatProjectDate(
                                                    project.start_date,
                                                )}
                                            </dd>
                                        </div>

                                        <div className="project-detail-page__definition-item">
                                            <dt>
                                                終了日
                                            </dt>

                                            <dd>
                                                {formatProjectDate(
                                                    project.end_date,
                                                )}
                                            </dd>
                                        </div>

                                        <div className="project-detail-page__definition-item">
                                            <dt>
                                                タスク数
                                            </dt>

                                            <dd>
                                                {
                                                    project.task_count
                                                }
                                            </dd>
                                        </div>
                                    </dl>
                                </div>

                                <div className="project-detail-page__card">
                                    <h2 className="project-detail-page__card-title">
                                        メンバー
                                    </h2>

                                    <ProjectMemberAvatars
                                        members={
                                            project.members
                                        }
                                        maxVisible={
                                            project.members
                                                .length
                                        }
                                    />

                                    <ul className="project-detail-page__member-list">
                                        {project.members.map(
                                            (
                                                member,
                                            ) => (
                                                <li
                                                    className="project-detail-page__member"
                                                    key={
                                                        member.user_id
                                                    }
                                                >
                                                    <span className="project-detail-page__member-name">
                                                        {
                                                            member.display_name
                                                        }
                                                    </span>

                                                    <span className="project-detail-page__member-role">
                                                        {
                                                            member.role
                                                        }
                                                    </span>
                                                </li>
                                            ),
                                        )}
                                    </ul>
                                </div>
                            </section>
                        </>
                    )}
                </main>
            </div>
        </div>
    );
}