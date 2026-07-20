import {
    useEffect,
    useState,
    type ReactNode,
} from 'react';
import {
    Navigate,
    useLocation,
} from 'react-router-dom';

import {
    ApiError,
    getCurrentUser,
} from '../../services/authApi';
import {
    getAccessToken,
    removeAccessToken,
} from '../../services/authStorage';

type ProtectedRouteProps = {
    children: ReactNode;
};

type AuthenticationStatus =
    | 'checking'
    | 'authenticated'
    | 'unauthenticated';

export function ProtectedRoute({
    children,
}: ProtectedRouteProps) {
    const location = useLocation();

    const [
        authenticationStatus,
        setAuthenticationStatus,
    ] = useState<AuthenticationStatus>('checking');

    useEffect(() => {
        let isMounted = true;

        async function verifyAuthentication() {
            const accessToken = getAccessToken();

            if (!accessToken) {
                if (isMounted) {
                    setAuthenticationStatus('unauthenticated');
                }

                return;
            }

            try {
                await getCurrentUser(accessToken);

                if (isMounted) {
                    setAuthenticationStatus('authenticated');
                }
            } catch (error) {
                if (
                    error instanceof ApiError &&
                    error.status === 401
                ) {
                    removeAccessToken();
                }

                if (isMounted) {
                    setAuthenticationStatus('unauthenticated');
                }
            }
        }

        void verifyAuthentication();

        return () => {
            isMounted = false;
        };
    }, []);

    if (authenticationStatus === 'checking') {
        return (
            <div role="status">
                ログイン状態を確認しています。
            </div>
        );
    }

    if (authenticationStatus === 'unauthenticated') {
        return (
            <Navigate
                to="/login"
                replace
                state={{
                    from: location.pathname,
                }}
            />
        );
    }

    return children;
}