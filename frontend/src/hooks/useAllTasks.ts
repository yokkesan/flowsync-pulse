import {
    useEffect,
    useState,
} from 'react';

import { ApiError } from '../services/authApi';
import { getAccessToken } from '../services/authStorage';
import { getTasks } from '../services/taskApi';
import type { Task } from '../types/task';

type UseAllTasksResult = {
    tasks: Task[];
    isLoading: boolean;
    errorMessage: string;
};

export function useAllTasks():
UseAllTasksResult {
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

        async function loadAllTasks():
        Promise<void> {
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
                    await getTasks(
                        accessToken,
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

                if (
                    error instanceof ApiError
                ) {
                    setErrorMessage(
                        error.message,
                    );
                    return;
                }

                setErrorMessage(
                    '会社内のタスク一覧取得中に予期しないエラーが発生しました。',
                );
            } finally {
                if (isMounted) {
                    setIsLoading(false);
                }
            }
        }

        void loadAllTasks();

        return () => {
            isMounted = false;
        };
    }, []);

    return {
        tasks,
        isLoading,
        errorMessage,
    };
}