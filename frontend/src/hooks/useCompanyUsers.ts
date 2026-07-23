import {
    useCallback,
    useEffect,
    useState,
} from 'react';

import { ApiError } from '../services/authApi';
import { getAccessToken } from '../services/authStorage';
import { getCompanyUsers } from '../services/userApi';

import type {
    CompanyUser,
} from '../types/user';

type UseCompanyUsersResult = {
    users: CompanyUser[];
    isLoading: boolean;
    errorMessage: string;
    reloadUsers: () => Promise<void>;
};

export function useCompanyUsers():
UseCompanyUsersResult {
    const [users, setUsers] =
        useState<CompanyUser[]>([]);

    const [isLoading, setIsLoading] =
        useState(true);

    const [errorMessage, setErrorMessage] =
        useState('');

    const reloadUsers =
        useCallback(async () => {
            const accessToken =
                getAccessToken();

            if (!accessToken) {
                setUsers([]);
                setErrorMessage(
                    'ログイン情報が見つかりません。',
                );
                setIsLoading(false);
                return;
            }

            setIsLoading(true);
            setErrorMessage('');

            try {
                const response =
                    await getCompanyUsers(
                        accessToken,
                    );

                setUsers(response.users);
            } catch (error) {
                setUsers([]);

                if (error instanceof ApiError) {
                    setErrorMessage(
                        error.message,
                    );
                    return;
                }

                setErrorMessage(
                    '所属ユーザー情報の取得中に予期しないエラーが発生しました。',
                );
            } finally {
                setIsLoading(false);
            }
        }, []);

    useEffect(() => {
        void reloadUsers();
    }, [reloadUsers]);

    return {
        users,
        isLoading,
        errorMessage,
        reloadUsers,
    };
}