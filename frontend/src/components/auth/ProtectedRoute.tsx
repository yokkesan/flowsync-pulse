import type { ReactNode } from 'react';
import {
    Navigate,
    useLocation,
} from 'react-router-dom';

import { useAuth } from '../../contexts/AuthContext';

type ProtectedRouteProps = {
    children: ReactNode;
};

export function ProtectedRoute({
    children,
}: ProtectedRouteProps) {
    const location = useLocation();
    const { status } = useAuth();

    if (status === 'checking') {
        return (
            <div role="status">
                ログイン状態を確認しています。
            </div>
        );
    }

    if (
        status === 'unauthenticated' ||
        status === 'error'
    ) {
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