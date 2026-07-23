import { useNavigate } from 'react-router-dom';

import { ProjectMemberAvatars } from './ProjectMemberAvatars';
import { ProjectStatusBadge } from './ProjectStatusBadge';

import { formatProjectDate } from '../../utils/projectFormatters';

import type { Project } from '../../types/project';

type ProjectTableRowProps = {
    project: Project;
};

export function ProjectTableRow({
    project,
}: ProjectTableRowProps) {
    const navigate = useNavigate();

    function handleProjectSelect(): void {
        navigate(
            `/projects/${project.project_id}`,
        );
    }

    return (
        <tr className="project-table__row">
            <td className="project-table__project-cell">
                <button
                    className="project-table__project-link"
                    type="button"
                    onClick={handleProjectSelect}
                >
                    {project.name}
                </button>

                <p className="project-table__project-description">
                    {project.description ||
                        project.slug}
                </p>
            </td>

            <td>
                <ProjectStatusBadge
                    status={project.status}
                />
            </td>

            <td>
                {formatProjectDate(
                    project.end_date,
                )}
            </td>

            <td>
                <span className="project-table__task-count">
                    {project.task_count}
                </span>
            </td>

            <td>
                <ProjectMemberAvatars
                    members={project.members}
                />
            </td>

            <td>
                <button
                    className="project-table__menu-button"
                    type="button"
                    aria-label={`${project.name}の操作メニュー`}
                >
                    •••
                </button>
            </td>
        </tr>
    );
}