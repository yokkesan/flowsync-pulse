import { getProjectStatusView } from '../../utils/projectFormatters';

type ProjectStatusBadgeProps = {
    status: string;
};

export function ProjectStatusBadge({
    status,
}: ProjectStatusBadgeProps) {
    const statusView =
        getProjectStatusView(status);

    return (
        <span
            className={[
                'project-status-badge',
                `project-status-badge--${statusView.modifier}`,
            ].join(' ')}
        >
            {statusView.label}
        </span>
    );
}