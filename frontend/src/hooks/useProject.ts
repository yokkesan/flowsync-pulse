import {
    useCallback,
    useEffect,
    useState,
} from 'react';

import { ApiError } from '../services/authApi';
import { getAccessToken } from '../services/authStorage';
import { getProject } from '../services/projectApi';

import type { Project } from '../types/project';

type UseProjectResult = {
    project: Project | null;
    isLoading: boolean;
    errorMessage: string;
    reloadProject: () => Promise<void>;
};

export function useProject(
    projectId: number,
): UseProjectResult {
    const [project, setProject] =
        useState<Project | null>(null);

    const [isLoading, setIsLoading] =
        useState(true);

    const [errorMessage, setErrorMessage] =
        useState('');

    const reloadProject =
        useCallback(async () => {
            const accessToken =
                getAccessToken();

            if (!accessToken) {
                setProject(null);
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
                    await getProject(
                        accessToken,
                        projectId,
                    );

                setProject(response);
            } catch (error) {
                setProject(null);

                if (error instanceof ApiError) {
                    setErrorMessage(
                        error.message,
                    );
                    return;
                }

                setErrorMessage(
                    'プロジェクト詳細の取得中に予期しないエラーが発生しました。',
                );
            } finally {
                setIsLoading(false);
            }
        }, [projectId]);

    useEffect(() => {
        void reloadProject();
    }, [reloadProject]);

    return {
        project,
        isLoading,
        errorMessage,
        reloadProject,
    };
}