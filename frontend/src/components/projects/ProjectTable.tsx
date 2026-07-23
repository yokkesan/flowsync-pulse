import { ProjectPagination } from './ProjectPagination';
import { ProjectTableRow } from './ProjectTableRow';

import type { Project } from '../../types/project';

type ProjectTableProps = {
    projects: Project[];
    isLoading: boolean;
    errorMessage: string;
    currentPage: number;
    totalPages: number;
    totalItems: number;
    firstItem: number;
    lastItem: number;
    onPageChange: (
        page: number,
    ) => void;
};

export function ProjectTable({
    projects,
    isLoading,
    errorMessage,
    currentPage,
    totalPages,
    totalItems,
    firstItem,
    lastItem,
    onPageChange,
}: ProjectTableProps) {
    if (isLoading) {
        return (
            <section
                className="project-table"
                aria-label="プロジェクト一覧"
            >
                <div
                    className="project-table__state"
                    role="status"
                >
                    プロジェクトを読み込んでいます。
                </div>
            </section>
        );
    }

    if (errorMessage) {
        return (
            <section
                className="project-table"
                aria-label="プロジェクト一覧"
            >
                <div
                    className="project-table__state project-table__state--error"
                    role="alert"
                >
                    {errorMessage}
                </div>
            </section>
        );
    }

    return (
        <section
            className="project-table"
            aria-label="プロジェクト一覧"
        >
            <div className="project-table__wrapper">
                <table className="project-table__table">
                    <thead>
                        <tr>
                            <th scope="col">
                                プロジェクト
                            </th>

                            <th scope="col">
                                ステータス
                            </th>

                            <th scope="col">
                                期限
                            </th>

                            <th scope="col">
                                タスク
                            </th>

                            <th scope="col">
                                担当者
                            </th>

                            <th scope="col">
                                操作
                            </th>
                        </tr>
                    </thead>

                    <tbody>
                        {projects.length === 0 ? (
                            <tr>
                                <td
                                    className="project-table__empty-cell"
                                    colSpan={6}
                                >
                                    条件に一致するプロジェクトはありません。
                                </td>
                            </tr>
                        ) : (
                            projects.map(
                                (project) => (
                                    <ProjectTableRow
                                        key={
                                            project.project_id
                                        }
                                        project={
                                            project
                                        }
                                    />
                                ),
                            )
                        )}
                    </tbody>
                </table>
            </div>

            <ProjectPagination
                currentPage={currentPage}
                totalPages={totalPages}
                totalItems={totalItems}
                firstItem={firstItem}
                lastItem={lastItem}
                onPageChange={onPageChange}
            />
        </section>
    );
}