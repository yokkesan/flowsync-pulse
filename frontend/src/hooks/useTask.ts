import {
    useEffect,
    useState,
} from 'react';

import { ApiError } from '../services/authApi';
import { getAccessToken } from '../services/authStorage';
import { getTask } from '../services/taskApi';
import type { Task } from '../types/task';

type UseTaskResult = {
    task: Task | null;
    isLoading: boolean;
    errorMessage: string;
};

export function useTask(
    projectId: number,
    taskId: number,
): UseTaskResult {
    const [task, setTask] =
        useState<Task | null>(null);

    const [isLoading, setIsLoading] =
        useState(true);

    const [
        errorMessage,
        setErrorMessage,
    ] = useState('');

    useEffect(() => {
        let isMounted = true;

        async function loadTask():
        Promise<void> {
            if (
                !Number.isSafeInteger(projectId) ||
                projectId <= 0 ||
                !Number.isSafeInteger(taskId) ||
                taskId <= 0
            ) {
                if (isMounted) {
                    setTask(null);
                    setErrorMessage(
                        'プロジェクトIDまたはタスクIDが正しくありません。',
                    );
                    setIsLoading(false);
                }

                return;
            }

            const accessToken =
                getAccessToken();

            if (!accessToken) {
                if (isMounted) {
                    setTask(null);
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
                const response = await getTask(
                    accessToken,
                    projectId,
                    taskId,
                );

                if (!isMounted) {
                    return;
                }

                setTask(response);
            } catch (error) {
                if (!isMounted) {
                    return;
                }

                setTask(null);

                if (error instanceof ApiError) {
                    setErrorMessage(
                        error.message,
                    );
                    return;
                }

                setErrorMessage(
                    'タスク詳細の取得中に予期しないエラーが発生しました。',
                );
            } finally {
                if (isMounted) {
                    setIsLoading(false);
                }
            }
        }

        void loadTask();

        return () => {
            isMounted = false;
        };
    }, [
        projectId,
        taskId,
    ]);

    return {
        task,
        isLoading,
        errorMessage,
    };
}