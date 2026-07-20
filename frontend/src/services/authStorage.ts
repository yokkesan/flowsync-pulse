const ACCESS_TOKEN_KEY =
    'flowsync_pulse_access_token';

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