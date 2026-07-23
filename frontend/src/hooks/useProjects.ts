import {
    useCallback,
    useEffect,
    useState,
} from 'react';

import { ApiError } from '../services/authApi';
import { getAccessToken } from '../services/authStorage';
import { getProjects } from '../services/projectApi';

import type { Project } from '../types/project';

type UseProjectsResult = {
    projects: Project[];
    isLoading: boolean;
    errorMessage: string;
    reloadProjects: () => Promise<void>;
};

export function useProjects(): UseProjectsResult {
    const [projects, setProjects] =
        useState<Project[]>([]);

    const [isLoading, setIsLoading] =
        useState(true);

    const [errorMessage, setErrorMessage] =
        useState('');

    const reloadProjects =
        useCallback(async () => {
            const accessToken =
                getAccessToken();

            if (!accessToken) {
                setProjects([]);
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
                    await getProjects(
                        accessToken,
                    );

                setProjects(
                    response.projects,
                );
            } catch (error) {
                setProjects([]);

                if (error instanceof ApiError) {
                    setErrorMessage(
                        error.message,
                    );
                    return;
                }

                setErrorMessage(
                    'プロジェクト情報の取得中に予期しないエラーが発生しました。',
                );
            } finally {
                setIsLoading(false);
            }
        }, []);

    useEffect(() => {
        void reloadProjects();
    }, [reloadProjects]);

    return {
        projects,
        isLoading,
        errorMessage,
        reloadProjects,
    };
}