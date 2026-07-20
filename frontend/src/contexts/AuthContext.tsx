import {
    createContext,
    useCallback,
    useContext,
    useEffect,
    useMemo,
    useState,
    type ReactNode,
} from 'react';

import {
    getCurrentUser,
    login as loginApi,
} from '../services/authApi';
import {
    clearAuthStorage,
    getAccessToken,
    removeLegacyAuthData,
    saveAccessToken,
} from '../services/authStorage';

import type {
    AuthCompany,
    AuthUser,
    LoginRequest,
} from '../types/auth';

type AuthenticationStatus =
    | 'checking'
    | 'authenticated'
    | 'unauthenticated'
    | 'error';

type AuthContextValue = {
    status: AuthenticationStatus;
    user: AuthUser | null;
    company: AuthCompany | null;
    login: (request: LoginRequest) => Promise<void>;
    logout: () => void;
    refreshAuthentication: () => Promise<void>;
};

type AuthProviderProps = {
    children: ReactNode;
};

const AuthContext =
    createContext<AuthContextValue | null>(null);

export function AuthProvider({
    children,
}: AuthProviderProps) {
    const [status, setStatus] =
        useState<AuthenticationStatus>('checking');

    const [user, setUser] =
        useState<AuthUser | null>(null);

    const [company, setCompany] =
        useState<AuthCompany | null>(null);

    const clearAuthentication = useCallback(() => {
        clearAuthStorage();
        setUser(null);
        setCompany(null);
        setStatus('unauthenticated');
    }, []);

    const refreshAuthentication =
        useCallback(async (): Promise<void> => {
            const accessToken = getAccessToken();

            if (!accessToken) {
                removeLegacyAuthData();
                setUser(null);
                setCompany(null);
                setStatus('unauthenticated');
                return;
            }

            setStatus('checking');

            try {
                const response =
                    await getCurrentUser(accessToken);

                setUser(response.user);
                setCompany(response.company);
                setStatus('authenticated');

                removeLegacyAuthData();
            } catch {
                clearAuthentication();
            }
        }, [clearAuthentication]);

    const login = useCallback(
        async (
            request: LoginRequest,
        ): Promise<void> => {
            const response =
                await loginApi(request);

            saveAccessToken(
                response.access_token,
            );

            try {
                const currentUserResponse =
                    await getCurrentUser(
                        response.access_token,
                    );

                setUser(currentUserResponse.user);
                setCompany(
                    currentUserResponse.company,
                );
                setStatus('authenticated');

                removeLegacyAuthData();
            } catch (error) {
                clearAuthentication();
                throw error;
            }
        },
        [clearAuthentication],
    );

    const logout = useCallback(() => {
        clearAuthentication();
    }, [clearAuthentication]);

    useEffect(() => {
        void refreshAuthentication();
    }, [refreshAuthentication]);

    const contextValue =
        useMemo<AuthContextValue>(
            () => ({
                status,
                user,
                company,
                login,
                logout,
                refreshAuthentication,
            }),
            [
                status,
                user,
                company,
                login,
                logout,
                refreshAuthentication,
            ],
        );

    return (
        <AuthContext.Provider
            value={contextValue}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth(): AuthContextValue {
    const context = useContext(AuthContext);

    if (!context) {
        throw new Error(
            'useAuthはAuthProvider内で使用してください。',
        );
    }

    return context;
}