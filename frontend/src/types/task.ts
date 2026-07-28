export type TaskStatus =
    | 'not_started'
    | 'in_progress'
    | 'completed'
    | 'suspended';

export type TaskPriority =
    | 'high'
    | 'medium'
    | 'low';

export type Task = {
    task_id: number;
    project_id: number;
    project_name: string;
    name: string;
    description: string | null;
    assignee_user_id: number;
    assignee_name: string;
    status: TaskStatus;
    priority: TaskPriority;
    start_date: string | null;
    due_date: string | null;
    completed_at: string | null;
    created_at: string;
    updated_at: string;
};

export type TaskListResponse = {
    tasks: Task[];
};

export type TaskWriteRequest = {
    name: string;
    description: string | null;
    assignee_user_id: number;
    status: TaskStatus;
    priority: TaskPriority;
    start_date: string | null;
    due_date: string | null;
};

export type CreateTaskRequest =
    TaskWriteRequest;

export type UpdateTaskRequest =
    TaskWriteRequest;