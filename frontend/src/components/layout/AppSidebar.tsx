import { NavLink } from 'react-router-dom';

type AppSidebarProps = {
    displayName: string;
    role: string;
};

function getUserInitial(
    displayName: string,
): string {
    const normalizedDisplayName =
        displayName.trim();

    if (!normalizedDisplayName) {
        return '?';
    }

    return Array.from(
        normalizedDisplayName,
    )[0];
}

function getRoleLabel(
    role: string,
): string {
    switch (role) {
        case 'owner':
            return 'オーナー';

        case 'admin':
            return '管理者';

        case 'member':
            return 'メンバー';

        default:
            return role || '権限未設定';
    }
}

export function AppSidebar({
    displayName,
    role,
}: AppSidebarProps) {
    const userInitial =
        getUserInitial(displayName);

    return (
        <aside className="app-sidebar">
            <div className="app-sidebar__brand">
                <span
                    className="app-sidebar__brand-mark"
                    aria-hidden="true"
                >
                    F
                </span>

                <span className="app-sidebar__brand-name">
                    FlowSync Pulse
                </span>
            </div>

            <nav
                className="app-sidebar__navigation"
                aria-label="メインナビゲーション"
            >
                <NavLink
                    className={({ isActive }) =>
                        [
                            'app-sidebar__navigation-link',
                            isActive
                                ? 'app-sidebar__navigation-link--active'
                                : '',
                        ]
                            .filter(Boolean)
                            .join(' ')
                    }
                    to="/office"
                >
                    <span
                        className="app-sidebar__navigation-icon"
                        aria-hidden="true"
                    >
                        ◉
                    </span>

                    <span>オフィス</span>
                </NavLink>

                <NavLink
                    className={({ isActive }) =>
                        [
                            'app-sidebar__navigation-link',
                            isActive
                                ? 'app-sidebar__navigation-link--active'
                                : '',
                        ]
                            .filter(Boolean)
                            .join(' ')
                    }
                    to="/projects"
                >
                    <span
                        className="app-sidebar__navigation-icon"
                        aria-hidden="true"
                    >
                        ▣
                    </span>

                    <span>
                        プロジェクト
                    </span>
                </NavLink>

                <span
                    className={[
                        'app-sidebar__navigation-link',
                        'app-sidebar__navigation-link--disabled',
                    ].join(' ')}
                    aria-disabled="true"
                >
                    <span
                        className="app-sidebar__navigation-icon"
                        aria-hidden="true"
                    >
                        ✓
                    </span>

                    <span>タスク</span>
                </span>
            </nav>

            <div className="app-sidebar__user">
                <span
                    className="app-sidebar__user-avatar"
                    aria-hidden="true"
                >
                    {userInitial}
                </span>

                <div className="app-sidebar__user-information">
                    <p className="app-sidebar__user-name">
                        {displayName}
                    </p>

                    <p className="app-sidebar__user-role">
                        {getRoleLabel(role)}
                    </p>
                </div>
            </div>
        </aside>
    );
}