import {
    useEffect,
    useState,
} from 'react';

import { ApiError } from '../services/authApi';
import { getAccessToken } from '../services/authStorage';
import { getTasks } from '../services/taskApi';
import type { Task } from '../types/task';

type UseTasksResult = {
    tasks: Task[];
    isLoading: boolean;
    errorMessage: string;
};

export function useTasks(
    projectId: number,
): UseTasksResult {
    const [tasks, setTasks] = useState<Task[]>([]);

    const [isLoading, setIsLoading] =
        useState(true);

    const [
        errorMessage,
        setErrorMessage,
    ] = useState('');

    useEffect(() => {
        let isMounted = true;

        async function loadTasks():
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
                const response = await getTasks(
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
                    'タスク一覧の取得中に予期しないエラーが発生しました。',
                );
            } finally {
                if (isMounted) {
                    setIsLoading(false);
                }
            }
        }

        void loadTasks();

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