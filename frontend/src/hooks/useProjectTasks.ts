import {
    useEffect,
    useState,
} from 'react';

import { ApiError } from '../services/authApi';
import { getAccessToken } from '../services/authStorage';
import { getProjectTasks } from '../services/taskApi';
import type { Task } from '../types/task';

type UseProjectTasksResult = {
    tasks: Task[];
    isLoading: boolean;
    errorMessage: string;
};

export function useProjectTasks(
    projectId: number,
): UseProjectTasksResult {
    const [tasks, setTasks] =
        useState<Task[]>([]);

    const [isLoading, setIsLoading] =
        useState(true);

    const [
        errorMessage,
        setErrorMessage,
    ] = useState('');

    useEffect(() => {
        let isMounted = true;

        async function loadProjectTasks():
        Promise<void> {
            if (
                !Number.isSafeInteger(projectId) ||
                projectId <= 0
            ) {
                if (isMounted) {
                    setTasks([]);
                    setErrorMessage(
                        'プロジェクトIDが正しくありません。',
                    );
                    setIsLoading(false);
                }

                return;
            }

            const accessToken =
                getAccessToken();

            if (!accessToken) {
                if (isMounted) {
                    setTasks([]);
                    setErrorMessage(
                        'ログイン情報が見つかりません。',
                    );
                    setIsLoading(false);
                }

                return;
            }

            if (isMounted) {
                setIsLoading(true);
                setErrorMessage('');
            }

            try {
                const response =
                    await getProjectTasks(
                        accessToken,
                        projectId,
                    );

                if (!isMounted) {
                    return;
                }

                setTasks(response.tasks);
            } catch (error) {
                if (!isMounted) {
                    return;
                }

                setTasks([]);

                if (error instanceof ApiError) {
                    setErrorMessage(
                        error.message,
                    );
                    return;
                }

                setErrorMessage(
                    'プロジェクトのタスク一覧取得中に予期しないエラーが発生しました。',
                );
            } finally {
                if (isMounted) {
                    setIsLoading(false);
                }
            }
        }

        void loadProjectTasks();

        return () => {
            isMounted = false;
        };
    }, [projectId]);

    return {
        tasks,
        isLoading,
        errorMessage,
    };
}