import { getMemberInitial } from '../../utils/projectFormatters';

import type { ProjectMember } from '../../types/project';

type ProjectMemberAvatarsProps = {
    members: ProjectMember[];
    maxVisible?: number;
};

export function ProjectMemberAvatars({
    members,
    maxVisible = 3,
}: ProjectMemberAvatarsProps) {
    const visibleMembers =
        members.slice(0, maxVisible);

    const remainingCount =
        members.length -
        visibleMembers.length;

    if (members.length === 0) {
        return (
            <span className="project-member-avatars__empty">
                未設定
            </span>
        );
    }

    return (
        <div className="project-member-avatars">
            <div className="project-member-avatars__list">
                {visibleMembers.map((member) => (
                    <span
                        className="project-member-avatars__avatar"
                        key={member.user_id}
                        title={member.display_name}
                    >
                        {getMemberInitial(
                            member.display_name,
                        )}
                    </span>
                ))}

                {remainingCount > 0 && (
                    <span className="project-member-avatars__remaining">
                        +{remainingCount}
                    </span>
                )}
            </div>
        </div>
    );
}