import {
    useEffect,
    useMemo,
    useState,
} from 'react';
import {
    Navigate,
    useNavigate,
} from 'react-router-dom';

import { AppHeader } from '../../components/layout/AppHeader';
import { AppSidebar } from '../../components/layout/AppSidebar';
import { ProjectFilters } from '../../components/projects/ProjectFilters';
import { ProjectTable } from '../../components/projects/ProjectTable';
import { useAuth } from '../../contexts/AuthContext';
import { useProjects } from '../../hooks/useProjects';

import type {
    ProjectStatusFilter,
} from '../../types/project';

const PROJECTS_PER_PAGE = 10;

export function ProjectListPage() {
    const navigate = useNavigate();
    const { status, user } = useAuth();

    const {
        projects,
        isLoading,
        errorMessage,
    } = useProjects();

    const [searchKeyword, setSearchKeyword] =
        useState('');

    const [statusFilter, setStatusFilter] =
        useState<ProjectStatusFilter>('all');

    const [currentPage, setCurrentPage] =
        useState(1);

    const filteredProjects = useMemo(() => {
        const normalizedKeyword =
            searchKeyword
                .trim()
                .toLowerCase();

        return projects.filter((project) => {
            const normalizedStatus =
                project.status === 'stopped'
                    ? 'paused'
                    : project.status;

            const matchesStatus =
                statusFilter === 'all' ||
                normalizedStatus ===
                    statusFilter;

            if (!normalizedKeyword) {
                return matchesStatus;
            }

            const matchesKeyword = [
                project.name,
                project.slug,
                project.description,
            ].some((value) =>
                value
                    .toLowerCase()
                    .includes(
                        normalizedKeyword,
                    ),
            );

            return (
                matchesKeyword &&
                matchesStatus
            );
        });
    }, [
        projects,
        searchKeyword,
        statusFilter,
    ]);

    const totalItems =
        filteredProjects.length;

    const totalPages = Math.max(
        1,
        Math.ceil(
            totalItems /
                PROJECTS_PER_PAGE,
        ),
    );

    const displayedProjects =
        useMemo(() => {
            const startIndex =
                (currentPage - 1) *
                PROJECTS_PER_PAGE;

            return filteredProjects.slice(
                startIndex,
                startIndex +
                    PROJECTS_PER_PAGE,
            );
        }, [
            currentPage,
            filteredProjects,
        ]);

    const firstItem =
        totalItems === 0
            ? 0
            : (currentPage - 1) *
                  PROJECTS_PER_PAGE +
              1;

    const lastItem = Math.min(
        currentPage *
            PROJECTS_PER_PAGE,
        totalItems,
    );

    useEffect(() => {
        setCurrentPage(1);
    }, [
        searchKeyword,
        statusFilter,
    ]);

    useEffect(() => {
        if (currentPage > totalPages) {
            setCurrentPage(totalPages);
        }
    }, [
        currentPage,
        totalPages,
    ]);

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
        <div className="project-list-page">
            <AppSidebar
                displayName={
                    user.display_name
                }
                role={user.role}
            />

            <div className="project-list-page__body">
                <AppHeader
                    title="プロジェクト"
                    displayName={
                        user.display_name
                    }
                />

                <main className="project-list-page__main">
                    <section className="project-list-page__heading">
                        <div>
                            <h1 className="project-list-page__title">
                                プロジェクト
                            </h1>

                            <p className="project-list-page__description">
                                参加中のプロジェクトを確認・管理できます。
                            </p>
                        </div>

                        <button
                            className="project-list-page__add-button"
                            type="button"
                            onClick={() => {
                                navigate(
                                    '/projects/new',
                                );
                            }}
                        >
                            <span aria-hidden="true">
                                ＋
                            </span>

                            <span>
                                プロジェクトを追加
                            </span>
                        </button>
                    </section>

                    <ProjectFilters
                        searchKeyword={
                            searchKeyword
                        }
                        statusFilter={
                            statusFilter
                        }
                        onSearchKeywordChange={
                            setSearchKeyword
                        }
                        onStatusFilterChange={
                            setStatusFilter
                        }
                    />

                    <ProjectTable
                        projects={
                            displayedProjects
                        }
                        isLoading={isLoading}
                        errorMessage={
                            errorMessage
                        }
                        currentPage={
                            currentPage
                        }
                        totalPages={
                            totalPages
                        }
                        totalItems={
                            totalItems
                        }
                        firstItem={
                            firstItem
                        }
                        lastItem={
                            lastItem
                        }
                        onPageChange={
                            setCurrentPage
                        }
                    />
                </main>
            </div>
        </div>
    );
}