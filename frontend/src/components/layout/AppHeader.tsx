type AppHeaderProps = {
    title: string;
    displayName: string;
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

export function AppHeader({
    title,
    displayName,
}: AppHeaderProps) {
    const userInitial =
        getUserInitial(displayName);

    return (
        <header className="app-header">
            <button
                className="app-header__page-selector"
                type="button"
            >
                <span>{title}</span>

                <span aria-hidden="true">
                    ▾
                </span>
            </button>

            <div className="app-header__actions">
                <button
                    className="app-header__icon-button"
                    type="button"
                    aria-label="通知"
                >
                    ♢
                </button>

                <button
                    className="app-header__profile-button"
                    type="button"
                    aria-label={`${displayName}のプロフィールメニュー`}
                >
                    <span
                        className="app-header__profile-avatar"
                        aria-hidden="true"
                    >
                        {userInitial}
                    </span>

                    <span aria-hidden="true">
                        ▾
                    </span>
                </button>
            </div>
        </header>
    );
}