const ACCESS_TOKEN_KEY =
    'flowsync_pulse_access_token';

const LEGACY_CURRENT_USER_KEY =
    'currentUser';

const LEGACY_AUTH_USER_KEY =
    'flowsync_pulse_current_user';

export function saveAccessToken(
    accessToken: string,
): void {
    sessionStorage.setItem(
        ACCESS_TOKEN_KEY,
        accessToken,
    );
}

export function getAccessToken(): string | null {
    return sessionStorage.getItem(
        ACCESS_TOKEN_KEY,
    );
}

export function removeAccessToken(): void {
    sessionStorage.removeItem(
        ACCESS_TOKEN_KEY,
    );
}

export function removeLegacyAuthData(): void {
    sessionStorage.removeItem(
        LEGACY_CURRENT_USER_KEY,
    );

    sessionStorage.removeItem(
        LEGACY_AUTH_USER_KEY,
    );
}

export function clearAuthStorage(): void {
    removeAccessToken();
    removeLegacyAuthData();
}